package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deltaswapio/deltachain/x/wormhole/types"
)

// SetConsensusPhylaxSetIndex set consensusPhylaxSetIndex in the store
func (k Keeper) SetConsensusPhylaxSetIndex(ctx sdk.Context, consensusPhylaxSetIndex types.ConsensusPhylaxSetIndex) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ConsensusPhylaxSetIndexKey))
	b := k.cdc.MustMarshal(&consensusPhylaxSetIndex)
	store.Set([]byte{0}, b)
}

// GetConsensusPhylaxSetIndex returns consensusPhylaxSetIndex
func (k Keeper) GetConsensusPhylaxSetIndex(ctx sdk.Context) (val types.ConsensusPhylaxSetIndex, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ConsensusPhylaxSetIndexKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveConsensusPhylaxSetIndex removes consensusPhylaxSetIndex from the store
func (k Keeper) RemoveConsensusPhylaxSetIndex(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ConsensusPhylaxSetIndexKey))
	store.Delete([]byte{0})
}
