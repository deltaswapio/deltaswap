//nolint:unparam
package adminrpc

import (
	"context"
	"crypto/ecdsa"
	"testing"
	"time"

	nodev1 "github.com/deltaswapio/deltaswap/node/pkg/proto/node/v1"
	"github.com/deltaswapio/deltaswap/node/pkg/watchers/evm/connectors"
	"github.com/deltaswapio/deltaswap/node/pkg/watchers/evm/connectors/ethabi"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockEVMConnector struct {
	phylaxAddrs    []common.Address
	phylaxSetIndex uint32
}

func (m mockEVMConnector) GetCurrentPhylaxSetIndex(ctx context.Context) (uint32, error) {
	return m.phylaxSetIndex, nil
}

func (m mockEVMConnector) GetPhylaxSet(ctx context.Context, index uint32) (ethabi.StructsPhylaxSet, error) {
	return ethabi.StructsPhylaxSet{
		Keys:           m.phylaxAddrs,
		ExpirationTime: 0,
	}, nil
}

func (m mockEVMConnector) NetworkName() string {
	panic("unimplemented")
}

func (m mockEVMConnector) ContractAddress() common.Address {
	panic("unimplemented")
}

func (m mockEVMConnector) WatchLogMessagePublished(ctx context.Context, errC chan error, sink chan<- *ethabi.AbiLogMessagePublished) (event.Subscription, error) {
	panic("unimplemented")
}

func (m mockEVMConnector) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	panic("unimplemented")
}

func (m mockEVMConnector) TimeOfBlockByHash(ctx context.Context, hash common.Hash) (uint64, error) {
	panic("unimplemented")
}

func (m mockEVMConnector) ParseLogMessagePublished(log types.Log) (*ethabi.AbiLogMessagePublished, error) {
	panic("unimplemented")
}

func (m mockEVMConnector) SubscribeForBlocks(ctx context.Context, errC chan error, sink chan<- *connectors.NewBlock) (ethereum.Subscription, error) {
	panic("unimplemented")
}

func (m mockEVMConnector) RawCallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	panic("unimplemented")
}

func (m mockEVMConnector) RawBatchCallContext(ctx context.Context, b []ethRpc.BatchElem) error {
	panic("unimplemented")
}

func generateGS(num int) (keys []*ecdsa.PrivateKey, addrs []common.Address) {
	for i := 0; i < num; i++ {
		key, err := ethcrypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		keys = append(keys, key)
		addrs = append(addrs, ethcrypto.PubkeyToAddress(key.PublicKey))
	}
	return
}

func addrsToHexStrings(addrs []common.Address) (out []string) {
	for _, addr := range addrs {
		out = append(out, addr.String())
	}
	return
}

func generateMockVAA(gsIndex uint32, gsKeys []*ecdsa.PrivateKey) []byte {
	v := &vaa.VAA{
		Version:          1,
		PhylaxSetIndex:   gsIndex,
		Signatures:       nil,
		Timestamp:        time.Now(),
		Nonce:            3,
		Sequence:         79,
		ConsistencyLevel: 1,
		EmitterChain:     1,
		EmitterAddress:   vaa.Address{},
		Payload:          []byte("test"),
	}
	for i, key := range gsKeys {
		v.AddSignature(key, uint8(i))
	}

	vBytes, err := v.Marshal()
	if err != nil {
		panic(err)
	}
	return vBytes
}

func setupAdminServerForVAASigning(gsIndex uint32, gsAddrs []common.Address) *nodePrivilegedService {
	gk, err := ethcrypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	connector := mockEVMConnector{
		phylaxAddrs:    gsAddrs,
		phylaxSetIndex: gsIndex,
	}

	return &nodePrivilegedService{
		db:            nil,
		injectC:       nil,
		obsvReqSendC:  nil,
		logger:        zap.L(),
		signedInC:     nil,
		governor:      nil,
		evmConnector:  connector,
		gk:            gk,
		phylaxAddress: ethcrypto.PubkeyToAddress(gk.PublicKey),
	}
}

func TestSignExistingVAA_NoVAA(t *testing.T) {
	s := setupAdminServerForVAASigning(0, []common.Address{})

	_, err := s.SignExistingVAA(context.Background(), &nodev1.SignExistingVAARequest{
		Vaa:               nil,
		NewPhylaxAddrs:    nil,
		NewPhylaxSetIndex: 0,
	})
	require.ErrorContains(t, err, "failed to unmarshal VAA")
}

