package sui

import (
	"github.com/deltaswapio/deltaswap/node/pkg/common"
	gossipv1 "github.com/deltaswapio/deltaswap/node/pkg/proto/gossip/v1"
	"github.com/deltaswapio/deltaswap/node/pkg/supervisor"
	"github.com/deltaswapio/deltaswap/node/pkg/watchers"
	"github.com/deltaswapio/deltaswap/node/pkg/watchers/interfaces"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
)

type WatcherConfig struct {
	NetworkID        watchers.NetworkID // human readable name
	ChainID          vaa.ChainID        // ChainID
	Rpc              string
	Websocket        string
	SuiMoveEventType string
}

func (wc *WatcherConfig) GetNetworkID() watchers.NetworkID {
	return wc.NetworkID
}

func (wc *WatcherConfig) GetChainID() vaa.ChainID {
	return wc.ChainID
}

func (wc *WatcherConfig) RequiredL1Finalizer() watchers.NetworkID {
	return ""
}

func (wc *WatcherConfig) SetL1Finalizer(l1finalizer interfaces.L1Finalizer) {
	// empty
}

func (wc *WatcherConfig) Create(
	msgC chan<- *common.MessagePublication,
	obsvReqC <-chan *gossipv1.ObservationRequest,
	_ chan<- *common.PhylaxSet,
	env common.Environment,
) (interfaces.L1Finalizer, supervisor.Runnable, error) {
	var devMode bool = (env == common.UnsafeDevNet)

	return nil, NewWatcher(wc.Rpc, wc.Websocket, wc.SuiMoveEventType, devMode, msgC, obsvReqC).Run, nil
}
