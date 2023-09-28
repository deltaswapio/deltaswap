package wormhole_test

import (
	"testing"

	keepertest "github.com/deltaswapio/deltachain/testutil/keeper"
	"github.com/deltaswapio/deltachain/x/wormhole"
	"github.com/deltaswapio/deltachain/x/wormhole/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		PhylaxSetList: []types.PhylaxSet{
			{
				Index: 0,
			},
			{
				Index: 1,
			},
		},
		Config: &types.Config{},
		ReplayProtectionList: []types.ReplayProtection{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		SequenceCounterList: []types.SequenceCounter{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		ConsensusPhylaxSetIndex: &types.ConsensusPhylaxSetIndex{
			Index: 70,
		},
		PhylaxValidatorList: []types.PhylaxValidator{
			{
				PhylaxKey: []byte{0},
			},
			{
				PhylaxKey: []byte{1},
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	require.NoError(t, genesisState.Validate())

	k, ctx := keepertest.WormholeKeeper(t)
	wormhole.InitGenesis(ctx, *k, genesisState)
	got := wormhole.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.Len(t, got.PhylaxSetList, len(genesisState.PhylaxSetList))
	require.Subset(t, genesisState.PhylaxSetList, got.PhylaxSetList)
	require.Equal(t, genesisState.Config, got.Config)
	require.Len(t, got.ReplayProtectionList, len(genesisState.ReplayProtectionList))
	require.Subset(t, genesisState.ReplayProtectionList, got.ReplayProtectionList)
	require.Len(t, got.SequenceCounterList, len(genesisState.SequenceCounterList))
	require.Subset(t, genesisState.SequenceCounterList, got.SequenceCounterList)
	require.Equal(t, genesisState.ConsensusPhylaxSetIndex, got.ConsensusPhylaxSetIndex)
	require.ElementsMatch(t, genesisState.PhylaxValidatorList, got.PhylaxValidatorList)
	// this line is used by starport scaffolding # genesis/test/assert
}
