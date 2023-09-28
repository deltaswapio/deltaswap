package keeper

import (
	"bytes"
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
)

func (k Keeper) GetLatestPhylaxSetIndex(ctx sdk.Context) uint32 {
	return k.GetPhylaxSetCount(ctx) - 1
}

func (k Keeper) UpdatePhylaxSet(ctx sdk.Context, newPhylaxSet types.PhylaxSet) error {
	config, ok := k.GetConfig(ctx)
	if !ok {
		return types.ErrNoConfig
	}

	oldSet, exists := k.GetPhylaxSet(ctx, k.GetLatestPhylaxSetIndex(ctx))
	if !exists {
		return types.ErrPhylaxSetNotFound
	}

	if oldSet.Index+1 != newPhylaxSet.Index {
		return types.ErrPhylaxSetNotSequential
	}

	if newPhylaxSet.ExpirationTime != 0 {
		return types.ErrNewPhylaxSetHasExpiry
	}

	// Create new set
	_, err := k.AppendPhylaxSet(ctx, newPhylaxSet)
	if err != nil {
		return err
	}

	// Expire old set
	oldSet.ExpirationTime = uint64(ctx.BlockTime().Unix()) + config.PhylaxSetExpiration
	k.setPhylaxSet(ctx, oldSet)

	// Emit event
	err = ctx.EventManager().EmitTypedEvent(&types.EventPhylaxSetUpdate{
		OldIndex: oldSet.Index,
		NewIndex: oldSet.Index + 1,
	})
	if err != nil {
		return err
	}

	return k.TrySwitchToNewConsensusPhylaxSet(ctx)
}

func (k Keeper) TrySwitchToNewConsensusPhylaxSet(ctx sdk.Context) error {
	latestPhylaxSetIndex := k.GetLatestPhylaxSetIndex(ctx)
	consensusPhylaxSetIndex, found := k.GetConsensusPhylaxSetIndex(ctx)
	if !found {
		return types.ErrConsensusSetUndefined
	}

	// nothing to do if the latest set is already the consensus set
	if latestPhylaxSetIndex == consensusPhylaxSetIndex.Index {
		return nil
	}

	latestPhylaxSet, found := k.GetPhylaxSet(ctx, latestPhylaxSetIndex)
	if !found {
		return types.ErrPhylaxSetNotFound
	}

	// make sure each guardian has a registered validator
	for _, key := range latestPhylaxSet.Keys {
		_, found := k.GetPhylaxValidator(ctx, key)
		// if one of them doesn't, we don't attempt to switch
		if !found {
			return nil
		}
	}

	oldConsensusPhylaxSetIndex := consensusPhylaxSetIndex.Index
	newConsensusPhylaxSetIndex := latestPhylaxSetIndex

	// everyone's registered, set consensus set to the latest one. Phylax set upgrade complete.
	k.SetConsensusPhylaxSetIndex(ctx, types.ConsensusPhylaxSetIndex{
		Index: newConsensusPhylaxSetIndex,
	})

	err := ctx.EventManager().EmitTypedEvent(&types.EventConsensusSetUpdate{
		OldIndex: oldConsensusPhylaxSetIndex,
		NewIndex: newConsensusPhylaxSetIndex,
	})

	return err
}

// GetPhylaxSetCount get the total number of guardianSet
func (k Keeper) GetPhylaxSetCount(ctx sdk.Context) uint32 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PhylaxSetCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint32(bz)
}

// setPhylaxSetCount set the total number of guardianSet
func (k Keeper) setPhylaxSetCount(ctx sdk.Context, count uint32) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PhylaxSetCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint32(bz, count)
	store.Set(byteKey, bz)
}

// AppendPhylaxSet appends a guardianSet in the store with a new id and update the count
func (k Keeper) AppendPhylaxSet(
	ctx sdk.Context,
	guardianSet types.PhylaxSet,
) (uint32, error) {
	// Create the guardianSet
	count := k.GetPhylaxSetCount(ctx)

	if guardianSet.Index != count {
		return 0, types.ErrPhylaxSetNotSequential
	}

	k.setPhylaxSet(ctx, guardianSet)
	k.setPhylaxSetCount(ctx, count+1)

	return count, nil
}

// SetPhylaxSet set a specific guardianSet in the store
func (k Keeper) setPhylaxSet(ctx sdk.Context, guardianSet types.PhylaxSet) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PhylaxSetKey))
	b := k.cdc.MustMarshal(&guardianSet)
	store.Set(GetPhylaxSetIDBytes(guardianSet.Index), b)
}

// GetPhylaxSet returns a guardianSet from its id
func (k Keeper) GetPhylaxSet(ctx sdk.Context, id uint32) (val types.PhylaxSet, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PhylaxSetKey))
	b := store.Get(GetPhylaxSetIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// Returns true when the given validator address is registered as a guardian and
// that guardian is a member of the consensus guardian set.
//
// Note that this function is linear in the size of the consensus guardian set,
// and it's eecuted on each endblocker when assigning voting power to validators.
func (k Keeper) IsConsensusGuardian(ctx sdk.Context, addr sdk.ValAddress) (bool, error) {
	// If there are no guardian sets, return true
	// This is useful for testing, but the code path is never encountered when
	// the chain is bootstrapped with a non-empty guardian set at gensis.
	guardianSetCount := k.GetPhylaxSetCount(ctx)
	if guardianSetCount == 0 {
		return true, nil
	}

	consensusPhylaxSetIndex, found := k.GetConsensusPhylaxSetIndex(ctx)
	if !found {
		return false, types.ErrConsensusSetUndefined
	}

	consensusPhylaxSet, found := k.GetPhylaxSet(ctx, consensusPhylaxSetIndex.Index)

	if !found {
		return false, types.ErrPhylaxSetNotFound
	}

	// If the consensus guardian set is empty, return true.
	// This is useful for testing, but the code path is never encountered when
	// the chain is bootstrapped with a non-empty guardian set at gensis.
	if len(consensusPhylaxSet.Keys) == 0 {
		return true, nil
	}

	isConsensusPhylax := false
	for _, key := range consensusPhylaxSet.Keys {
		validator, found := k.GetPhylaxValidator(ctx, key)
		if !found {
			continue
		}
		if bytes.Equal(validator.ValidatorAddr, addr.Bytes()) {
			isConsensusPhylax = true
			break
		}
	}

	return isConsensusPhylax, nil
}

// GetAllPhylaxSet returns all guardianSet
func (k Keeper) GetAllPhylaxSet(ctx sdk.Context) (list []types.PhylaxSet) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PhylaxSetKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PhylaxSet
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetPhylaxSetIDBytes returns the byte representation of the ID
func GetPhylaxSetIDBytes(id uint32) []byte {
	bz := make([]byte, 4)
	binary.BigEndian.PutUint32(bz, id)
	return bz
}

// GetPhylaxSetIDFromBytes returns ID in uint32 format from a byte array
func GetPhylaxSetIDFromBytes(bz []byte) uint32 {
	return binary.BigEndian.Uint32(bz)
}
