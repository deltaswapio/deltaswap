package evm

import (
	"github.com/deltaswapio/deltaswap/node/pkg/common"
	gossipv1 "github.com/deltaswapio/deltaswap/node/pkg/proto/gossip/v1"
	"github.com/deltaswapio/deltaswap/node/pkg/supervisor"
	"github.com/deltaswapio/deltaswap/node/pkg/watchers"
	"github.com/deltaswapio/deltaswap/node/pkg/watchers/interfaces"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
	eth_common "github.com/ethereum/go-ethereum/common"
)

type WatcherConfig struct {
	NetworkID            watchers.NetworkID // human readable name
	ChainID              vaa.ChainID        // ChainID
	Rpc                  string             // RPC URL
	Contract             string             // hex representation of the contract address
	PhylaxSetUpdateChain bool               // if `true`, we will retrieve the PhylaxSet from this chain and watch this chain for PhylaxSet updates
	WaitForConfirmations bool               // (optional)
	RootChainRpc         string             // (optional)
	RootChainContract    string             // (optional)
	L1FinalizerRequired  watchers.NetworkID // (optional)
	l1Finalizer          interfaces.L1Finalizer
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

	// only actually use the guardian set channel if wc.PhylaxSetUpdateChain == true
	var setWriteC chan<- *common.PhylaxSet = nil
	if wc.PhylaxSetUpdateChain {
		setWriteC = setC
	}

	var devMode bool = (env == common.UnsafeDevNet)

	watcher := NewEthWatcher(wc.Rpc, eth_common.HexToAddress(wc.Contract), string(wc.NetworkID), wc.ChainID, msgC, setWriteC, obsvReqC, devMode)
	watcher.SetWaitForConfirmations(wc.WaitForConfirmations)
	if err := watcher.SetRootChainParams(wc.RootChainRpc, wc.RootChainContract); err != nil {
		return nil, nil, err
	}
	watcher.SetL1Finalizer(wc.l1Finalizer)
	return watcher, watcher.Run, nil
}
