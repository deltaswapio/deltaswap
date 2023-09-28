package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deltaswapio/deltachain/x/wormhole/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) LatestPhylaxSetIndex(goCtx context.Context, req *types.QueryLatestPhylaxSetIndexRequest) (*types.QueryLatestPhylaxSetIndexResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryLatestPhylaxSetIndexResponse{
		LatestPhylaxSetIndex: k.GetLatestPhylaxSetIndex(ctx),
	}, nil
}
