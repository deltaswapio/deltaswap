package phylaxd

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
	"github.com/spf13/pflag"
	"github.com/tendermint/tendermint/libs/rand"
	"github.com/wormhole-foundation/wormhole/sdk/vaa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"

	"github.com/certusone/wormhole/node/pkg/devnet"
	nodev1 "github.com/certusone/wormhole/node/pkg/proto/node/v1"
)

var setUpdateNumGuardians *int
var templateGuardianIndex *int
var chainID *string
var address *string
var module *string

var circleIntegrationChainID *string
var circleIntegrationFinality *string
var circleIntegrationForeignEmitterChainID *string
var circleIntegrationForeignEmitterAddress *string
var circleIntegrationCircleDomain *string
var circleIntegrationNewImplementationAddress *string

var deltachainStoreCodeWasmHash *string

var deltachainInstantiateContractCodeId *string
var deltachainInstantiateContractInstantiationMsg *string
var deltachainInstantiateContractLabel *string

var deltachainMigrateContractCodeId *string
var deltachainMigrateContractContractAddress *string
var deltachainMigrateContractInstantiationMsg *string

var deltachainWasmInstantiateAllowlistCodeId *string
var deltachainWasmInstantiateAllowlistContractAddress *string

var gatewayScheduleUpgradeName *string
var gatewayScheduleUpgradeHeight *string
var gatewayIbcComposabilityMwContractAddress *string

var ibcUpdateChannelChainTargetChainId *string
var ibcUpdateChannelChainChannelId *string
var ibcUpdateChannelChainChainId *string

