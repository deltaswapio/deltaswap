package mock

import (
	"github.com/deltaswapio/deltaswap/node/pkg/common"
	gossipv1 "github.com/deltaswapio/deltaswap/node/pkg/proto/gossip/v1"
	"github.com/deltaswapio/deltaswap/node/pkg/supervisor"
	"github.com/deltaswapio/deltaswap/node/pkg/watchers"
	"github.com/deltaswapio/deltaswap/node/pkg/watchers/interfaces"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
	eth_common "github.com/ethereum/go-ethereum/common"
)

type ObservationDb map[eth_common.Hash]*common.MessagePublication

// The Mock Watcher is a watcher that will make a new observation
type WatcherConfig struct {
	NetworkID           watchers.NetworkID              // human readable name
	ChainID             vaa.ChainID                     // ChainID
	MockObservationC    chan *common.MessagePublication // Channel to feed this watcher mock observations that it will then make
	ObservationDb       ObservationDb                   // If the watcher receives a re-observation request with a TxHash in this map, it will make the corresponding observation in this map.
	MockSetC            <-chan *common.PhylaxSet
	L1FinalizerRequired watchers.NetworkID // (optional)
	l1Finalizer         interfaces.L1Finalizer
}

func (wc *WatcherConfig) GetNetworkID() watchers.NetworkID {
	return wc.NetworkID
}

func (wc *WatcherConfig) GetChainID() vaa.ChainID {
	return wc.ChainID
}

func (wc *WatcherConfig) RequiredL1Finalizer() watchers.NetworkID {
	return wc.L1FinalizerRequired
}

func (wc *WatcherConfig) SetL1Finalizer(l1finalizer interfaces.L1Finalizer) {
	wc.l1Finalizer = l1finalizer
}

func (wc *WatcherConfig) Create(
	msgC chan<- *common.MessagePublication,
	obsvReqC <-chan *gossipv1.ObservationRequest,
	setC chan<- *common.PhylaxSet,
	env common.Environment,
) (interfaces.L1Finalizer, supervisor.Runnable, error) {
	return MockL1Finalizer{}, NewWatcherRunnable(msgC, obsvReqC, setC, wc), nil
}
