package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deltaswapio/deltachain/x/wormhole/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ConsensusPhylaxSetIndex(c context.Context, req *types.QueryGetConsensusPhylaxSetIndexRequest) (*types.QueryGetConsensusPhylaxSetIndexResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetConsensusPhylaxSetIndex(ctx)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetConsensusPhylaxSetIndexResponse{ConsensusPhylaxSetIndex: val}, nil
}
