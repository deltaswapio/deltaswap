package types_test

import (
	"testing"

	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
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
					Index: 14,
				},
				PhylaxValidatorList: []types.PhylaxValidator{
					{
						PhylaxKey:     []byte{0},
						ValidatorAddr: []byte{3},
					},
					{
						PhylaxKey:     []byte{1},
						ValidatorAddr: []byte{4},
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated guardianSet",
			genState: &types.GenesisState{
				PhylaxSetList: []types.PhylaxSet{
					{
						Index: 0,
					},
					{
						Index: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated replayProtection",
			genState: &types.GenesisState{
				ReplayProtectionList: []types.ReplayProtection{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated sequenceCounter",
			genState: &types.GenesisState{
				SequenceCounterList: []types.SequenceCounter{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated phylaxValidator",
			genState: &types.GenesisState{
				PhylaxValidatorList: []types.PhylaxValidator{
					{
						PhylaxKey:     []byte{0},
						ValidatorAddr: []byte{10},
					},
					{
						PhylaxKey:     []byte{1},
						ValidatorAddr: []byte{10},
					},
				},
			},
			valid: true,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
