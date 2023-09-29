package adminrpc

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/deltaswapio/deltaswap/node/pkg/watchers/evm/connectors"
	"github.com/holiman/uint256"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/exp/slices"

	"github.com/deltaswapio/deltaswap/node/pkg/db"
	"github.com/deltaswapio/deltaswap/node/pkg/governor"
	gossipv1 "github.com/deltaswapio/deltaswap/node/pkg/proto/gossip/v1"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/deltaswapio/deltaswap/node/pkg/common"
	nodev1 "github.com/deltaswapio/deltaswap/node/pkg/proto/node/v1"
	"github.com/deltaswapio/deltaswap/sdk/vaa"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

var (
	vaaInjectionsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "wormhole_vaa_injections_total",
			Help: "Total number of injected VAA queued for broadcast",
		})
)

type nodePrivilegedService struct {
	nodev1.UnimplementedNodePrivilegedServiceServer
	db            *db.Database
	injectC       chan<- *common.MessagePublication
	obsvReqSendC  chan<- *gossipv1.ObservationRequest
	logger        *zap.Logger
	signedInC     chan<- *gossipv1.SignedVAAWithQuorum
	governor      *governor.ChainGovernor
	evmConnector  connectors.Connector
	gsCache       sync.Map
	gk            *ecdsa.PrivateKey
	phylaxAddress ethcommon.Address
	rpcMap        map[string]string
}

func NewPrivService(
	db *db.Database,
	injectC chan<- *common.MessagePublication,
	obsvReqSendC chan<- *gossipv1.ObservationRequest,
	logger *zap.Logger,
	signedInC chan<- *gossipv1.SignedVAAWithQuorum,
	governor *governor.ChainGovernor,
	evmConnector connectors.Connector,
	gk *ecdsa.PrivateKey,
	phylaxAddress ethcommon.Address,
	rpcMap map[string]string,

) *nodePrivilegedService {
	return &nodePrivilegedService{
		db:            db,
		injectC:       injectC,
		obsvReqSendC:  obsvReqSendC,
		logger:        logger,
		signedInC:     signedInC,
		governor:      governor,
		evmConnector:  evmConnector,
		gk:            gk,
		phylaxAddress: phylaxAddress,
		rpcMap:        rpcMap,
	}
}

// adminPhylaxSetUpdateToVAA converts a nodev1.PhylaxSetUpdate message to its canonical VAA representation.
// Returns an error if the data is invalid.
func adminPhylaxSetUpdateToVAA(req *nodev1.PhylaxSetUpdate, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) {
	if len(req.Phylaxs) == 0 {
		return nil, errors.New("empty phylax set specified")
	}

	if len(req.Phylaxs) > common.MaxPhylaxCount {
		return nil, fmt.Errorf("too many phylaxs - %d, maximum is %d", len(req.Phylaxs), common.MaxPhylaxCount)
	}

	addrs := make([]ethcommon.Address, len(req.Phylaxs))
	for i, g := range req.Phylaxs {
		if !ethcommon.IsHexAddress(g.Pubkey) {
			return nil, fmt.Errorf("invalid pubkey format at index %d (%s)", i, g.Name)
		}

		ethAddr := ethcommon.HexToAddress(g.Pubkey)
		for j, pk := range addrs {
			if pk == ethAddr {
				return nil, fmt.Errorf("duplicate pubkey at index %d (duplicate of %d): %s", i, j, g.Name)
			}
		}

		addrs[i] = ethAddr
	}

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyPhylaxSetUpdate{
			Keys:     addrs,
			NewIndex: phylaxSetIndex + 1,
		}.Serialize())

	return v, nil
}

// adminContractUpgradeToVAA converts a nodev1.ContractUpgrade message to its canonical VAA representation.
// Returns an error if the data is invalid.
func adminContractUpgradeToVAA(req *nodev1.ContractUpgrade, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) {
	b, err := hex.DecodeString(req.NewContract)
	if err != nil {
		return nil, errors.New("invalid new contract address encoding (expected hex)")
	}

	if len(b) != 32 {
		return nil, errors.New("invalid new_contract address")
	}

	if req.ChainId > math.MaxUint16 {
		return nil, errors.New("invalid chain_id")
	}

	newContractAddress := vaa.Address{}
	copy(newContractAddress[:], b)

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyContractUpgrade{
			ChainID:     vaa.ChainID(req.ChainId),
			NewContract: newContractAddress,
		}.Serialize())

	return v, nil
}

