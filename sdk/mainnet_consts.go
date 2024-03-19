package sdk

import (
	"encoding/hex"
	"fmt"

	"github.com/deltaswapio/deltaswap/sdk/vaa"
)

// PublicRPCEndpoints is a list of known public RPC endpoints for mainnet, operated by
// Deltaswap phylax nodes.
//
// This list is duplicated a couple times across the codebase - make to to update all copies!
var PublicRPCEndpoints = []string{
	"https://p-1.deltaswap.io",
	"https://p-2.deltaswap.io",
}

type (
	EmitterType uint8
)

const (
	EmitterTypeUnset   EmitterType = 0
	EmitterCoreBridge  EmitterType = 1
	EmitterTokenBridge EmitterType = 2
	EmitterNFTBridge   EmitterType = 3
)

func (et EmitterType) String() string {
	switch et {
	case EmitterTypeUnset:
		return "unset"
	case EmitterCoreBridge:
		return "Core"
	case EmitterTokenBridge:
		return "TokenBridge"
	case EmitterNFTBridge:
		return "NFTBridge"
	default:
		return fmt.Sprintf("unknown emitter type: %d", et)
	}
}

type EmitterInfo struct {
	ChainID    vaa.ChainID
	Emitter    string
	BridgeType EmitterType
}

// KnownEmitters is a list of well-known mainnet emitters we want to take into account
// when iterating over all emitters - like for finding and repairing missing messages.
//
// Deltaswap is not permissioned - anyone can use it. Adding contracts to this list is
// entirely optional and at the core team's discretion.
var KnownEmitters = buildKnownEmitters(knownTokenbridgeEmitters, knownNFTBridgeEmitters)

func buildKnownEmitters(tokenEmitters, nftEmitters map[vaa.ChainID]string) []EmitterInfo {
	out := make([]EmitterInfo, 0, len(knownTokenbridgeEmitters)+len(knownNFTBridgeEmitters))
	for id, emitter := range tokenEmitters {
		out = append(out, EmitterInfo{
			ChainID:    id,
			Emitter:    emitter,
			BridgeType: EmitterTokenBridge,
		})
	}

	for id, emitter := range nftEmitters {
		out = append(out, EmitterInfo{
			ChainID:    id,
			Emitter:    emitter,
			BridgeType: EmitterNFTBridge,
		})
	}

	return out
}

func buildEmitterMap(hexmap map[vaa.ChainID]string) map[vaa.ChainID][]byte {
	out := make(map[vaa.ChainID][]byte)
	for id, emitter := range hexmap {
		e, err := hex.DecodeString(emitter)
		if err != nil {
			panic(fmt.Sprintf("Failed to decode emitter address %v: %v", emitter, err))
		}
		out[id] = e
	}

	return out
}

// KnownTokenbridgeEmitters is a list of well-known mainnet emitters for the tokenbridge.
var KnownTokenbridgeEmitters = buildEmitterMap(knownTokenbridgeEmitters)
var knownTokenbridgeEmitters = map[vaa.ChainID]string{
	vaa.ChainIDSolana:     "ec7372995d5cc8732397fb0ad35c0121e0eaa90d26f828a534cab54391b3a4f5",
	vaa.ChainIDEthereum:   "0000000000000000000000003ee18b2214aff97000d974cf647e7c347e8fa585",
	vaa.ChainIDTerra:      "0000000000000000000000007cf7b764e38a0a5e967972c1df77d432510564e2",
	vaa.ChainIDTerra2:     "a463ad028fb79679cfc8ce1efba35ac0e77b35080a1abe9bebe83461f176b0a3",
	vaa.ChainIDBSC:        "000000000000000000000000C891aBa0b42818fb4c975Bf6461033c62BCE75ff",
	vaa.ChainIDPolygon:    "0000000000000000000000008Eb8fB8B3c3d50140a19D35A71B3046543B37097",
	vaa.ChainIDPlanq:      "0000000000000000000000004FD8625cfE4B0034642140005b78291D26183df1",
	vaa.ChainIDTron:       "0000000000000000000000000000000000000000000000000000000000000000",
	vaa.ChainIDAvalanche:  "0000000000000000000000000e082f06ff657d94310cb8ce8b0d9a04541d8052",
	vaa.ChainIDOasis:      "0000000000000000000000005848c791e09901b40a9ef749f2a6735b418d7564",
	vaa.ChainIDAlgorand:   "67e93fa6c8ac5c819990aa7340c0c16b508abb1178be9b30d024b8ac25193d45",
	vaa.ChainIDAptos:      "0000000000000000000000000000000000000000000000000000000000000001",
	vaa.ChainIDAurora:     "00000000000000000000000051b5123a7b0F9b2bA265f9c4C8de7D78D52f510F",
	vaa.ChainIDFantom:     "0000000000000000000000007C9Fc5741288cDFdD83CeB07f3ea7e22618D79D2",
	vaa.ChainIDKarura:     "000000000000000000000000ae9d7fe007b3327AA64A32824Aaac52C42a6E624",
	vaa.ChainIDAcala:      "000000000000000000000000ae9d7fe007b3327AA64A32824Aaac52C42a6E624",
	vaa.ChainIDKlaytn:     "0000000000000000000000005b08ac39EAED75c0439FC750d9FE7E1F9dD0193F",
	vaa.ChainIDCelo:       "000000000000000000000000796Dff6D74F3E27060B71255Fe517BFb23C93eed",
	vaa.ChainIDNear:       "148410499d3fcda4dcfd68a1ebfcdddda16ab28326448d4aae4d2f0465cdfcb7",
	vaa.ChainIDMoonbeam:   "000000000000000000000000B1731c586ca89a23809861c6103F0b96B3F57D92",
	vaa.ChainIDArbitrum:   "0000000000000000000000000b2402144Bb366A632D14B83F244D2e0e21bD39c",
	vaa.ChainIDOptimism:   "0000000000000000000000001D68124e65faFC907325e3EDbF8c4d84499DAa8b",
	vaa.ChainIDBase:       "000000000000000000000000547169332126A398F67E02D1006d0F3955Eb552C",
	vaa.ChainIDXpla:       "8f9cf727175353b17a5f574270e370776123d90fd74956ae4277962b4fdee24c",
	vaa.ChainIDInjective:  "00000000000000000000000045dbea4617971d93188eda21530bc6503d153313",
	vaa.ChainIDSui:        "ccceeb29348f71bdd22ffef43a2a19c1f5b5e17c5cca5411529120182672ade5",
	vaa.ChainIDSei:        "86c5fd957e2db8389553e1728f9c27964b22a8154091ccba54d75f4b10c61f5e",
	vaa.ChainIDDeltachain: "aeb534c45c3049d380b9d9b966f9895f53abd4301bfaff407fa09dea8ae7a924",
}

