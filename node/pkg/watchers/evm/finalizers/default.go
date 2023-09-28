package finalizers

import (
	"context"

	"github.com/deltaswapio/deltaswap/node/pkg/watchers/evm/connectors"
)

// DefaultFinalizer assumes all blocks to be finalized.
type DefaultFinalizer struct {
}

func NewDefaultFinalizer() *DefaultFinalizer {
	return &DefaultFinalizer{}
}

func (d *DefaultFinalizer) IsBlockFinalized(ctx context.Context, block *connectors.NewBlock) (bool, error) {
	return true, nil
}
