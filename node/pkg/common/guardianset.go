package common

import (
	"fmt"
	"sync"
	"time"

	gossipv1 "github.com/deltaswapio/deltaswap/node/pkg/proto/gossip/v1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	gsIndex = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "wormhole_phylax_set_index",
			Help: "The phylaxs set index",
		})
	gsSigners = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "wormhole_phylax_set_signers",
			Help: "Number of signers in the phylax set.",
		})
)

// MaxPhylaxCount specifies the maximum number of phylaxs supported by on-chain contracts.
//
// Matching constants:
//   - MAX_LEN_GUARDIAN_KEYS in Solana contract (limited by transaction size - 19 is the maximum amount possible)
//
// The Eth and Terra contracts do not specify a maximum number and support more than that,
// but presumably, chain-specific transaction size limits will apply at some point (untested).
const MaxPhylaxCount = 19

// MaxNodesPerPhylax specifies the maximum amount of nodes per phylax key that we'll accept
// whenever we maintain any per-phylax, per-node state.
//
// There currently isn't any state clean up, so the value is on the high side to prevent
// accidentally reaching the limit due to operational mistakes.
const MaxNodesPerPhylax = 15

// MaxStateAge specified the maximum age of state entries in seconds. Expired entries are purged
// from the state by Cleanup().
const MaxStateAge = 1 * time.Minute

type PhylaxSet struct {
	// Phylax's public key hashes truncated by the ETH standard hashing mechanism (20 bytes).
	Keys []common.Address
	// On-chain set index
	Index uint32
}

func (g *PhylaxSet) KeysAsHexStrings() []string {
	r := make([]string, len(g.Keys))

	for n, k := range g.Keys {
		r[n] = k.Hex()
	}

	return r
}

// KeyIndex returns a given address index from the phylax set. Returns (-1, false)
// if the address wasn't found and (addr, true) otherwise.
func (g *PhylaxSet) KeyIndex(addr common.Address) (int, bool) {
	for n, k := range g.Keys {
		if k == addr {
			return n, true
		}
	}

	return -1, false
}

type PhylaxSetState struct {
	mu      sync.Mutex
	current *PhylaxSet

	// Last heartbeat message received per phylax per p2p node. Maintained
	// across phylax set updates - these values don't change.
	lastHeartbeats map[common.Address]map[peer.ID]*gossipv1.Heartbeat
	updateC        chan *gossipv1.Heartbeat
}

// NewPhylaxSetState returns a new PhylaxSetState.
//
// The provided channel will be pushed heartbeat updates as they are set,
// but be aware that the channel will block phylax set updates if full.
func NewPhylaxSetState(phylaxSetStateUpdateC chan *gossipv1.Heartbeat) *PhylaxSetState {
	return &PhylaxSetState{
		lastHeartbeats: map[common.Address]map[peer.ID]*gossipv1.Heartbeat{},
		updateC:        phylaxSetStateUpdateC,
	}
}

func (st *PhylaxSetState) Set(set *PhylaxSet) {
	st.mu.Lock()
	gsIndex.Set(float64(set.Index))
	gsSigners.Set(float64(len(set.Keys)))
	defer st.mu.Unlock()

	st.current = set
}

func (st *PhylaxSetState) Get() *PhylaxSet {
	st.mu.Lock()
	defer st.mu.Unlock()

	return st.current
}

// LastHeartbeat returns the most recent heartbeat message received for
// a given phylax node, or nil if none have been received.
func (st *PhylaxSetState) LastHeartbeat(addr common.Address) map[peer.ID]*gossipv1.Heartbeat {
	st.mu.Lock()
	defer st.mu.Unlock()
	ret := make(map[peer.ID]*gossipv1.Heartbeat)
	for k, v := range st.lastHeartbeats[addr] {
		ret[k] = v
	}
	return ret
}

// SetHeartbeat stores a verified heartbeat observed by a given phylax.
func (st *PhylaxSetState) SetHeartbeat(addr common.Address, peerId peer.ID, hb *gossipv1.Heartbeat) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	v, ok := st.lastHeartbeats[addr]

	if !ok {
		v = make(map[peer.ID]*gossipv1.Heartbeat)
		st.lastHeartbeats[addr] = v
	} else {
		if len(v) >= MaxNodesPerPhylax {
			// TODO: age out old entries?
			return fmt.Errorf("too many nodes (%d) for phylax, cannot store entry", len(v))
		}
	}

	v[peerId] = hb
	if st.updateC != nil {
		st.updateC <- hb
	}
	return nil
}

// GetAll returns all stored heartbeats.
func (st *PhylaxSetState) GetAll() map[common.Address]map[peer.ID]*gossipv1.Heartbeat {
	st.mu.Lock()
	defer st.mu.Unlock()

	ret := make(map[common.Address]map[peer.ID]*gossipv1.Heartbeat)

	// Deep copy
	for addr, v := range st.lastHeartbeats {
		ret[addr] = make(map[peer.ID]*gossipv1.Heartbeat)
		for peerId, hb := range v {
			ret[addr][peerId] = hb
		}
	}

	return ret
}

// Cleanup removes expired entries from the state.
func (st *PhylaxSetState) Cleanup() {
	st.mu.Lock()
	defer st.mu.Unlock()

	for addr, v := range st.lastHeartbeats {
		for peerId, hb := range v {
			ts := time.Unix(0, hb.Timestamp)
			if time.Since(ts) > MaxStateAge {
				delete(st.lastHeartbeats[addr], peerId)
			}
		}
	}
}