func init() {
	governanceFlagSet := pflag.NewFlagSet("governance", pflag.ExitOnError)
	chainID = governanceFlagSet.String("chain-id", "", "Chain ID")
	address = governanceFlagSet.String("new-address", "", "New address (hex, base58 or bech32)")

	moduleFlagSet := pflag.NewFlagSet("module", pflag.ExitOnError)
	module = moduleFlagSet.String("module", "", "Module name")

	templateGuardianIndex = TemplateCmd.PersistentFlags().Int("idx", 3, "Default current guardian set index")

	setUpdateNumGuardians = AdminClientGuardianSetTemplateCmd.Flags().Int("num", 1, "Number of devnet guardians in example file")
	TemplateCmd.AddCommand(AdminClientGuardianSetTemplateCmd)

	AdminClientContractUpgradeTemplateCmd.Flags().AddFlagSet(governanceFlagSet)
	TemplateCmd.AddCommand(AdminClientContractUpgradeTemplateCmd)

	AdminClientTokenBridgeRegisterChainCmd.Flags().AddFlagSet(governanceFlagSet)
	AdminClientTokenBridgeRegisterChainCmd.Flags().AddFlagSet(moduleFlagSet)
	TemplateCmd.AddCommand(AdminClientTokenBridgeRegisterChainCmd)

	AdminClientTokenBridgeUpgradeContractCmd.Flags().AddFlagSet(governanceFlagSet)
	AdminClientTokenBridgeUpgradeContractCmd.Flags().AddFlagSet(moduleFlagSet)
	TemplateCmd.AddCommand(AdminClientTokenBridgeUpgradeContractCmd)

	AdminClientWormholeRelayerSetDefaultDeliveryProviderCmd.Flags().AddFlagSet(governanceFlagSet)
	TemplateCmd.AddCommand(AdminClientWormholeRelayerSetDefaultDeliveryProviderCmd)

	circleIntegrationChainIDFlagSet := pflag.NewFlagSet("circle-integ", pflag.ExitOnError)
	circleIntegrationChainID = circleIntegrationChainIDFlagSet.String("chain-id", "", "Target chain ID")

	circleIntegrationFinalityFlagSet := pflag.NewFlagSet("finality", pflag.ExitOnError)
	circleIntegrationFinality = circleIntegrationFinalityFlagSet.String("finality", "", "Desired wormhole finality")
	AdminClientCircleIntegrationUpdateWormholeFinalityCmd.Flags().AddFlagSet(circleIntegrationChainIDFlagSet)
	AdminClientCircleIntegrationUpdateWormholeFinalityCmd.Flags().AddFlagSet(circleIntegrationFinalityFlagSet)
	TemplateCmd.AddCommand(AdminClientCircleIntegrationUpdateWormholeFinalityCmd)

	circleIntegrationRegisterEmitterFlagSet := pflag.NewFlagSet("register", pflag.ExitOnError)
	circleIntegrationForeignEmitterChainID = circleIntegrationRegisterEmitterFlagSet.String("foreign-emitter-chain-id", "", "Foreign emitter chain ID")
	circleIntegrationForeignEmitterAddress = circleIntegrationRegisterEmitterFlagSet.String("foreign-emitter-address", "", "Foreign emitter address (hex, base58 or bech32)")
	circleIntegrationCircleDomain = circleIntegrationRegisterEmitterFlagSet.String("circle-domain", "", "Circle domain")
	AdminClientCircleIntegrationRegisterEmitterAndDomainCmd.Flags().AddFlagSet(circleIntegrationChainIDFlagSet)
	AdminClientCircleIntegrationRegisterEmitterAndDomainCmd.Flags().AddFlagSet(circleIntegrationRegisterEmitterFlagSet)
	TemplateCmd.AddCommand(AdminClientCircleIntegrationRegisterEmitterAndDomainCmd)

	circleIntegrationUpgradeContractImplementationFlagSet := pflag.NewFlagSet("upgrade", pflag.ExitOnError)
	circleIntegrationNewImplementationAddress = circleIntegrationUpgradeContractImplementationFlagSet.String("new-implementation-address", "", "New implementation address (hex, base58 or bech32)")
	AdminClientCircleIntegrationUpgradeContractImplementationCmd.Flags().AddFlagSet(circleIntegrationChainIDFlagSet)
	AdminClientCircleIntegrationUpgradeContractImplementationCmd.Flags().AddFlagSet(circleIntegrationUpgradeContractImplementationFlagSet)
	TemplateCmd.AddCommand(AdminClientCircleIntegrationUpgradeContractImplementationCmd)

	deltachainStoreCodeFlagSet := pflag.NewFlagSet("deltachain-store-code", pflag.ExitOnError)
	deltachainStoreCodeWasmHash = deltachainStoreCodeFlagSet.String("wasm-hash", "", "WASM Hash of the stored code")
	AdminClientWormchainStoreCodeCmd.Flags().AddFlagSet(deltachainStoreCodeFlagSet)
	TemplateCmd.AddCommand(AdminClientWormchainStoreCodeCmd)

	deltachainInstantiateContractFlagSet := pflag.NewFlagSet("deltachain-instantiate-contract", pflag.ExitOnError)
	deltachainInstantiateContractCodeId = deltachainInstantiateContractFlagSet.String("code-id", "", "code ID of the stored code")
	deltachainInstantiateContractLabel = deltachainInstantiateContractFlagSet.String("label", "", "label")
	deltachainInstantiateContractInstantiationMsg = deltachainInstantiateContractFlagSet.String("instantiation-msg", "", "instantiate message")
	AdminClientWormchainInstantiateContractCmd.Flags().AddFlagSet(deltachainInstantiateContractFlagSet)
	TemplateCmd.AddCommand(AdminClientWormchainInstantiateContractCmd)

	deltachainMigrateContractFlagSet := pflag.NewFlagSet("deltachain-migrate-contract", pflag.ExitOnError)
	deltachainMigrateContractCodeId = deltachainMigrateContractFlagSet.String("code-id", "", "code ID of the stored code")
	deltachainMigrateContractContractAddress = deltachainMigrateContractFlagSet.String("contract-address", "", "contract address")
	deltachainMigrateContractInstantiationMsg = deltachainMigrateContractFlagSet.String("instantiation-msg", "", "instantiate message")
	AdminClientWormchainMigrateContractCmd.Flags().AddFlagSet(deltachainMigrateContractFlagSet)
	TemplateCmd.AddCommand(AdminClientWormchainMigrateContractCmd)

	// flags for the deltachain add/delete wasm instantiate allowlist commands
	deltachainWasmInstantiateAllowlistFlagSet := pflag.NewFlagSet("deltachain-wasm-instantiate-allowlist", pflag.ExitOnError)
	deltachainWasmInstantiateAllowlistCodeId = deltachainWasmInstantiateAllowlistFlagSet.String("code-id", "", "code ID of the stored code to add/delete allowlist wasm instantiate for")
	deltachainWasmInstantiateAllowlistContractAddress = deltachainWasmInstantiateAllowlistFlagSet.String("contract-address", "", "contract address to add/delete allowlist wasm instantiate for")
	AdminClientWormchainAddWasmInstantiateAllowlistCmd.Flags().AddFlagSet(deltachainWasmInstantiateAllowlistFlagSet)
	AdminClientWormchainDeleteWasmInstantiateAllowlistCmd.Flags().AddFlagSet(deltachainWasmInstantiateAllowlistFlagSet)
	TemplateCmd.AddCommand(AdminClientWormchainAddWasmInstantiateAllowlistCmd)
	TemplateCmd.AddCommand(AdminClientWormchainDeleteWasmInstantiateAllowlistCmd)

	// flags for the gateway-ibc-composability-mw-set-contract command
	gatewayIbcComposabilityMwFlagSet := pflag.NewFlagSet("gateway-ibc-composability-mw-set-contract", pflag.ExitOnError)
	gatewayIbcComposabilityMwContractAddress = gatewayIbcComposabilityMwFlagSet.String("contract-address", "", "contract address to set in the ibc composability middleware")
	AdminClientGatewayIbcComposabilityMwSetContractCmd.Flags().AddFlagSet(gatewayIbcComposabilityMwFlagSet)
	TemplateCmd.AddCommand(AdminClientGatewayIbcComposabilityMwSetContractCmd)

	// flags for the gateway-schedule-upgrade command
	gatewayScheduleUpgradeFlagSet := pflag.NewFlagSet("gateway-schedule-upgrade", pflag.ExitOnError)
	gatewayScheduleUpgradeName = gatewayScheduleUpgradeFlagSet.String("name", "", "Scheduled upgrade name")
	gatewayScheduleUpgradeHeight = gatewayScheduleUpgradeFlagSet.String("height", "", "Scheduled upgrade height")
	AdminClientGatewayScheduleUpgradeCmd.Flags().AddFlagSet(gatewayScheduleUpgradeFlagSet)
	TemplateCmd.AddCommand(AdminClientGatewayScheduleUpgradeCmd)

	// AdminClientGatewayCancelUpgradeCmd doesn't have any flags
	TemplateCmd.AddCommand(AdminClientGatewayCancelUpgradeCmd)

	// flags for the ibc-receiver-update-channel-chain and ibc-translator-update-channel-chain commands
	ibcUpdateChannelChainFlagSet := pflag.NewFlagSet("ibc-mapping", pflag.ExitOnError)
	ibcUpdateChannelChainTargetChainId = ibcUpdateChannelChainFlagSet.String("target-chain-id", "", "Target Chain ID for the governance VAA")
	ibcUpdateChannelChainChannelId = ibcUpdateChannelChainFlagSet.String("channel-id", "", "IBC Channel ID on Wormchain")
	ibcUpdateChannelChainChainId = ibcUpdateChannelChainFlagSet.String("chain-id", "", "IBC Chain ID that the channel ID corresponds to")
	AdminClientIbcReceiverUpdateChannelChainCmd.Flags().AddFlagSet(ibcUpdateChannelChainFlagSet)
	AdminClientIbcTranslatorUpdateChannelChainCmd.Flags().AddFlagSet(ibcUpdateChannelChainFlagSet)
	TemplateCmd.AddCommand(AdminClientIbcReceiverUpdateChannelChainCmd)
	TemplateCmd.AddCommand(AdminClientIbcTranslatorUpdateChannelChainCmd)
}

var TemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Guardian governance VAA template commands ",
}

var AdminClientGuardianSetTemplateCmd = &cobra.Command{
	Use:   "guardian-set-update",
	Short: "Generate an empty guardian set template",
	Run:   runGuardianSetTemplate,
}

var AdminClientContractUpgradeTemplateCmd = &cobra.Command{
	Use:   "contract-upgrade",
	Short: "Generate an empty contract upgrade template",
	Run:   runContractUpgradeTemplate,
}

var AdminClientTokenBridgeRegisterChainCmd = &cobra.Command{
	Use:   "token-bridge-register-chain",
	Short: "Generate an empty token bridge chain registration template at specified path",
	Run:   runTokenBridgeRegisterChainTemplate,
}

var AdminClientTokenBridgeUpgradeContractCmd = &cobra.Command{
	Use:   "token-bridge-upgrade-contract",
	Short: "Generate an empty token bridge contract upgrade template at specified path",
	Run:   runTokenBridgeUpgradeContractTemplate,
}

var AdminClientCircleIntegrationUpdateWormholeFinalityCmd = &cobra.Command{
	Use:   "circle-integration-update-wormhole-finality",
	Short: "Generate an empty circle integration update wormhole finality template at specified path",
	Run:   runCircleIntegrationUpdateWormholeFinalityTemplate,
}

