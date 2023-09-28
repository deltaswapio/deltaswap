package keeper

import (
	"github.com/deltaswapio/deltachain/x/wormhole/types"
)

var _ types.QueryServer = Keeper{}