// tokenBridgeRegisterChain converts a nodev1.TokenBridgeRegisterChain message to its canonical VAA representation.
// Returns an error if the data is invalid.
func tokenBridgeRegisterChain(req *nodev1.BridgeRegisterChain, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) {
	if req.ChainId > math.MaxUint16 {
		return nil, errors.New("invalid chain_id")
	}

	b, err := hex.DecodeString(req.EmitterAddress)
	if err != nil {
		return nil, errors.New("invalid emitter address encoding (expected hex)")
	}

	if len(b) != 32 {
		return nil, errors.New("invalid emitter address (expected 32 bytes)")
	}

	emitterAddress := vaa.Address{}
	copy(emitterAddress[:], b)

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyTokenBridgeRegisterChain{
			Module:         req.Module,
			ChainID:        vaa.ChainID(req.ChainId),
			EmitterAddress: emitterAddress,
		}.Serialize())

	return v, nil
}

// accountantModifyBalance converts a nodev1.AccountantModifyBalance message to its canonical VAA representation.
// Returns an error if the data is invalid.
func accountantModifyBalance(req *nodev1.AccountantModifyBalance, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) {
	if req.TargetChainId > math.MaxUint16 {
		return nil, errors.New("invalid target_chain_id")
	}
	if req.ChainId > math.MaxUint16 {
		return nil, errors.New("invalid chain_id")
	}
	if req.TokenChain > math.MaxUint16 {
		return nil, errors.New("invalid token_chain")
	}

	b, err := hex.DecodeString(req.TokenAddress)
	if err != nil {
		return nil, errors.New("invalid token address (expected hex)")
	}

	if len(b) != 32 {
		return nil, errors.New("invalid new token address (expected 32 bytes)")
	}

	if len(req.Reason) > 32 {
		return nil, errors.New("the reason should not be larger than 32 bytes")
	}

	amount_big := big.NewInt(0)
	amount_big, ok := amount_big.SetString(req.Amount, 10)
	if !ok {
		return nil, errors.New("invalid amount")
	}

	// uint256 has Bytes32 method for easier serialization
	amount, overflow := uint256.FromBig(amount_big)
	if overflow {
		return nil, errors.New("amount overflow")
	}

	tokenAdress := vaa.Address{}
	copy(tokenAdress[:], b)

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyAccountantModifyBalance{
			Module:        req.Module,
			TargetChainID: vaa.ChainID(req.TargetChainId),

			Sequence:     req.Sequence,
			ChainId:      vaa.ChainID(req.ChainId),
			TokenChain:   vaa.ChainID(req.TokenChain),
			TokenAddress: tokenAdress,
			Kind:         uint8(req.Kind),
			Amount:       amount,
			Reason:       req.Reason,
		}.Serialize())

	return v, nil
}

// tokenBridgeUpgradeContract converts a nodev1.TokenBridgeRegisterChain message to its canonical VAA representation.
// Returns an error if the data is invalid.
func tokenBridgeUpgradeContract(req *nodev1.BridgeUpgradeContract, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) {
	if req.TargetChainId > math.MaxUint16 {
		return nil, errors.New("invalid target_chain_id")
	}

	b, err := hex.DecodeString(req.NewContract)
	if err != nil {
		return nil, errors.New("invalid new contract address (expected hex)")
	}

	if len(b) != 32 {
		return nil, errors.New("invalid new contract address (expected 32 bytes)")
	}

	newContract := vaa.Address{}
	copy(newContract[:], b)

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyTokenBridgeUpgradeContract{
			Module:        req.Module,
			TargetChainID: vaa.ChainID(req.TargetChainId),
			NewContract:   newContract,
		}.Serialize())

	return v, nil
}

// deltachainStoreCode converts a nodev1.DeltachainStoreCode to its canonical VAA representation
// Returns an error if the data is invalid
func deltachainStoreCode(req *nodev1.DeltachainStoreCode, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) {
	// validate the length of the hex passed in
	b, err := hex.DecodeString(req.WasmHash)
	if err != nil {
		return nil, fmt.Errorf("invalid cosmwasm bytecode hash (expected hex): %w", err)
	}

	if len(b) != 32 {
		return nil, fmt.Errorf("invalid cosmwasm bytecode hash (expected 32 bytes but received %d bytes)", len(b))
	}

	wasmHash := [32]byte{}
	copy(wasmHash[:], b)

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyDeltachainStoreCode{
			WasmHash: wasmHash,
		}.Serialize())

	return v, nil
}

