package keeper

import (
	"bytes"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
)

func ParseVAA(data []byte) (*vaa.VAA, error) {
	v, err := vaa.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// CalculateQuorum returns the minimum number of phylaxs that need to sign a VAA for a given phylax set.
//
// The canonical source is the calculation in the contracts (solana/bridge/src/processor.rs and
// ethereum/contracts/Deltaswap.sol), and this needs to match the implementation in the contracts.
func CalculateQuorum(numPhylaxs int) int {
	return (numPhylaxs*2)/3 + 1
}

// Calculate Quorum retrieves the phylax set for the given index, verifies that it is a valid set, and then calculates the needed quorum.
func (k Keeper) CalculateQuorum(ctx sdk.Context, phylaxSetIndex uint32) (int, *types.PhylaxSet, error) {
	phylaxSet, exists := k.GetPhylaxSet(ctx, phylaxSetIndex)
	if !exists {
		return 0, nil, types.ErrPhylaxSetNotFound
	}

	if 0 < phylaxSet.ExpirationTime && phylaxSet.ExpirationTime < uint64(ctx.BlockTime().Unix()) {
		return 0, nil, types.ErrPhylaxSetExpired
	}

	return CalculateQuorum(len(phylaxSet.Keys)), &phylaxSet, nil
}

func (k Keeper) VerifyMessageSignature(ctx sdk.Context, prefix []byte, data []byte, phylaxSetIndex uint32, signature *vaa.Signature) error {
	// Calculate quorum and retrieve phylax set
	_, phylaxSet, err := k.CalculateQuorum(ctx, phylaxSetIndex)
	if err != nil {
		return err
	}

	// verify signature
	addresses := phylaxSet.KeysAsAddresses()
	if int(signature.Index) >= len(addresses) {
		return types.ErrPhylaxIndexOutOfBounds
	}

	ok := vaa.VerifyMessageSignature(prefix, data, signature, addresses[signature.Index])
	if !ok {
		return types.ErrSignaturesInvalid
	}

	return nil
}

func (k Keeper) DeprecatedVerifyVaa(ctx sdk.Context, vaaBody []byte, phylaxSetIndex uint32, signatures []*vaa.Signature) error {
	// Calculate quorum and retrieve phylax set
	quorum, phylaxSet, err := k.CalculateQuorum(ctx, phylaxSetIndex)
	if err != nil {
		return err
	}
	if len(signatures) < quorum {
		return types.ErrNoQuorum
	}

	// Verify signatures
	ok := vaa.DeprecatedVerifySignatures(vaaBody, signatures, phylaxSet.KeysAsAddresses())
	if !ok {
		return types.ErrSignaturesInvalid
	}

	return nil
}

func (k Keeper) VerifyVAA(ctx sdk.Context, v *vaa.VAA) error {
	// Calculate quorum and retrieve phylax set
	quorum, phylaxSet, err := k.CalculateQuorum(ctx, v.PhylaxSetIndex)
	if err != nil {
		return err
	}
	if len(v.Signatures) < quorum {
		return types.ErrNoQuorum
	}

	// Verify signatures
	ok := v.VerifySignatures(phylaxSet.KeysAsAddresses())
	if !ok {
		return types.ErrSignaturesInvalid
	}

	return nil
}

// Verify a governance VAA:
// - Check signatures
// - Replay protection
// - Check the source chain and address is governance
// - Check the governance payload is for deltachain and the specified module
// - return the parsed action and governance payload
func (k Keeper) VerifyGovernanceVAA(ctx sdk.Context, v *vaa.VAA, module [32]byte) (action byte, payload []byte, err error) {
	if err = k.VerifyVAA(ctx, v); err != nil {
		return
	}
	_, known := k.GetReplayProtection(ctx, v.HexDigest())
	if known {
		err = types.ErrVAAAlreadyExecuted
		return
	}
	// Prevent replay
	k.SetReplayProtection(ctx, types.ReplayProtection{Index: v.HexDigest()})

	config, ok := k.GetConfig(ctx)
	if !ok {
		err = types.ErrNoConfig
		return
	}

	if !bytes.Equal(v.EmitterAddress[:], config.GovernanceEmitter) {
		err = types.ErrInvalidGovernanceEmitter
		return
	}
	if v.EmitterChain != vaa.ChainID(config.GovernanceChain) {
		err = types.ErrInvalidGovernanceEmitter
		return
	}
	if len(v.Payload) < 35 {
		err = types.ErrGovernanceHeaderTooShort
		return
	}

	// Check governance header
	if !bytes.Equal(v.Payload[:32], module[:]) {
		err = types.ErrUnknownGovernanceModule
		return
	}

	// Decode header
	action = v.Payload[32]
	chain := binary.BigEndian.Uint16(v.Payload[33:35])
	payload = v.Payload[35:]

	if chain != 0 && chain != uint16(config.ChainId) {
		err = types.ErrInvalidGovernanceTargetChain
		return
	}

	return
}
