package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/deltaswapio/deltachain/testutil/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
)

func createTestConfig(keeper *keeper.Keeper, ctx sdk.Context) types.Config {
	item := types.Config{}
	keeper.SetConfig(ctx, item)
	return item
}

func TestConfigGet(t *testing.T) {
	keeper, ctx := keepertest.DeltaswapKeeper(t)
	item := createTestConfig(keeper, ctx)
	rst, found := keeper.GetConfig(ctx)
	require.True(t, found)
	require.Equal(t, item, rst)
}
func TestConfigRemove(t *testing.T) {
	keeper, ctx := keepertest.DeltaswapKeeper(t)
	createTestConfig(keeper, ctx)
	keeper.RemoveConfig(ctx)
	_, found := keeper.GetConfig(ctx)
	require.False(t, found)
}