// deltachainInstantiateContract converts a nodev1.DeltachainInstantiateContract to its canonical VAA representation
// Returns an error if the data is invalid
func deltachainInstantiateContract(req *nodev1.DeltachainInstantiateContract, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) { //nolint:unparam // error is always nil but kept to mirror function signature of other functions
	instantiationParams_hash := vaa.CreateInstatiateCosmwasmContractHash(req.CodeId, req.Label, []byte(req.InstantiationMsg))

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyDeltachainInstantiateContract{
			InstantiationParamsHash: instantiationParams_hash,
		}.Serialize())

	return v, nil
}

// deltachainMigrateContract converts a nodev1.DeltachainMigrateContract to its canonical VAA representation
func deltachainMigrateContract(req *nodev1.DeltachainMigrateContract, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) { //nolint:unparam // error is always nil but kept to mirror function signature of other functions
	instantiationParams_hash := vaa.CreateMigrateCosmwasmContractHash(req.CodeId, req.Contract, []byte(req.InstantiationMsg))

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyDeltachainMigrateContract{
			MigrationParamsHash: instantiationParams_hash,
		}.Serialize())

	return v, nil
}

func deltachainWasmInstantiateAllowlist(
	req *nodev1.DeltachainWasmInstantiateAllowlist,
	timestamp time.Time,
	phylaxSetIndex uint32,
	nonce uint32,
	sequence uint64,
) (*vaa.VAA, error) { //nolint:unparam // error is always nil but kept to mirror function signature of other functions
	decodedAddr, err := sdktypes.GetFromBech32(req.Contract, "wormhole")
	if err != nil {
		return nil, err
	}

	var action vaa.GovernanceAction
	if req.Action == nodev1.DeltachainWasmInstantiateAllowlistAction_DELTACHAIN_WASM_INSTANTIATE_ALLOWLIST_ACTION_ADD {
		action = vaa.ActionAddWasmInstantiateAllowlist
	} else if req.Action == nodev1.DeltachainWasmInstantiateAllowlistAction_DELTACHAIN_WASM_INSTANTIATE_ALLOWLIST_ACTION_DELETE {
		action = vaa.ActionDeleteWasmInstantiateAllowlist
	} else {
		return nil, fmt.Errorf("unrecognized wasm instantiate allowlist action")
	}

	var decodedAddr32 [32]byte
	copy(decodedAddr32[:], decodedAddr)

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex, vaa.BodyDeltachainWasmAllowlistInstantiate{
		ContractAddr: decodedAddr32,
		CodeId:       req.CodeId,
	}.Serialize(action))

	return v, nil
}

func gatewayScheduleUpgrade(
	req *nodev1.GatewayScheduleUpgrade,
	timestamp time.Time,
	phylaxSetIndex uint32,
	nonce uint32,
	sequence uint64,
) (*vaa.VAA, error) { //nolint:unparam // error is always nil but kept to mirror function signature of other functions
	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex, vaa.BodyGatewayScheduleUpgrade{
		Name:   req.Name,
		Height: req.Height,
	}.Serialize())

	return v, nil
}

func gatewayCancelUpgrade(
	timestamp time.Time,
	phylaxSetIndex uint32,
	nonce uint32,
	sequence uint64,
) (*vaa.VAA, error) { //nolint:unparam // error is always nil but kept to mirror function signature of other functions
	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.EmptyPayloadVaa(vaa.GatewayModuleStr, vaa.ActionCancelUpgrade, vaa.ChainIDDeltachain),
	)

	return v, nil
}

func gatewayIbcComposabilityMwSetContract(
	req *nodev1.GatewayIbcComposabilityMwSetContract,
	timestamp time.Time,
	phylaxSetIndex uint32,
	nonce uint32,
	sequence uint64,
) (*vaa.VAA, error) {
	decodedAddr, err := sdktypes.GetFromBech32(req.Contract, "wormhole")
	if err != nil {
		return nil, err
	}

	var decodedAddr32 [32]byte
	copy(decodedAddr32[:], decodedAddr)

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex, vaa.BodyGatewayIbcComposabilityMwContract{
		ContractAddr: decodedAddr32,
	}.Serialize())

	return v, nil
}

// circleIntegrationUpdateWormholeFinality converts a nodev1.CircleIntegrationUpdateWormholeFinality to its canonical VAA representation
// Returns an error if the data is invalid
func circleIntegrationUpdateWormholeFinality(req *nodev1.CircleIntegrationUpdateWormholeFinality, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) {
	if req.TargetChainId > math.MaxUint16 {
		return nil, fmt.Errorf("invalid target chain id, must be <= %d", math.MaxUint16)
	}
	if req.Finality > math.MaxUint8 {
		return nil, fmt.Errorf("invalid finality, must be <= %d", math.MaxUint8)
	}
	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyCircleIntegrationUpdateWormholeFinality{
			TargetChainID: vaa.ChainID(req.TargetChainId),
			Finality:      uint8(req.Finality),
		}.Serialize())

	return v, nil
}

