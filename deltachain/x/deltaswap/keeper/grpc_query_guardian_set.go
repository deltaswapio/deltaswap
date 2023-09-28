package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PhylaxSetAll(c context.Context, req *types.QueryAllPhylaxSetRequest) (*types.QueryAllPhylaxSetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var guardianSets []types.PhylaxSet
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	guardianSetStore := prefix.NewStore(store, types.KeyPrefix(types.PhylaxSetKey))

	pageRes, err := query.Paginate(guardianSetStore, req.Pagination, func(key []byte, value []byte) error {
		var guardianSet types.PhylaxSet
		if err := k.cdc.Unmarshal(value, &guardianSet); err != nil {
			return err
		}

		guardianSets = append(guardianSets, guardianSet)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPhylaxSetResponse{PhylaxSet: guardianSets, Pagination: pageRes}, nil
}

func (k Keeper) PhylaxSet(c context.Context, req *types.QueryGetPhylaxSetRequest) (*types.QueryGetPhylaxSetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	guardianSet, found := k.GetPhylaxSet(ctx, req.Index)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetPhylaxSetResponse{PhylaxSet: guardianSet}, nil
}
