package cosmwasm

import (
	"github.com/deltaswapio/deltaswap/node/pkg/common"
	gossipv1 "github.com/deltaswapio/deltaswap/node/pkg/proto/gossip/v1"
	"github.com/deltaswapio/deltaswap/node/pkg/query"
	"github.com/deltaswapio/deltaswap/node/pkg/supervisor"
	"github.com/deltaswapio/deltaswap/node/pkg/watchers"
	"github.com/deltaswapio/deltaswap/node/pkg/watchers/interfaces"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
)

type WatcherConfig struct {
	NetworkID watchers.NetworkID // human readable name
	ChainID   vaa.ChainID        // ChainID
	Websocket string             // Websocket URL
	Lcd       string             // LCD
	Contract  string             // hex representation of the contract address
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
	_ <-chan *query.PerChainQueryInternal,
	_ chan<- *query.PerChainQueryResponseInternal,
	_ chan<- *common.PhylaxSet,
	env common.Environment,
) (interfaces.L1Finalizer, supervisor.Runnable, error) {
	return nil, NewWatcher(wc.Websocket, wc.Lcd, wc.Contract, msgC, obsvReqC, wc.ChainID, env).Run, nil
}
