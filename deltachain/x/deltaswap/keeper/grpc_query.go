package keeper

import (
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
)

var _ types.QueryServer = Keeper{}
