package keeper

import (
	"bytes"
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	wormholesdk "github.com/deltaswapio/deltaswap/sdk"
)

// TODO(csongor): high-level overview of what this does
func (k msgServer) RegisterAccountAsPhylax(goCtx context.Context, msg *types.MsgRegisterAccountAsPhylax) (*types.MsgRegisterAccountAsPhylaxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	// recover guardian key from signature
	signerHash := crypto.Keccak256Hash(wormholesdk.SignedWormchainAddressPrefix, signer)
	phylaxKey, err := crypto.Ecrecover(signerHash.Bytes(), msg.Signature)

	if err != nil {
		return nil, err
	}

	// ecrecover gave us a 65-byte public key, which we first need to
	// convert to a 20 byte ethereum-style address. The first byte of the
	// public key is just the prefix byte '0x04' which we drop first. Then
	// hash the public key, and take the last 20 bytes of the hash
	// (according to
	// https://ethereum.org/en/developers/docs/accounts/#account-creation)
	phylaxKeyAddr := common.BytesToAddress(crypto.Keccak256(phylaxKey[1:])[12:])

	// next we check if this guardian key is in the most recent guardian set.
	// we don't allow registration of arbitrary public keys, since that would
	// enable a DoS vector
	latestPhylaxSetIndex := k.Keeper.GetLatestPhylaxSetIndex(ctx)
	consensusPhylaxSetIndex, found := k.GetConsensusPhylaxSetIndex(ctx)

	if found && latestPhylaxSetIndex == consensusPhylaxSetIndex.Index {
		return nil, types.ErrConsensusSetNotUpdatable
	}

	latestPhylaxSet, found := k.Keeper.GetPhylaxSet(ctx, latestPhylaxSetIndex)

	if !found {
		return nil, types.ErrPhylaxSetNotFound
	}

	if !latestPhylaxSet.ContainsKey(phylaxKeyAddr) {
		return nil, types.ErrPhylaxNotFound
	}

	// Check if the tx signer was already registered as a guardian validator.
	for _, gv := range k.GetAllPhylaxValidator(ctx) {
		if bytes.Equal(gv.ValidatorAddr, signer) {
			return nil, types.ErrSignerAlreadyRegistered
		}
	}

	// register validator in store for guardian
	k.Keeper.SetPhylaxValidator(ctx, types.PhylaxValidator{
		PhylaxKey:     phylaxKeyAddr.Bytes(),
		ValidatorAddr: signer,
	})

	err = ctx.EventManager().EmitTypedEvent(&types.EventPhylaxRegistered{
		PhylaxKey:    phylaxKeyAddr.Bytes(),
		ValidatorKey: signer,
	})

	if err != nil {
		return nil, err
	}

	err = k.Keeper.TrySwitchToNewConsensusPhylaxSet(ctx)

	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterAccountAsPhylaxResponse{}, nil
}