var AdminClientCircleIntegrationRegisterEmitterAndDomainCmd = &cobra.Command{
	Use:   "circle-integration-register-emitter-and-domain",
	Short: "Generate an empty circle integration register emitter and domain template at specified path",
	Run:   runCircleIntegrationRegisterEmitterAndDomainTemplate,
}

var AdminClientCircleIntegrationUpgradeContractImplementationCmd = &cobra.Command{
	Use:   "circle-integration-upgrade-contract-implementation",
	Short: "Generate an empty circle integration upgrade contract implementation template at specified path",
	Run:   runCircleIntegrationUpgradeContractImplementationTemplate,
}

var AdminClientWormchainStoreCodeCmd = &cobra.Command{
	Use:   "deltachain-store-code",
	Short: "Generate an empty deltachain store code template at specified path",
	Run:   runWormchainStoreCodeTemplate,
}

var AdminClientWormchainInstantiateContractCmd = &cobra.Command{
	Use:   "deltachain-instantiate-contract",
	Short: "Generate an empty deltachain instantiate contract template at specified path",
	Run:   runWormchainInstantiateContractTemplate,
}

var AdminClientWormchainMigrateContractCmd = &cobra.Command{
	Use:   "deltachain-migrate-contract",
	Short: "Generate an empty deltachain migrate contract template at specified path",
	Run:   runWormchainMigrateContractTemplate,
}

var AdminClientWormchainAddWasmInstantiateAllowlistCmd = &cobra.Command{
	Use:   "deltachain-add-wasm-instantiate-allowlist",
	Short: "Generate an empty deltachain add wasm instantiate allowlist template at specified path",
	Run:   runWormchainAddWasmInstantiateAllowlistTemplate,
}

var AdminClientWormchainDeleteWasmInstantiateAllowlistCmd = &cobra.Command{
	Use:   "deltachain-delete-wasm-instantiate-allowlist",
	Short: "Generate an empty deltachain delete wasm instantiate allowlist template at specified path",
	Run:   runWormchainDeleteWasmInstantiateAllowlistTemplate,
}

var AdminClientGatewayScheduleUpgradeCmd = &cobra.Command{
	Use:   "gateway-schedule-upgrade",
	Short: "Schedule an upgrade on Gateway with a specified name for a specified height",
	Run:   runGatewayScheduleUpgradeTemplate,
}

var AdminClientGatewayCancelUpgradeCmd = &cobra.Command{
	Use:   "gateway-cancel-upgrade",
	Short: "Cancel a scheduled upgrade on Gateway",
	Run:   runGatewayCancelUpgradeTemplate,
}

var AdminClientGatewayIbcComposabilityMwSetContractCmd = &cobra.Command{
	Use:   "gateway-ibc-composability-mw-set-contract",
	Short: "Set the contract that the IBC Composability middleware will query",
	Run:   runGatewayIbcComposabilityMwSetContractTemplate,
}

