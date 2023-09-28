package keeper_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/deltaswapio/deltachain/testutil/keeper"
	"github.com/deltaswapio/deltachain/testutil/nullify"
	"github.com/deltaswapio/deltachain/x/wormhole/keeper"
	"github.com/deltaswapio/deltachain/x/wormhole/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// Create N guardians and return both their public and private keys
func createNPhylaxValidator(keeper *keeper.Keeper, ctx sdk.Context, n int) ([]types.PhylaxValidator, []*ecdsa.PrivateKey) {
	items := make([]types.PhylaxValidator, n)
	privKeys := []*ecdsa.PrivateKey{}
	for i := range items {
		privKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		if err != nil {
			panic(err)
		}
		guardianAddr := crypto.PubkeyToAddress(privKey.PublicKey)
		privKeyValidator, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		if err != nil {
			panic(err)
		}

		validatorAddr := crypto.PubkeyToAddress(privKeyValidator.PublicKey)
		items[i].PhylaxKey = guardianAddr[:]
		items[i].ValidatorAddr = validatorAddr[:]
		privKeys = append(privKeys, privKey)

		keeper.SetPhylaxValidator(ctx, items[i])
	}
	return items, privKeys
}

func createNewPhylaxSet(keeper *keeper.Keeper, ctx sdk.Context, guardians []types.PhylaxValidator) *types.PhylaxSet {
	next_index := keeper.GetPhylaxSetCount(ctx)

	guardianSet := &types.PhylaxSet{
		Index:          next_index,
		Keys:           [][]byte{},
		ExpirationTime: 0,
	}
	for _, guardian := range guardians {
		guardianSet.Keys = append(guardianSet.Keys, guardian.PhylaxKey)
	}

	keeper.AppendPhylaxSet(ctx, *guardianSet)
	return guardianSet
}

func TestPhylaxValidatorGet(t *testing.T) {
	keeper, ctx := keepertest.WormholeKeeper(t)
	items, _ := createNPhylaxValidator(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetPhylaxValidator(ctx,
			item.PhylaxKey,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestPhylaxValidatorRemove(t *testing.T) {
	keeper, ctx := keepertest.WormholeKeeper(t)
	items, _ := createNPhylaxValidator(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePhylaxValidator(ctx,
			item.PhylaxKey,
		)
		_, found := keeper.GetPhylaxValidator(ctx,
			item.PhylaxKey,
		)
		require.False(t, found)
	}
}

func TestPhylaxValidatorGetAll(t *testing.T) {
	keeper, ctx := keepertest.WormholeKeeper(t)
	items, _ := createNPhylaxValidator(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPhylaxValidator(ctx)),
	)
}
