package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/deltaswapio/deltachain/testutil/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/stretchr/testify/require"
)

func createNPhylaxSet(t *testing.T, keeper *keeper.Keeper, ctx sdk.Context, n int) []types.PhylaxSet {
	items := make([]types.PhylaxSet, n)
	for i := range items {
		items[i].Index = uint32(i)
		_, err := keeper.AppendPhylaxSet(ctx, items[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	return items
}

func TestPhylaxSetGet(t *testing.T) {
	keeper, ctx := keepertest.WormholeKeeper(t)
	items := createNPhylaxSet(t, keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetPhylaxSet(ctx, item.Index)
		require.True(t, found)
		require.Equal(t, item, got)
	}
}

func TestPhylaxSetGetAll(t *testing.T) {
	keeper, ctx := keepertest.WormholeKeeper(t)
	items := createNPhylaxSet(t, keeper, ctx, 10)
	require.ElementsMatch(t, items, keeper.GetAllPhylaxSet(ctx))
}

func TestPhylaxSetCount(t *testing.T) {
	keeper, ctx := keepertest.WormholeKeeper(t)
	items := createNPhylaxSet(t, keeper, ctx, 10)
	count := uint32(len(items))
	require.Equal(t, count, keeper.GetPhylaxSetCount(ctx))
}