func TestSignExistingVAA_NotPhylax(t *testing.T) {
	gsKeys, gsAddrs := generateGS(5)
	s := setupAdminServerForVAASigning(0, gsAddrs)

	v := generateMockVAA(0, gsKeys)

	_, err := s.SignExistingVAA(context.Background(), &nodev1.SignExistingVAARequest{
		Vaa:               v,
		NewPhylaxAddrs:    addrsToHexStrings(gsAddrs),
		NewPhylaxSetIndex: 1,
	})
	require.ErrorContains(t, err, "local phylax is not a member of the new phylax set")
}

func TestSignExistingVAA_InvalidVAA(t *testing.T) {
	gsKeys, gsAddrs := generateGS(5)
	s := setupAdminServerForVAASigning(0, gsAddrs)

	v := generateMockVAA(0, gsKeys[:2])

	gsAddrs = append(gsAddrs, s.phylaxAddress)
	_, err := s.SignExistingVAA(context.Background(), &nodev1.SignExistingVAARequest{
		Vaa:               v,
		NewPhylaxAddrs:    addrsToHexStrings(gsAddrs),
		NewPhylaxSetIndex: 1,
	})
	require.ErrorContains(t, err, "failed to verify existing VAA")
}

func TestSignExistingVAA_DuplicatePhylax(t *testing.T) {
	gsKeys, gsAddrs := generateGS(5)
	s := setupAdminServerForVAASigning(0, gsAddrs)

	v := generateMockVAA(0, gsKeys)

	gsAddrs = append(gsAddrs, s.phylaxAddress)
	gsAddrs = append(gsAddrs, s.phylaxAddress)
	_, err := s.SignExistingVAA(context.Background(), &nodev1.SignExistingVAARequest{
		Vaa:               v,
		NewPhylaxAddrs:    addrsToHexStrings(gsAddrs),
		NewPhylaxSetIndex: 1,
	})
	require.ErrorContains(t, err, "duplicate phylaxs in the phylax set")
}

func TestSignExistingVAA_AlreadyPhylax(t *testing.T) {
	gsKeys, gsAddrs := generateGS(5)
	s := setupAdminServerForVAASigning(0, gsAddrs)
	s.evmConnector = mockEVMConnector{
		phylaxAddrs:    append(gsAddrs, s.phylaxAddress),
		phylaxSetIndex: 0,
	}

	v := generateMockVAA(0, append(gsKeys, s.gk))

	gsAddrs = append(gsAddrs, s.phylaxAddress)
	_, err := s.SignExistingVAA(context.Background(), &nodev1.SignExistingVAARequest{
		Vaa:               v,
		NewPhylaxAddrs:    addrsToHexStrings(gsAddrs),
		NewPhylaxSetIndex: 1,
	})
	require.ErrorContains(t, err, "local phylax is already on the old set")
}

func TestSignExistingVAA_NotAFuturePhylax(t *testing.T) {
	gsKeys, gsAddrs := generateGS(5)
	s := setupAdminServerForVAASigning(0, gsAddrs)

	v := generateMockVAA(0, gsKeys)

	_, err := s.SignExistingVAA(context.Background(), &nodev1.SignExistingVAARequest{
		Vaa:               v,
		NewPhylaxAddrs:    addrsToHexStrings(gsAddrs),
		NewPhylaxSetIndex: 1,
	})
	require.ErrorContains(t, err, "local phylax is not a member of the new phylax set")
}

func TestSignExistingVAA_CantReachQuorum(t *testing.T) {
	gsKeys, gsAddrs := generateGS(5)
	s := setupAdminServerForVAASigning(0, gsAddrs)

	v := generateMockVAA(0, gsKeys)

	gsAddrs = append(gsAddrs, s.phylaxAddress)
	_, err := s.SignExistingVAA(context.Background(), &nodev1.SignExistingVAARequest{
		Vaa:               v,
		NewPhylaxAddrs:    addrsToHexStrings(append(gsAddrs, common.Address{0, 1}, common.Address{3, 1}, common.Address{8, 1})),
		NewPhylaxSetIndex: 1,
	})
	require.ErrorContains(t, err, "cannot reach quorum on new phylax set with the local signature")
}

func TestSignExistingVAA_Valid(t *testing.T) {
	gsKeys, gsAddrs := generateGS(5)
	s := setupAdminServerForVAASigning(0, gsAddrs)

	v := generateMockVAA(0, gsKeys)

	gsAddrs = append(gsAddrs, s.phylaxAddress)
	res, err := s.SignExistingVAA(context.Background(), &nodev1.SignExistingVAARequest{
		Vaa:               v,
		NewPhylaxAddrs:    addrsToHexStrings(gsAddrs),
		NewPhylaxSetIndex: 1,
	})

	require.NoError(t, err)
	v2 := generateMockVAA(1, append(gsKeys, s.gk))
	require.Equal(t, v2, res.Vaa)
}