// circleIntegrationRegisterEmitterAndDomain converts a nodev1.CircleIntegrationRegisterEmitterAndDomain to its canonical VAA representation
// Returns an error if the data is invalid
func circleIntegrationRegisterEmitterAndDomain(req *nodev1.CircleIntegrationRegisterEmitterAndDomain, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) {
	if req.TargetChainId > math.MaxUint16 {
		return nil, fmt.Errorf("invalid target chain id, must be <= %d", math.MaxUint16)
	}
	if req.ForeignEmitterChainId > math.MaxUint16 {
		return nil, fmt.Errorf("invalid foreign emitter chain id, must be <= %d", math.MaxUint16)
	}
	b, err := hex.DecodeString(req.ForeignEmitterAddress)
	if err != nil {
		return nil, errors.New("invalid foreign emitter address encoding (expected hex)")
	}

	if len(b) != 32 {
		return nil, errors.New("invalid foreign emitter address (expected 32 bytes)")
	}

	foreignEmitterAddress := vaa.Address{}
	copy(foreignEmitterAddress[:], b)

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyCircleIntegrationRegisterEmitterAndDomain{
			TargetChainID:         vaa.ChainID(req.TargetChainId),
			ForeignEmitterChainId: vaa.ChainID(req.ForeignEmitterChainId),
			ForeignEmitterAddress: foreignEmitterAddress,
			CircleDomain:          req.CircleDomain,
		}.Serialize())

	return v, nil
}

// circleIntegrationUpgradeContractImplementation converts a nodev1.CircleIntegrationUpgradeContractImplementation to its canonical VAA representation
// Returns an error if the data is invalid
func circleIntegrationUpgradeContractImplementation(req *nodev1.CircleIntegrationUpgradeContractImplementation, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) {
	if req.TargetChainId > math.MaxUint16 {
		return nil, fmt.Errorf("invalid target chain id, must be <= %d", math.MaxUint16)
	}
	b, err := hex.DecodeString(req.NewImplementationAddress)
	if err != nil {
		return nil, errors.New("invalid new implementation address encoding (expected hex)")
	}

	if len(b) != 32 {
		return nil, errors.New("invalid new implementation address (expected 32 bytes)")
	}

	newImplementationAddress := vaa.Address{}
	copy(newImplementationAddress[:], b)

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyCircleIntegrationUpgradeContractImplementation{
			TargetChainID:            vaa.ChainID(req.TargetChainId),
			NewImplementationAddress: newImplementationAddress,
		}.Serialize())

	return v, nil
}

func ibcUpdateChannelChain(
	req *nodev1.IbcUpdateChannelChain,
	timestamp time.Time,
	phylaxSetIndex uint32,
	nonce uint32,
	sequence uint64,
) (*vaa.VAA, error) {
	// validate parameters
	if req.TargetChainId > math.MaxUint16 {
		return nil, fmt.Errorf("invalid target chain id, must be <= %d", math.MaxUint16)
	}

	if req.ChainId > math.MaxUint16 {
		return nil, fmt.Errorf("invalid chain id, must be <= %d", math.MaxUint16)
	}

	if len(req.ChannelId) > 64 {
		return nil, fmt.Errorf("invalid channel ID length, must be <= 64")
	}
	channelId := vaa.LeftPadIbcChannelId(req.ChannelId)

	var module string
	if req.Module == nodev1.IbcUpdateChannelChainModule_IBC_UPDATE_CHANNEL_CHAIN_MODULE_RECEIVER {
		module = vaa.IbcReceiverModuleStr
	} else if req.Module == nodev1.IbcUpdateChannelChainModule_IBC_UPDATE_CHANNEL_CHAIN_MODULE_TRANSLATOR {
		module = vaa.IbcTranslatorModuleStr
	} else {
		return nil, fmt.Errorf("unrecognized ibc update channel chain module")
	}

	// create governance VAA
	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyIbcUpdateChannelChain{
			TargetChainId: vaa.ChainID(req.TargetChainId),
			ChannelId:     channelId,
			ChainId:       vaa.ChainID(req.ChainId),
		}.Serialize(module))

	return v, nil
}

