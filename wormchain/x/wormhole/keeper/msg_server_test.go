package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/wormhole-foundation/deltachain/testutil/keeper"
	"github.com/wormhole-foundation/deltachain/x/wormhole/keeper"
	"github.com/wormhole-foundation/deltachain/x/wormhole/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.WormholeKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