var AdminClientIbcReceiverUpdateChannelChainCmd = &cobra.Command{
	Use:   "ibc-receiver-update-channel-chain",
	Short: "Generate an empty ibc receiver channelId to chainId mapping update template at specified path",
	Run:   runIbcReceiverUpdateChannelChainTemplate,
}

var AdminClientIbcTranslatorUpdateChannelChainCmd = &cobra.Command{
	Use:   "ibc-translator-update-channel-chain",
	Short: "Generate an empty ibc translator channelId to chainId mapping update template at specified path",
	Run:   runIbcTranslatorUpdateChannelChainTemplate,
}

var AdminClientWormholeRelayerSetDefaultDeliveryProviderCmd = &cobra.Command{
	Use:   "wormhole-relayer-set-default-delivery-provider",
	Short: "Generate a 'set default delivery provider' template for specified chain and address",
	Run:   runWormholeRelayerSetDefaultDeliveryProviderTemplate,
}

func runGuardianSetTemplate(cmd *cobra.Command, args []string) {
	// Use deterministic devnet addresses as examples in the template, such that this doubles as a test fixture.
	guardians := make([]*nodev1.GuardianSetUpdate_Guardian, *setUpdateNumGuardians)
	for i := 0; i < *setUpdateNumGuardians; i++ {
		k := devnet.InsecureDeterministicEcdsaKeyByIndex(crypto.S256(), uint64(i))
		guardians[i] = &nodev1.GuardianSetUpdate_Guardian{
			Pubkey: crypto.PubkeyToAddress(k.PublicKey).Hex(),
			Name:   fmt.Sprintf("Example validator %d", i),
		}
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_GuardianSet{
					GuardianSet: &nodev1.GuardianSetUpdate{Guardians: guardians},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runContractUpgradeTemplate(cmd *cobra.Command, args []string) {
	address, err := parseAddress(*address)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := parseChainID(*chainID)
	if err != nil {
		log.Fatal(err)
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_ContractUpgrade{
					ContractUpgrade: &nodev1.ContractUpgrade{
						ChainId:     uint32(chainID),
						NewContract: address,
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}
func runTokenBridgeRegisterChainTemplate(cmd *cobra.Command, args []string) {
	address, err := parseAddress(*address)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := parseChainID(*chainID)
	if err != nil {
		log.Fatal(err)
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_BridgeRegisterChain{
					BridgeRegisterChain: &nodev1.BridgeRegisterChain{
						Module:         *module,
						ChainId:        uint32(chainID),
						EmitterAddress: address,
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runTokenBridgeUpgradeContractTemplate(cmd *cobra.Command, args []string) {
	address, err := parseAddress(*address)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := parseChainID(*chainID)
	if err != nil {
		log.Fatal(err)
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_BridgeContractUpgrade{
					BridgeContractUpgrade: &nodev1.BridgeUpgradeContract{
						Module:        *module,
						TargetChainId: uint32(chainID),
						NewContract:   address,
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runCircleIntegrationUpdateWormholeFinalityTemplate(cmd *cobra.Command, args []string) {
	if *circleIntegrationChainID == "" {
		log.Fatal("--chain-id must be specified.")
	}
	chainID, err := parseChainID(*circleIntegrationChainID)
	if err != nil {
		log.Fatal("failed to parse chain id:", err)
	}
	if *circleIntegrationFinality == "" {
		log.Fatal("--finality must be specified.")
	}
	finality, err := strconv.ParseUint(*circleIntegrationFinality, 10, 8)
	if err != nil {
		log.Fatal("failed to parse finality as uint8: ", err)
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_CircleIntegrationUpdateWormholeFinality{
					CircleIntegrationUpdateWormholeFinality: &nodev1.CircleIntegrationUpdateWormholeFinality{
						TargetChainId: uint32(chainID),
						Finality:      uint32(finality),
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runCircleIntegrationRegisterEmitterAndDomainTemplate(cmd *cobra.Command, args []string) {
	if *circleIntegrationChainID == "" {
		log.Fatal("--chain-id must be specified.")
	}
	chainID, err := parseChainID(*circleIntegrationChainID)
	if err != nil {
		log.Fatal("failed to parse chain id:", err)
	}
	if *circleIntegrationForeignEmitterChainID == "" {
		log.Fatal("--foreign-emitter-chain-id must be specified.")
	}
	foreignEmitterChainId, err := parseChainID(*circleIntegrationForeignEmitterChainID)
	if err != nil {
		log.Fatal("failed to parse foreign emitter chain id as uint8:", err)
	}
	if *circleIntegrationForeignEmitterAddress == "" {
		log.Fatal("--foreign-emitter-address must be specified.")
	}
	foreignEmitterAddress, err := parseAddress(*circleIntegrationForeignEmitterAddress)
	if err != nil {
		log.Fatal("failed to parse foreign emitter address: ", err)
	}
	if *circleIntegrationCircleDomain == "" {
		log.Fatal("--circle-domain must be specified.")
	}
	circleDomain, err := strconv.ParseUint(*circleIntegrationCircleDomain, 10, 32)
	if err != nil {
		log.Fatal("failed to parse circle domain as uint32:", err)
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_CircleIntegrationRegisterEmitterAndDomain{
					CircleIntegrationRegisterEmitterAndDomain: &nodev1.CircleIntegrationRegisterEmitterAndDomain{
						TargetChainId:         uint32(chainID),
						ForeignEmitterChainId: uint32(foreignEmitterChainId),
						ForeignEmitterAddress: foreignEmitterAddress,
						CircleDomain:          uint32(circleDomain),
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runCircleIntegrationUpgradeContractImplementationTemplate(cmd *cobra.Command, args []string) {
	if *circleIntegrationChainID == "" {
		log.Fatal("--chain-id must be specified.")
	}
	chainID, err := parseChainID(*circleIntegrationChainID)
	if err != nil {
		log.Fatal("failed to parse chain id:", err)
	}
	if *circleIntegrationNewImplementationAddress == "" {
		log.Fatal("--new-implementation-address must be specified.")
	}
	newImplementationAddress, err := parseAddress(*circleIntegrationNewImplementationAddress)
	if err != nil {
		log.Fatal("failed to parse new implementation address: ", err)
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_CircleIntegrationUpgradeContractImplementation{
					CircleIntegrationUpgradeContractImplementation: &nodev1.CircleIntegrationUpgradeContractImplementation{
						TargetChainId:            uint32(chainID),
						NewImplementationAddress: newImplementationAddress,
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runWormchainStoreCodeTemplate(cmd *cobra.Command, args []string) {
	if *deltachainStoreCodeWasmHash == "" {
		log.Fatal("--wasm-hash must be specified.")
	}

	// Validate the string is valid hex.
	buf, err := hex.DecodeString(*deltachainStoreCodeWasmHash)
	if err != nil {
		log.Fatal("invalid wasm-hash (expected hex): %w", err)
	}

	// Validate the string is the correct length.
	if len(buf) != 32 {
		log.Fatalf("wasm-hash (expected 32 bytes but received %d bytes)", len(buf))
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_WormchainStoreCode{
					WormchainStoreCode: &nodev1.WormchainStoreCode{
						WasmHash: string(*deltachainStoreCodeWasmHash),
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runWormchainInstantiateContractTemplate(cmd *cobra.Command, args []string) {
	if *deltachainInstantiateContractCodeId == "" {
		log.Fatal("--code-id must be specified.")
	}
	codeId, err := strconv.ParseUint(*deltachainInstantiateContractCodeId, 10, 64)
	if err != nil {
		log.Fatal("failed to parse code-id as uint64: ", err)
	}
	if *deltachainInstantiateContractLabel == "" {
		log.Fatal("--label must be specified.")
	}
	if *deltachainInstantiateContractInstantiationMsg == "" {
		log.Fatal("--instantiation-msg must be specified.")
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_WormchainInstantiateContract{
					WormchainInstantiateContract: &nodev1.WormchainInstantiateContract{
						CodeId:           codeId,
						Label:            *deltachainInstantiateContractLabel,
						InstantiationMsg: *deltachainInstantiateContractInstantiationMsg,
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runWormchainMigrateContractTemplate(cmd *cobra.Command, args []string) {
	if *deltachainMigrateContractCodeId == "" {
		log.Fatal("--code-id must be specified.")
	}
	codeId, err := strconv.ParseUint(*deltachainMigrateContractCodeId, 10, 64)
	if err != nil {
		log.Fatal("failed to parse code-id as uint64: ", err)
	}
	if *deltachainMigrateContractContractAddress == "" {
		log.Fatal("--contract-address must be specified.")
	}
	if *deltachainMigrateContractInstantiationMsg == "" {
		log.Fatal("--instantiate-msg must be specified.")
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_WormchainMigrateContract{
					WormchainMigrateContract: &nodev1.WormchainMigrateContract{
						CodeId:           codeId,
						Contract:         *deltachainMigrateContractContractAddress,
						InstantiationMsg: *deltachainMigrateContractInstantiationMsg,
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runWormchainAddWasmInstantiateAllowlistTemplate(cmd *cobra.Command, args []string) {
	runWormchainWasmInstantiateAllowlistTemplate(nodev1.WormchainWasmInstantiateAllowlistAction_WORMCHAIN_WASM_INSTANTIATE_ALLOWLIST_ACTION_ADD)
}

func runWormchainDeleteWasmInstantiateAllowlistTemplate(cmd *cobra.Command, args []string) {
	runWormchainWasmInstantiateAllowlistTemplate(nodev1.WormchainWasmInstantiateAllowlistAction_WORMCHAIN_WASM_INSTANTIATE_ALLOWLIST_ACTION_DELETE)
}

func runWormchainWasmInstantiateAllowlistTemplate(action nodev1.WormchainWasmInstantiateAllowlistAction) {
	if *deltachainWasmInstantiateAllowlistCodeId == "" {
		log.Fatal("--code-id must be specified")
	}
	codeId, err := strconv.ParseUint(*deltachainWasmInstantiateAllowlistCodeId, 10, 64)
	if err != nil {
		log.Fatal("failed to parse code-id as utin64: ", err)
	}
	if *deltachainWasmInstantiateAllowlistContractAddress == "" {
		log.Fatal("--contract-address must be specified")
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_WormchainWasmInstantiateAllowlist{
					WormchainWasmInstantiateAllowlist: &nodev1.WormchainWasmInstantiateAllowlist{
						CodeId:   codeId,
						Contract: *deltachainWasmInstantiateAllowlistContractAddress,
						Action:   action,
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runGatewayScheduleUpgradeTemplate(cmd *cobra.Command, args []string) {
	if *gatewayScheduleUpgradeName == "" {
		log.Fatal("--name must be specified")
	}

	if *gatewayScheduleUpgradeHeight == "" {
		log.Fatal("--height must be specified")
	}

	height, err := strconv.ParseUint(*gatewayScheduleUpgradeHeight, 10, 64)
	if err != nil {
		log.Fatal("failed to parse height as uint64: ", err)
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_GatewayScheduleUpgrade{
					GatewayScheduleUpgrade: &nodev1.GatewayScheduleUpgrade{
						Name:   *gatewayScheduleUpgradeName,
						Height: height,
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runGatewayCancelUpgradeTemplate(cmd *cobra.Command, args []string) {
	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload:  &nodev1.GovernanceMessage_GatewayCancelUpgrade{},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runGatewayIbcComposabilityMwSetContractTemplate(cmd *cobra.Command, args []string) {
	if *gatewayIbcComposabilityMwContractAddress == "" {
		log.Fatal("--contract-address must be specified")
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_GatewayIbcComposabilityMwSetContract{
					GatewayIbcComposabilityMwSetContract: &nodev1.GatewayIbcComposabilityMwSetContract{
						Contract: *gatewayIbcComposabilityMwContractAddress,
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

func runIbcReceiverUpdateChannelChainTemplate(cmd *cobra.Command, args []string) {
	runIbcUpdateChannelChainTemplate(nodev1.IbcUpdateChannelChainModule_IBC_UPDATE_CHANNEL_CHAIN_MODULE_RECEIVER)
}

func runIbcTranslatorUpdateChannelChainTemplate(cmd *cobra.Command, args []string) {
	runIbcUpdateChannelChainTemplate(nodev1.IbcUpdateChannelChainModule_IBC_UPDATE_CHANNEL_CHAIN_MODULE_TRANSLATOR)
}

func runIbcUpdateChannelChainTemplate(module nodev1.IbcUpdateChannelChainModule) {
	if *ibcUpdateChannelChainTargetChainId == "" {
		log.Fatal("--target-chain-id must be specified")
	}
	targetChainId, err := parseChainID(*ibcUpdateChannelChainTargetChainId)
	if err != nil {
		log.Fatal("failed to parse chain id: ", err)
	}

	if *ibcUpdateChannelChainChannelId == "" {
		log.Fatal("--channel-id must be specified")
	}
	if len(*ibcUpdateChannelChainChannelId) > 64 {
		log.Fatal("invalid channel id length, must be <= 64")
	}

	if *ibcUpdateChannelChainChainId == "" {
		log.Fatal("--chain-id must be specified")
	}
	chainId, err := parseChainID(*ibcUpdateChannelChainChainId)
	if err != nil {
		log.Fatal("failed to parse chain id: ", err)
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_IbcUpdateChannelChain{
					IbcUpdateChannelChain: &nodev1.IbcUpdateChannelChain{
						TargetChainId: uint32(targetChainId),
						ChannelId:     *ibcUpdateChannelChainChannelId,
						ChainId:       uint32(chainId),
						Module:        module,
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))

}

func runWormholeRelayerSetDefaultDeliveryProviderTemplate(cmd *cobra.Command, args []string) {
	address, err := parseAddress(*address)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := parseChainID(*chainID)
	if err != nil {
		log.Fatal(err)
	}

	m := &nodev1.InjectGovernanceVAARequest{
		CurrentSetIndex: uint32(*templateGuardianIndex),
		Messages: []*nodev1.GovernanceMessage{
			{
				Sequence: rand.Uint64(),
				Nonce:    rand.Uint32(),
				Payload: &nodev1.GovernanceMessage_WormholeRelayerSetDefaultDeliveryProvider{
					WormholeRelayerSetDefaultDeliveryProvider: &nodev1.WormholeRelayerSetDefaultDeliveryProvider{
						ChainId:                           uint32(chainID),
						NewDefaultDeliveryProviderAddress: address,
					},
				},
			},
		},
	}

	b, err := prototext.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

// parseAddress parses either a hex-encoded address and returns
// a left-padded 32 byte hex string.
func parseAddress(s string) (string, error) {
	// try base58
	b, err := base58.Decode(s)
	if err == nil {
		return leftPadAddress(b)
	}

	// try bech32
	_, b, err = bech32.Decode(s)
	if err == nil {
		return leftPadAddress(b)
	}

	// try hex
	if len(s) > 2 && strings.ToLower(s[:2]) == "0x" {
		s = s[2:]
	}

	a, err := hex.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("invalid hex address: %v", err)
	}
	return leftPadAddress(a)
}

func leftPadAddress(a []byte) (string, error) {
	if len(a) > 32 {
		return "", fmt.Errorf("address longer than 32 bytes")
	}
	return hex.EncodeToString(common.LeftPadBytes(a, 32)), nil
}

// parseChainID parses a human-readable chain name or a chain ID.
func parseChainID(name string) (vaa.ChainID, error) {
	s, err := vaa.ChainIDFromString(name)
	if err == nil {
		return s, nil
	}

	// parse as uint32
	i, err := strconv.ParseUint(name, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse as name or uint32: %v", err)
	}

	return vaa.ChainID(i), nil
}