// wormholeRelayerSetDefaultDeliveryProvider converts a nodev1.DeltaswapRelayerSetDefaultDeliveryProvider message to its canonical VAA representation.
// Returns an error if the data is invalid.
func wormholeRelayerSetDefaultDeliveryProvider(req *nodev1.DeltaswapRelayerSetDefaultDeliveryProvider, timestamp time.Time, phylaxSetIndex uint32, nonce uint32, sequence uint64) (*vaa.VAA, error) {
	if req.ChainId > math.MaxUint16 {
		return nil, errors.New("invalid target_chain_id")
	}

	b, err := hex.DecodeString(req.NewDefaultDeliveryProviderAddress)
	if err != nil {
		return nil, errors.New("invalid new default delivery provider address (expected hex)")
	}

	if len(b) != 32 {
		return nil, errors.New("invalid new default delivery provider address (expected 32 bytes)")
	}

	NewDefaultDeliveryProviderAddress := vaa.Address{}
	copy(NewDefaultDeliveryProviderAddress[:], b)

	v := vaa.CreateGovernanceVAA(timestamp, nonce, sequence, phylaxSetIndex,
		vaa.BodyDeltaswapRelayerSetDefaultDeliveryProvider{
			ChainID:                           vaa.ChainID(req.ChainId),
			NewDefaultDeliveryProviderAddress: NewDefaultDeliveryProviderAddress,
		}.Serialize())

	return v, nil
}

