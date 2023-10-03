package near

import (
	"context"
	"time"

	"github.com/deltaswapio/deltaswap/node/pkg/p2p"
	"github.com/deltaswapio/deltaswap/node/pkg/supervisor"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type eventType int

const (
	EVENT_FINALIZED_CACHE_MISS eventType = iota
	EVENT_NEAR_MESSAGE_CONFIRMED
	EVENT_NEAR_API_HTTP_ERR // NEAR API returned a status code other than 200
	EVENT_NEAR_WATCHER_TOO_FAR_BEHIND
)

func (e *Watcher) runMetrics(ctx context.Context) error {
	logger := supervisor.Logger(ctx)

	reg := prometheus.NewRegistry()

	deltaswapMsgCounter := 0
	var deltaswapMsgTotalTime time.Duration = 0

	deltaswapTxAvgDuration := promauto.With(reg).NewGauge(
		prometheus.GaugeOpts{
			Name: "deltaswap_near_tx_avg_duration",
			Help: "Average duration it takes for a deltaswap message to be processed in milliseconds",
		})

	txqueueLen := promauto.With(reg).NewGauge(
		prometheus.GaugeOpts{
			Name: "deltaswap_near_tx_queue",
			Help: "Current Near transaction processing queue length",
		})

	chunkqueueLen := promauto.With(reg).NewGauge(
		prometheus.GaugeOpts{
			Name: "deltaswap_near_chunk_queue",
			Help: "Current Near chunk processing queue length",
		})

	nearMessagesConfirmed := promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "deltaswap_near_observations_confirmed_total",
			Help: "Total number of verified Near observations found",
		})

	nearFinalizedCacheMisses := promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "deltaswap_near_finalized_cache_misses",
			Help: "Number of times the watcher needed to check the finalization status of a block by walking the chain forward because the finalization status was not cached.",
		})

	nearRpcErrorCounter := promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "deltaswap_near_rpc_error",
			Help: "NEAR RPC Error Counter",
		})

	metricsIntervalTimer := time.NewTicker(metricsInterval) // this is just one ms for the first iteration.

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-metricsIntervalTimer.C:
			// compute and publish periodic metrics
			l1 := e.transactionProcessingQueueCounter.Load()
			l2 := len(e.transactionProcessingQueue)
			l3 := len(e.chunkProcessingQueue)
			txqueueLen.Set(float64(l1))
			chunkqueueLen.Set(float64(l2))
			logger.Debug("metrics", zap.Int64("txqueueWaiting", l1), zap.Int("txqueueReady", l2), zap.Int("chunkqueueLen", l3))

		case event := <-e.eventChan:
			switch event {
			case EVENT_FINALIZED_CACHE_MISS:
				nearFinalizedCacheMisses.Inc()
			case EVENT_NEAR_MESSAGE_CONFIRMED:
				nearMessagesConfirmed.Inc()
			case EVENT_NEAR_WATCHER_TOO_FAR_BEHIND:
				logger.Error("NEAR Watcher fell behind too far", zap.String("log_msg_type", "watcher_behind"))
				p2p.DefaultRegistry.AddErrorCount(vaa.ChainIDNear, 1)
			case EVENT_NEAR_API_HTTP_ERR:
				nearRpcErrorCounter.Inc()
			}
		case d := <-e.eventChanTxProcessedDuration:
			deltaswapMsgCounter++
			deltaswapMsgTotalTime += d
			avgDurationMs := deltaswapMsgTotalTime.Milliseconds() / int64(deltaswapMsgCounter)
			deltaswapTxAvgDuration.Set(float64(avgDurationMs))
		}
	}
}
