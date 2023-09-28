package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/deltaswapio/deltachain/testutil/keeper"
	"github.com/deltaswapio/deltachain/testutil/nullify"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
)

func TestConsensusPhylaxSetIndexQuery(t *testing.T) {
	keeper, ctx := keepertest.DeltaswapKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestConsensusPhylaxSetIndex(keeper, ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetConsensusPhylaxSetIndexRequest
		response *types.QueryGetConsensusPhylaxSetIndexResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetConsensusPhylaxSetIndexRequest{},
			response: &types.QueryGetConsensusPhylaxSetIndexResponse{ConsensusPhylaxSetIndex: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.ConsensusPhylaxSetIndex(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