func GovMsgToVaa(message *nodev1.GovernanceMessage, currentSetIndex uint32, timestamp time.Time) (*vaa.VAA, error) {
	var (
		v   *vaa.VAA
		err error
	)

	switch payload := message.Payload.(type) {
	case *nodev1.GovernanceMessage_PhylaxSet:
		v, err = adminPhylaxSetUpdateToVAA(payload.PhylaxSet, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_ContractUpgrade:
		v, err = adminContractUpgradeToVAA(payload.ContractUpgrade, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_BridgeRegisterChain:
		v, err = tokenBridgeRegisterChain(payload.BridgeRegisterChain, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_BridgeContractUpgrade:
		v, err = tokenBridgeUpgradeContract(payload.BridgeContractUpgrade, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_AccountantModifyBalance:
		v, err = accountantModifyBalance(payload.AccountantModifyBalance, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_DeltachainStoreCode:
		v, err = deltachainStoreCode(payload.DeltachainStoreCode, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_DeltachainInstantiateContract:
		v, err = deltachainInstantiateContract(payload.DeltachainInstantiateContract, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_DeltachainMigrateContract:
		v, err = deltachainMigrateContract(payload.DeltachainMigrateContract, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_DeltachainWasmInstantiateAllowlist:
		v, err = deltachainWasmInstantiateAllowlist(payload.DeltachainWasmInstantiateAllowlist, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_GatewayScheduleUpgrade:
		v, err = gatewayScheduleUpgrade(payload.GatewayScheduleUpgrade, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_GatewayCancelUpgrade:
		v, err = gatewayCancelUpgrade(timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_GatewayIbcComposabilityMwSetContract:
		v, err = gatewayIbcComposabilityMwSetContract(payload.GatewayIbcComposabilityMwSetContract, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_CircleIntegrationUpdateWormholeFinality:
		v, err = circleIntegrationUpdateWormholeFinality(payload.CircleIntegrationUpdateWormholeFinality, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_CircleIntegrationRegisterEmitterAndDomain:
		v, err = circleIntegrationRegisterEmitterAndDomain(payload.CircleIntegrationRegisterEmitterAndDomain, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_CircleIntegrationUpgradeContractImplementation:
		v, err = circleIntegrationUpgradeContractImplementation(payload.CircleIntegrationUpgradeContractImplementation, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_IbcUpdateChannelChain:
		v, err = ibcUpdateChannelChain(payload.IbcUpdateChannelChain, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	case *nodev1.GovernanceMessage_DeltaswapRelayerSetDefaultDeliveryProvider:
		v, err = wormholeRelayerSetDefaultDeliveryProvider(payload.DeltaswapRelayerSetDefaultDeliveryProvider, timestamp, currentSetIndex, message.Nonce, message.Sequence)
	default:
		panic(fmt.Sprintf("unsupported VAA type: %T", payload))
	}

	return v, err
}

func (s *nodePrivilegedService) InjectGovernanceVAA(ctx context.Context, req *nodev1.InjectGovernanceVAARequest) (*nodev1.InjectGovernanceVAAResponse, error) {
	s.logger.Info("governance VAA injected via admin socket", zap.String("request", req.String()))

	var (
		v   *vaa.VAA
		err error
	)

	timestamp := time.Unix(int64(req.Timestamp), 0)

	digests := make([][]byte, len(req.Messages))

	for i, message := range req.Messages {
		v, err = GovMsgToVaa(message, req.CurrentSetIndex, timestamp)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		// Generate digest of the unsigned VAA.
		digest := v.SigningDigest()

		s.logger.Info("governance VAA constructed",
			zap.Any("vaa", v),
			zap.String("digest", digest.String()),
		)

		vaaInjectionsTotal.Inc()

		s.injectC <- &common.MessagePublication{
			TxHash:           ethcommon.Hash{},
			Timestamp:        v.Timestamp,
			Nonce:            v.Nonce,
			Sequence:         v.Sequence,
			ConsistencyLevel: v.ConsistencyLevel,
			EmitterChain:     v.EmitterChain,
			EmitterAddress:   v.EmitterAddress,
			Payload:          v.Payload,
			Unreliable:       false,
		}

		digests[i] = digest.Bytes()
	}

	return &nodev1.InjectGovernanceVAAResponse{Digests: digests}, nil
}

// fetchMissing attempts to backfill a gap by fetching and storing missing signed VAAs from the network.
// Returns true if the gap was filled, false otherwise.
func (s *nodePrivilegedService) fetchMissing(
	ctx context.Context,
	nodes []string,
	c *http.Client,
	chain vaa.ChainID,
	addr string,
	seq uint64) (bool, error) {

	// shuffle the list of public RPC endpoints
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	for _, node := range nodes {
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(
			"%s/v1/signed_vaa/%d/%s/%d", node, chain, addr, seq), nil)
		if err != nil {
			return false, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := c.Do(req)
		if err != nil {
			s.logger.Warn("failed to fetch missing VAA",
				zap.String("node", node),
				zap.String("chain", chain.String()),
				zap.String("address", addr),
				zap.Uint64("sequence", seq),
				zap.Error(err),
			)
			continue
		}

		switch resp.StatusCode {
		case http.StatusNotFound:
			resp.Body.Close()
			continue
		case http.StatusOK:
			type getVaaResp struct {
				VaaBytes string `json:"vaaBytes"`
			}
			var respBody getVaaResp
			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
				resp.Body.Close()
				s.logger.Warn("failed to decode VAA response",
					zap.String("node", node),
					zap.String("chain", chain.String()),
					zap.String("address", addr),
					zap.Uint64("sequence", seq),
					zap.Error(err),
				)
				continue
			}

			// base64 decode the VAA bytes
			vaaBytes, err := base64.StdEncoding.DecodeString(respBody.VaaBytes)
			if err != nil {
				resp.Body.Close()
				s.logger.Warn("failed to decode VAA body",
					zap.String("node", node),
					zap.String("chain", chain.String()),
					zap.String("address", addr),
					zap.Uint64("sequence", seq),
					zap.Error(err),
				)
				continue
			}

			s.logger.Info("backfilled VAA",
				zap.Uint16("chain", uint16(chain)),
				zap.String("address", addr),
				zap.Uint64("sequence", seq),
				zap.Int("numBytes", len(vaaBytes)),
			)

			// Inject into the gossip signed VAA receive path.
			// This has the same effect as if the VAA was received from the network
			// (verifying signature, storing in local DB...).
			s.signedInC <- &gossipv1.SignedVAAWithQuorum{
				Vaa: vaaBytes,
			}

			resp.Body.Close()
			return true, nil
		default:
			resp.Body.Close()
			return false, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
		}
	}

	return false, nil
}

func (s *nodePrivilegedService) FindMissingMessages(ctx context.Context, req *nodev1.FindMissingMessagesRequest) (*nodev1.FindMissingMessagesResponse, error) {
	b, err := hex.DecodeString(req.EmitterAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid emitter address encoding: %v", err)
	}
	emitterAddress := vaa.Address{}
	copy(emitterAddress[:], b)

	ids, first, last, err := s.db.FindEmitterSequenceGap(db.VAAID{
		EmitterChain:   vaa.ChainID(req.EmitterChain),
		EmitterAddress: emitterAddress,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database operation failed: %v", err)
	}

	if req.RpcBackfill {
		c := &http.Client{}
		unfilled := make([]uint64, 0, len(ids))
		for _, id := range ids {
			if ok, err := s.fetchMissing(ctx, req.BackfillNodes, c, vaa.ChainID(req.EmitterChain), emitterAddress.String(), id); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to backfill VAA: %v", err)
			} else if ok {
				continue
			}
			unfilled = append(unfilled, id)
		}
		ids = unfilled
	}

	resp := make([]string, len(ids))
	for i, v := range ids {
		resp[i] = fmt.Sprintf("%d/%s/%d", req.EmitterChain, emitterAddress, v)
	}
	return &nodev1.FindMissingMessagesResponse{
		MissingMessages: resp,
		FirstSequence:   first,
		LastSequence:    last,
	}, nil
}

func (s *nodePrivilegedService) SendObservationRequest(ctx context.Context, req *nodev1.SendObservationRequestRequest) (*nodev1.SendObservationRequestResponse, error) {
	if err := common.PostObservationRequest(s.obsvReqSendC, req.ObservationRequest); err != nil {
		return nil, err
	}

	s.logger.Info("sent observation request", zap.Any("request", req.ObservationRequest))
	return &nodev1.SendObservationRequestResponse{}, nil
}

func (s *nodePrivilegedService) ChainGovernorStatus(ctx context.Context, req *nodev1.ChainGovernorStatusRequest) (*nodev1.ChainGovernorStatusResponse, error) {
	if s.governor == nil {
		return nil, fmt.Errorf("chain governor is not enabled")
	}

	return &nodev1.ChainGovernorStatusResponse{
		Response: s.governor.Status(),
	}, nil
}

func (s *nodePrivilegedService) ChainGovernorReload(ctx context.Context, req *nodev1.ChainGovernorReloadRequest) (*nodev1.ChainGovernorReloadResponse, error) {
	if s.governor == nil {
		return nil, fmt.Errorf("chain governor is not enabled")
	}

	resp, err := s.governor.Reload()
	if err != nil {
		return nil, err
	}

	return &nodev1.ChainGovernorReloadResponse{
		Response: resp,
	}, nil
}

func (s *nodePrivilegedService) ChainGovernorDropPendingVAA(ctx context.Context, req *nodev1.ChainGovernorDropPendingVAARequest) (*nodev1.ChainGovernorDropPendingVAAResponse, error) {
	if s.governor == nil {
		return nil, fmt.Errorf("chain governor is not enabled")
	}

	if len(req.VaaId) == 0 {
		return nil, fmt.Errorf("the VAA id must be specified as \"chainId/emitterAddress/seqNum\"")
	}

	resp, err := s.governor.DropPendingVAA(req.VaaId)
	if err != nil {
		return nil, err
	}

	return &nodev1.ChainGovernorDropPendingVAAResponse{
		Response: resp,
	}, nil
}

func (s *nodePrivilegedService) ChainGovernorReleasePendingVAA(ctx context.Context, req *nodev1.ChainGovernorReleasePendingVAARequest) (*nodev1.ChainGovernorReleasePendingVAAResponse, error) {
	if s.governor == nil {
		return nil, fmt.Errorf("chain governor is not enabled")
	}

	if len(req.VaaId) == 0 {
		return nil, fmt.Errorf("the VAA id must be specified as \"chainId/emitterAddress/seqNum\"")
	}

	resp, err := s.governor.ReleasePendingVAA(req.VaaId)
	if err != nil {
		return nil, err
	}

	return &nodev1.ChainGovernorReleasePendingVAAResponse{
		Response: resp,
	}, nil
}

func (s *nodePrivilegedService) ChainGovernorResetReleaseTimer(ctx context.Context, req *nodev1.ChainGovernorResetReleaseTimerRequest) (*nodev1.ChainGovernorResetReleaseTimerResponse, error) {
	if s.governor == nil {
		return nil, fmt.Errorf("chain governor is not enabled")
	}

	if len(req.VaaId) == 0 {
		return nil, fmt.Errorf("the VAA id must be specified as \"chainId/emitterAddress/seqNum\"")
	}

	resp, err := s.governor.ResetReleaseTimer(req.VaaId)
	if err != nil {
		return nil, err
	}

	return &nodev1.ChainGovernorResetReleaseTimerResponse{
		Response: resp,
	}, nil
}

func (s *nodePrivilegedService) PurgePythNetVaas(ctx context.Context, req *nodev1.PurgePythNetVaasRequest) (*nodev1.PurgePythNetVaasResponse, error) {
	prefix := db.VAAID{EmitterChain: vaa.ChainIDPythNet}
	oldestTime := time.Now().Add(-time.Hour * 24 * time.Duration(req.DaysOld))
	resp, err := s.db.PurgeVaas(prefix, oldestTime, req.LogOnly)
	if err != nil {
		return nil, err
	}

	return &nodev1.PurgePythNetVaasResponse{
		Response: resp,
	}, nil
}

func (s *nodePrivilegedService) SignExistingVAA(ctx context.Context, req *nodev1.SignExistingVAARequest) (*nodev1.SignExistingVAAResponse, error) {
	v, err := vaa.Unmarshal(req.Vaa)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal VAA: %w", err)
	}

	if req.NewPhylaxSetIndex <= v.PhylaxSetIndex {
		return nil, errors.New("new phylax set index must be higher than provided VAA")
	}

	if s.evmConnector == nil {
		return nil, errors.New("the node needs to have an Ethereum connection configured to sign existing VAAs")
	}

	var gs *common.PhylaxSet
	if cachedGs, exists := s.gsCache.Load(v.PhylaxSetIndex); exists {
		var ok bool
		gs, ok = cachedGs.(*common.PhylaxSet)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}
	} else {
		evmGs, err := s.evmConnector.GetPhylaxSet(ctx, v.PhylaxSetIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to load phylax set [%d]: %w", v.PhylaxSetIndex, err)
		}
		gs = &common.PhylaxSet{
			Keys:  evmGs.Keys,
			Index: v.PhylaxSetIndex,
		}
		s.gsCache.Store(v.PhylaxSetIndex, gs)
	}

	if slices.Index(gs.Keys, s.phylaxAddress) != -1 {
		return nil, fmt.Errorf("local phylax is already on the old set")
	}

	// Verify VAA
	err = v.Verify(gs.Keys)
	if err != nil {
		return nil, fmt.Errorf("failed to verify existing VAA: %w", err)
	}

	if len(req.NewPhylaxAddrs) > 255 {
		return nil, errors.New("new phylax set has too many phylaxs")
	}
	newGS := make([]ethcommon.Address, len(req.NewPhylaxAddrs))
	for i, phylaxString := range req.NewPhylaxAddrs {
		phylaxAddress := ethcommon.HexToAddress(phylaxString)
		newGS[i] = phylaxAddress
	}

	// Make sure there are no duplicates. Compact needs to take a sorted slice to remove all duplicates.
	newGSSorted := slices.Clone(newGS)
	slices.SortFunc(newGSSorted, func(a, b ethcommon.Address) bool {
		return bytes.Compare(a[:], b[:]) < 0
	})
	newGsLen := len(newGSSorted)
	if len(slices.Compact(newGSSorted)) != newGsLen {
		return nil, fmt.Errorf("duplicate phylaxs in the phylax set")
	}

	localPhylaxIndex := slices.Index(newGS, s.phylaxAddress)
	if localPhylaxIndex == -1 {
		return nil, fmt.Errorf("local phylax is not a member of the new phylax set")
	}

	newVAA := &vaa.VAA{
		Version: v.Version,
		// Set the new phylax set index
		PhylaxSetIndex: req.NewPhylaxSetIndex,
		// Signatures will be repopulated
		Signatures:       nil,
		Timestamp:        v.Timestamp,
		Nonce:            v.Nonce,
		Sequence:         v.Sequence,
		ConsistencyLevel: v.ConsistencyLevel,
		EmitterChain:     v.EmitterChain,
		EmitterAddress:   v.EmitterAddress,
		Payload:          v.Payload,
	}

	// Copy original VAA signatures
	for _, sig := range v.Signatures {
		signerAddress := gs.Keys[sig.Index]
		newIndex := slices.Index(newGS, signerAddress)
		// Phylax is not part of the new set
		if newIndex == -1 {
			continue
		}
		newVAA.Signatures = append(newVAA.Signatures, &vaa.Signature{
			Index:     uint8(newIndex),
			Signature: sig.Signature,
		})
	}

	// Add our own signature only if the new phylax set would reach quorum
	if vaa.CalculateQuorum(len(newGS)) > len(newVAA.Signatures)+1 {
		return nil, errors.New("cannot reach quorum on new phylax set with the local signature")
	}

	// Add local signature
	newVAA.AddSignature(s.gk, uint8(localPhylaxIndex))

	// Sort VAA signatures by phylax ID
	slices.SortFunc(newVAA.Signatures, func(a, b *vaa.Signature) bool {
		return a.Index < b.Index
	})

	newVAABytes, err := newVAA.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal new VAA: %w", err)
	}

	return &nodev1.SignExistingVAAResponse{Vaa: newVAABytes}, nil
}

func (s *nodePrivilegedService) DumpRPCs(ctx context.Context, req *nodev1.DumpRPCsRequest) (*nodev1.DumpRPCsResponse, error) {
	return &nodev1.DumpRPCsResponse{
		Response: s.rpcMap,
	}, nil
}
