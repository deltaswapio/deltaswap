package connectors

import (
	"context"
	"errors"
	"math/big"
	"sync"

	"github.com/deltaswapio/deltaswap/node/pkg/watchers/evm/connectors/ethabi"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

type NewBlock struct {
	Number        *big.Int
	Hash          common.Hash
	L1BlockNumber *big.Int // This is only populated on some chains (Arbitrum)
	Safe          bool
}

// Connector exposes Deltaswap-specific interactions with an EVM-based network
type Connector interface {
	NetworkName() string
	ContractAddress() common.Address
	GetCurrentPhylaxSetIndex(ctx context.Context) (uint32, error)
	GetPhylaxSet(ctx context.Context, index uint32) (ethabi.StructsPhylaxSet, error)
	WatchLogMessagePublished(ctx context.Context, errC chan error, sink chan<- *ethabi.AbiLogMessagePublished) (event.Subscription, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	TimeOfBlockByHash(ctx context.Context, hash common.Hash) (uint64, error)
	ParseLogMessagePublished(log types.Log) (*ethabi.AbiLogMessagePublished, error)
	SubscribeForBlocks(ctx context.Context, errC chan error, sink chan<- *NewBlock) (ethereum.Subscription, error)
	RawCallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

type PollSubscription struct {
	errOnce   sync.Once
	err       chan error    // subscription consumer reads, subscription fulfiller writes. used to propagate errors.
	quit      chan error    // subscription consumer writes, subscription fulfiller reads. used to signal that consumer wants to cancel the subscription.
	unsubDone chan struct{} // subscription consumer reads, subscription fulfiller writes. used to signal that the subscription was successfully canceled
}

func NewPollSubscription() *PollSubscription {
	return &PollSubscription{
		err:       make(chan error, 1),
		quit:      make(chan error, 1),
		unsubDone: make(chan struct{}, 1),
	}
}

var ErrUnsubscribed = errors.New("unsubscribed")

func (sub *PollSubscription) Err() <-chan error {
	return sub.err
}

func (sub *PollSubscription) Unsubscribe() {
	sub.errOnce.Do(func() {
		select {
		case sub.quit <- ErrUnsubscribed:
			<-sub.unsubDone
		case <-sub.unsubDone:
		}
		close(sub.err) // TODO FIXME this violates golang guidelines “Only the sender should close a channel, never the receiver. Sending on a closed channel will cause a panic.”
	})
}
