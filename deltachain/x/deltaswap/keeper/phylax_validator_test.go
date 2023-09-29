package keeper_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/deltaswapio/deltachain/testutil/keeper"
	"github.com/deltaswapio/deltachain/testutil/nullify"
	"github.com/deltaswapio/deltachain/x/deltaswap/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// Create N phylaxs and return both their public and private keys
func createNPhylaxValidator(keeper *keeper.Keeper, ctx sdk.Context, n int) ([]types.PhylaxValidator, []*ecdsa.PrivateKey) {
	items := make([]types.PhylaxValidator, n)
	privKeys := []*ecdsa.PrivateKey{}
	for i := range items {
		privKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		if err != nil {
			panic(err)
		}
		phylaxAddr := crypto.PubkeyToAddress(privKey.PublicKey)
		privKeyValidator, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		if err != nil {
			panic(err)
		}

		validatorAddr := crypto.PubkeyToAddress(privKeyValidator.PublicKey)
		items[i].PhylaxKey = phylaxAddr[:]
		items[i].ValidatorAddr = validatorAddr[:]
		privKeys = append(privKeys, privKey)

		keeper.SetPhylaxValidator(ctx, items[i])
	}
	return items, privKeys
}

func createNewPhylaxSet(keeper *keeper.Keeper, ctx sdk.Context, phylaxs []types.PhylaxValidator) *types.PhylaxSet {
	next_index := keeper.GetPhylaxSetCount(ctx)

	phylaxSet := &types.PhylaxSet{
		Index:          next_index,
		Keys:           [][]byte{},
		ExpirationTime: 0,
	}
	for _, phylax := range phylaxs {
		phylaxSet.Keys = append(phylaxSet.Keys, phylax.PhylaxKey)
	}

	keeper.AppendPhylaxSet(ctx, *phylaxSet)
	return phylaxSet
}

func TestPhylaxValidatorGet(t *testing.T) {
	keeper, ctx := keepertest.DeltaswapKeeper(t)
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
	keeper, ctx := keepertest.DeltaswapKeeper(t)
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
	keeper, ctx := keepertest.DeltaswapKeeper(t)
	items, _ := createNPhylaxValidator(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPhylaxValidator(ctx)),
	)
}
