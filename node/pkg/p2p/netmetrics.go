package p2p

import (
	"math"
	"regexp"
	"strconv"

	gossipv1 "github.com/deltaswapio/deltaswap/node/pkg/proto/gossip/v1"
	"github.com/deltaswapio/deltaswap/node/pkg/version"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	deltaswapNetworkNodeHeight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deltaswap_network_node_height",
			Help: "Network height of the given phylax node per network",
		}, []string{"phylax_addr", "node_id", "node_name", "network"})
	deltaswapNetworkNodeErrors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deltaswap_network_node_errors_count",
			Help: "Number of errors the given phylax node encountered per network",
		}, []string{"phylax_addr", "node_id", "node_name", "network"})
	deltaswapNetworkVersion = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deltaswap_network_node_version",
			Help: "Network version of the given phylax node per network",
		}, []string{"phylax_addr", "node_id", "node_name", "network", "version"})
)

func collectNodeMetrics(addr common.Address, peerId peer.ID, hb *gossipv1.Heartbeat) {
	for _, n := range hb.Networks {
		if n == nil {
			continue
		}

		chain := vaa.ChainID(n.Id)

		deltaswapNetworkNodeHeight.WithLabelValues(
			addr.Hex(), peerId.Pretty(), hb.NodeName, chain.String()).Set(float64(n.Height))

		deltaswapNetworkNodeErrors.WithLabelValues(
			addr.Hex(), peerId.Pretty(), hb.NodeName, chain.String()).Set(float64(n.ErrorCount))

		deltaswapNetworkVersion.WithLabelValues(
			addr.Hex(), peerId.Pretty(), hb.NodeName, chain.String(),
			sanitizeVersion(hb.Version, version.Version())).Set(1)
	}
}

var (
	// Parse version string using regular expression.
	// The version string should be in the format of "vX.Y.Z"
	// where X, Y and Z are integers. Suffixes are ignored.
	reVersion = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)`)
)

// sanitizeVersion cleans up the version string to prevent an attacker from executing a cardinality attack.
func sanitizeVersion(version string, reference string) string {
	// Match groups of reVersion
	components := reVersion.FindStringSubmatch(version)
	referenceComponents := reVersion.FindStringSubmatch(reference)

	// Compare components of the version string with the reference and ensure
	// that the distance is less than 5.
	for i, c := range components {
		if len(referenceComponents) <= i {
			return "other"
		}

		cInt, _ := strconv.Atoi(c)
		cRefInt, _ := strconv.Atoi(referenceComponents[i])

		if math.Abs(float64(cInt-cRefInt)) > 5 {
			return "other"
		}
	}

	v := reVersion.FindString(version)
	if v == "" {
		return "other"
	}
	return v
}