// KnownNFTBridgeEmitters is a list of well-known mainnet emitters for the NFT bridge.
var KnownNFTBridgeEmitters = buildEmitterMap(knownNFTBridgeEmitters)
var knownNFTBridgeEmitters = map[vaa.ChainID]string{
	vaa.ChainIDSolana:    "0def15a24423e1edd1a5ab16f557b9060303ddbab8c803d2ee48f4b78a1cfd6b",
	vaa.ChainIDEthereum:  "0000000000000000000000006ffd7ede62328b3af38fcd61461bbfc52f5651fe",
	vaa.ChainIDBSC:       "0000000000000000000000002a1280866Fa742E50c93472B68B5026B558596e8",
	vaa.ChainIDPolygon:   "0000000000000000000000003379714f3720Dd57577456571D411Ec66C2f4551",
	vaa.ChainIDTron:      "0000000000000000000000000000000000000000000000000000000000000000",
	vaa.ChainIDPlanq:     "000000000000000000000000853348a8b1Db0db0A0F1e955fD7A90F84B03D050",
	vaa.ChainIDAvalanche: "000000000000000000000000f7b6737ca9c4e08ae573f75a97b73d7a813f5de5",
	vaa.ChainIDOasis:     "00000000000000000000000004952d522ff217f40b5ef3cbf659eca7b952a6c1",
	vaa.ChainIDAurora:    "0000000000000000000000006dcC0484472523ed9Cdc017F711Bcbf909789284",
	vaa.ChainIDFantom:    "000000000000000000000000A9c7119aBDa80d4a4E0C06C8F4d8cF5893234535",
	vaa.ChainIDKarura:    "000000000000000000000000b91e3638F82A1fACb28690b37e3aAE45d2c33808",
	vaa.ChainIDAcala:     "000000000000000000000000b91e3638F82A1fACb28690b37e3aAE45d2c33808",
	vaa.ChainIDKlaytn:    "0000000000000000000000003c3c561757BAa0b78c5C025CdEAa4ee24C1dFfEf",
	vaa.ChainIDCelo:      "000000000000000000000000A6A377d75ca5c9052c9a77ED1e865Cc25Bd97bf3",
	vaa.ChainIDMoonbeam:  "000000000000000000000000453cfBe096C0f8D763E8C5F24B441097d577bdE2",
	vaa.ChainIDArbitrum:  "0000000000000000000000003dD14D553cFD986EAC8e3bddF629d82073e188c8",
	vaa.ChainIDOptimism:  "000000000000000000000000fE8cD454b4A1CA468B57D79c0cc77Ef5B6f64585",
	vaa.ChainIDBase:      "00000000000000000000000033Ceb839c428c8244F30EDA589A0f52c747643e2",
	vaa.ChainIDAptos:     "0000000000000000000000000000000000000000000000000000000000000005",
}

func GetEmitterAddressForChain(chainID vaa.ChainID, emitterType EmitterType) (vaa.Address, error) {
	for _, emitter := range KnownEmitters {
		if emitter.ChainID == chainID && emitter.BridgeType == emitterType {
			emitterAddr, err := vaa.StringToAddress(emitter.Emitter)
			if err != nil {
				return vaa.Address{}, err
			}

			return emitterAddr, nil
		}
	}

	return vaa.Address{}, fmt.Errorf("lookup failed")
}
