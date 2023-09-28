package cli_test

import (
	"fmt"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/status"

	"github.com/deltaswapio/deltachain/testutil/network"
	"github.com/deltaswapio/deltachain/testutil/nullify"
	"github.com/deltaswapio/deltachain/x/deltaswap/client/cli"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
)

func networkWithConsensusPhylaxSetIndexObjects(t *testing.T) (*network.Network, types.ConsensusPhylaxSetIndex) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	consensusPhylaxSetIndex := &types.ConsensusPhylaxSetIndex{
		Index: 0,
	}
	nullify.Fill(&consensusPhylaxSetIndex)
	state.ConsensusPhylaxSetIndex = consensusPhylaxSetIndex

	phylaxSetList := []types.PhylaxSet{{
		Index:          0,
		Keys:           [][]byte{},
		ExpirationTime: 0,
	}}
	nullify.Fill(&phylaxSetList)
	state.PhylaxSetList = phylaxSetList

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), *state.ConsensusPhylaxSetIndex
}

func TestShowConsensusPhylaxSetIndex(t *testing.T) {
	net, obj := networkWithConsensusPhylaxSetIndexObjects(t)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		args []string
		err  error
		obj  types.ConsensusPhylaxSetIndex
	}{
		{
			desc: "get",
			args: common,
			obj:  obj,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowConsensusPhylaxSetIndex(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetConsensusPhylaxSetIndexResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.ConsensusPhylaxSetIndex)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.ConsensusPhylaxSetIndex),
				)
			}
		})
	}
}
