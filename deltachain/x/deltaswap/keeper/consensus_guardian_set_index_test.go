package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/deltaswapio/deltachain/testutil/keeper"
	"github.com/deltaswapio/deltachain/testutil/nullify"
	"github.com/deltaswapio/deltachain/x/deltaswap/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
)

func createTestConsensusPhylaxSetIndex(keeper *keeper.Keeper, ctx sdk.Context) types.ConsensusPhylaxSetIndex {
	item := types.ConsensusPhylaxSetIndex{}
	keeper.SetConsensusPhylaxSetIndex(ctx, item)
	return item
}

func TestConsensusPhylaxSetIndexGet(t *testing.T) {
	keeper, ctx := keepertest.DeltaswapKeeper(t)
	item := createTestConsensusPhylaxSetIndex(keeper, ctx)
	rst, found := keeper.GetConsensusPhylaxSetIndex(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestConsensusPhylaxSetIndexRemove(t *testing.T) {
	keeper, ctx := keepertest.DeltaswapKeeper(t)
	createTestConsensusPhylaxSetIndex(keeper, ctx)
	keeper.RemoveConsensusPhylaxSetIndex(ctx)
	_, found := keeper.GetConsensusPhylaxSetIndex(ctx)
	require.False(t, found)
}
