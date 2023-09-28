package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
)

// SetPhylaxValidator set a specific phylaxValidator in the store from its index
func (k Keeper) SetPhylaxValidator(ctx sdk.Context, phylaxValidator types.PhylaxValidator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PhylaxValidatorKeyPrefix))
	b := k.cdc.MustMarshal(&phylaxValidator)
	store.Set(types.PhylaxValidatorKey(
		phylaxValidator.PhylaxKey,
	), b)
}

// GetPhylaxValidator returns a phylaxValidator from its index
func (k Keeper) GetPhylaxValidator(
	ctx sdk.Context,
	phylaxKey []byte,

) (val types.PhylaxValidator, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PhylaxValidatorKeyPrefix))

	b := store.Get(types.PhylaxValidatorKey(
		phylaxKey,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePhylaxValidator removes a phylaxValidator from the store
func (k Keeper) RemovePhylaxValidator(
	ctx sdk.Context,
	phylaxKey []byte,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PhylaxValidatorKeyPrefix))
	store.Delete(types.PhylaxValidatorKey(
		phylaxKey,
	))
}

// GetAllPhylaxValidator returns all phylaxValidator
func (k Keeper) GetAllPhylaxValidator(ctx sdk.Context) (list []types.PhylaxValidator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PhylaxValidatorKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PhylaxValidator
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
