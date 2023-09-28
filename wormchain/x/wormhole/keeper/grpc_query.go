package keeper

import (
	"github.com/wormhole-foundation/deltachain/x/wormhole/types"
)

var _ types.QueryServer = Keeper{}
