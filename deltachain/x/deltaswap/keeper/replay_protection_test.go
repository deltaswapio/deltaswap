package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/deltaswapio/deltachain/testutil/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNReplayProtection(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.ReplayProtection {
	items := make([]types.ReplayProtection, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetReplayProtection(ctx, items[i])
	}
	return items
}

func TestReplayProtectionGet(t *testing.T) {
	keeper, ctx := keepertest.DeltaswapKeeper(t)
	items := createNReplayProtection(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetReplayProtection(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}
func TestReplayProtectionRemove(t *testing.T) {
	keeper, ctx := keepertest.DeltaswapKeeper(t)
	items := createNReplayProtection(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveReplayProtection(ctx,
			item.Index,
		)
		_, found := keeper.GetReplayProtection(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestReplayProtectionGetAll(t *testing.T) {
	keeper, ctx := keepertest.DeltaswapKeeper(t)
	items := createNReplayProtection(keeper, ctx, 10)
	require.ElementsMatch(t, items, keeper.GetAllReplayProtection(ctx))
}
