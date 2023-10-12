package watchers

import (
	"github.com/deltaswapio/deltaswap/node/pkg/common"
	gossipv1 "github.com/deltaswapio/deltaswap/node/pkg/proto/gossip/v1"
	"github.com/deltaswapio/deltaswap/node/pkg/query"
	"github.com/deltaswapio/deltaswap/node/pkg/supervisor"
	"github.com/deltaswapio/deltaswap/node/pkg/watchers/interfaces"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
)

// NetworkID is a unique identifier of a watcher that is used to link watchers together for the purpose of L1 Finalizers.
// This is different from vaa.ChainID because there could be multiple watchers for a single chain (e.g. solana-confirmed and solana-finalized)
type NetworkID string

type WatcherConfig interface {
	GetNetworkID() NetworkID
	GetChainID() vaa.ChainID
	RequiredL1Finalizer() NetworkID // returns NetworkID of the L1 Finalizer that should be used for this Watcher.
	SetL1Finalizer(l1finalizer interfaces.L1Finalizer)
	Create(
		msgC chan<- *common.MessagePublication,
		obsvReqC <-chan *gossipv1.ObservationRequest,
		queryReqC <-chan *query.PerChainQueryInternal,
		queryResponseC chan<- *query.PerChainQueryResponseInternal,
		setC chan<- *common.PhylaxSet,
		env common.Environment,
	) (interfaces.L1Finalizer, supervisor.Runnable, error)
}
