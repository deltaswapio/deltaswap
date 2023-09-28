package cli_test

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/deltaswapio/deltachain/testutil/network"
	"github.com/deltaswapio/deltachain/testutil/nullify"
	"github.com/deltaswapio/deltachain/x/deltaswap/client/cli"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithPhylaxValidatorObjects(t *testing.T, n int) (*network.Network, []types.PhylaxValidator) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		guardianValidator := types.PhylaxValidator{
			PhylaxKey:     []byte(strconv.Itoa(i)),
			ValidatorAddr: []byte(strconv.Itoa(i)),
		}
		state.PhylaxValidatorList = append(state.PhylaxValidatorList, guardianValidator)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.PhylaxValidatorList
}

func TestShowPhylaxValidator(t *testing.T) {
	net, objs := networkWithPhylaxValidatorObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc        string
		idPhylaxKey string

		args []string
		err  error
		obj  types.PhylaxValidator
	}{
		{
			desc:        "found",
			idPhylaxKey: hex.EncodeToString(objs[0].PhylaxKey),

			args: common,
			obj:  objs[0],
		},
		{
			desc:        "not found",
			idPhylaxKey: "0x100000",

			args: common,
			err:  status.Error(codes.InvalidArgument, "not found"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idPhylaxKey,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowPhylaxValidator(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetPhylaxValidatorResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.PhylaxValidator)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.PhylaxValidator),
				)
			}
		})
	}
}

func TestListPhylaxValidator(t *testing.T) {
	net, objs := networkWithPhylaxValidatorObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPhylaxValidator(), args)
			require.NoError(t, err)
			var resp types.QueryAllPhylaxValidatorResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.PhylaxValidator), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.PhylaxValidator),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPhylaxValidator(), args)
			require.NoError(t, err)
			var resp types.QueryAllPhylaxValidatorResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.PhylaxValidator), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.PhylaxValidator),
			)
			next = resp.Pagination.NextKey
		}
	})
	// TODO(csongor): this test is failing, figure out why
	// t.Run("Total", func(t *testing.T) {
	// 	args := request(nil, 0, uint64(len(objs)), true)
	// 	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPhylaxValidator(), args)
	// 	require.NoError(t, err)
	// 	var resp types.QueryAllPhylaxValidatorResponse
	// 	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	// 	require.NoError(t, err)
	// 	require.Equal(t, len(objs), int(resp.Pagination.Total))
	// 	require.ElementsMatch(t,
	// 		nullify.Fill(objs),
	// 		nullify.Fill(resp.PhylaxValidator),
	// 	)
	// })
}
