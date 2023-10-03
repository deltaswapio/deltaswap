package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PhylaxValidatorAll(c context.Context, req *types.QueryAllPhylaxValidatorRequest) (*types.QueryAllPhylaxValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var phylaxValidators []types.PhylaxValidator
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	phylaxValidatorStore := prefix.NewStore(store, types.KeyPrefix(types.PhylaxValidatorKeyPrefix))

	pageRes, err := query.Paginate(phylaxValidatorStore, req.Pagination, func(key []byte, value []byte) error {
		var phylaxValidator types.PhylaxValidator
		if err := k.cdc.Unmarshal(value, &phylaxValidator); err != nil {
			return err
		}

		phylaxValidators = append(phylaxValidators, phylaxValidator)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPhylaxValidatorResponse{PhylaxValidator: phylaxValidators, Pagination: pageRes}, nil
}

func (k Keeper) PhylaxValidator(c context.Context, req *types.QueryGetPhylaxValidatorRequest) (*types.QueryGetPhylaxValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetPhylaxValidator(
		ctx,
		req.PhylaxKey,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetPhylaxValidatorResponse{PhylaxValidator: val}, nil
}
