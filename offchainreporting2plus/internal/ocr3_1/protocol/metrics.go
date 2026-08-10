package protocol

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
	"github.com/smartcontractkit/libocr/internal/util"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

type attestedStateTransitionBlockWriter string

const (
	attestedStateTransitionBlockWriterOutcomeGeneration attestedStateTransitionBlockWriter = "outcome_generation"
	attestedStateTransitionBlockWriterStateSync         attestedStateTransitionBlockWriter = "state_sync"
)

// newAttestedStateTransitionBlocksWrittenTotal is shared by
// outcomeGenerationMetrics and stateSyncMetrics. The two variants are
// registered from separate structs against the same registerer, and prometheus
// only tolerates that if the name, help and label names match exactly.
func newAttestedStateTransitionBlocksWrittenTotal(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
	writer attestedStateTransitionBlockWriter,
) prometheus.Counter {
	c := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_experimental_attested_state_transition_blocks_written_total",
		Help: fmt.Sprintf("Total number of attested state transition blocks written with writer being one of %v",
			[]attestedStateTransitionBlockWriter{
				attestedStateTransitionBlockWriterOutcomeGeneration, attestedStateTransitionBlockWriterStateSync}),
		ConstLabels: prometheus.Labels{"writer": string(writer)},
	})
	metricshelper.RegisterOrLogError(logger, registerer, c, "ocr3_1_experimental_attested_state_transition_blocks_written_total")
	return c
}

type pacemakerMetrics struct {
	registerer             prometheus.Registerer
	epoch                  prometheus.Gauge
	leader                 prometheus.Gauge
	tProgressTimeoutsTotal prometheus.Counter
}

func newPacemakerMetrics(registerer prometheus.Registerer,
	logger commontypes.Logger) *pacemakerMetrics {

	epoch := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_epoch",
		Help: "The total number of initialized epochs",
	})
	metricshelper.RegisterOrLogError(logger, registerer, epoch, "ocr3_1_epoch")

	leader := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_experimental_leader_oid",
		Help: "The leader oracle id",
	})
	metricshelper.RegisterOrLogError(logger, registerer, leader, "ocr3_1_experimental_leader_oid")

	tProgressTimeoutsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_experimental_t_progress_timeouts_total",
		Help: "Total number of TProgress timeouts",
	})
	metricshelper.RegisterOrLogError(logger, registerer, tProgressTimeoutsTotal, "ocr3_1_experimental_t_progress_timeouts_total")

	return &pacemakerMetrics{
		registerer,
		epoch,
		leader,
		tProgressTimeoutsTotal,
	}
}

func (pm *pacemakerMetrics) Close() {
	pm.registerer.Unregister(pm.epoch)
	pm.registerer.Unregister(pm.leader)
	pm.registerer.Unregister(pm.tProgressTimeoutsTotal)
}

type outcomeGenerationMetrics struct {
	registerer                 prometheus.Registerer
	committedSeqNr             prometheus.Gauge
	sentObservationsTotal      prometheus.Counter
	includedObservationsTotal  prometheus.Counter
	ledCommittedRoundsTotal    prometheus.Counter
	attestedBlocksWrittenTotal prometheus.Counter
}

func newOutcomeGenerationMetrics(registerer prometheus.Registerer,
	logger commontypes.Logger) *outcomeGenerationMetrics {

	committedSeqNr := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_committed_sequence_number",
		Help: "The committed sequence number",
	})
	metricshelper.RegisterOrLogError(logger, registerer, committedSeqNr, "ocr3_1_committed_sequence_number")

	sentObservationsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_sent_observations_total",
		Help: "The total number of observations by this oracle sent to the leader. Note that a " +
			"sent observation might not arrive at the leader in time, or not be included in a " +
			"proposal for other reasons. This metric is useful for checking an oracle's ability " +
			"to make observations.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, sentObservationsTotal, "ocr3_1_sent_observations_total")

	includedObservationsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_included_observations_total",
		Help: "The total number of (valid) observations by this oracle included in a proposal " +
			"from the leader. Note that there is no guarantee that the proposal will actually get " +
			"committed; for instance, because the leader crashes or maliciously equivocates to " +
			"make this oracle believe that an observation was included. This metric is useful in " +
			"comparison with ocr3_1_sent_observations_total to check whether an oracle is able to " +
			"regularly make observations that are included in proposals.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, includedObservationsTotal, "ocr3_1_included_observations_total")

	ledCommittedRoundsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_led_committed_rounds_total",
		Help: "The total number of rounds committed that were led by this oracle. This metric is " +
			"useful for checking an oracle's ability to act as leader.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, ledCommittedRoundsTotal, "ocr3_1_led_committed_rounds_total")

	attestedBlocksWrittenTotal := newAttestedStateTransitionBlocksWrittenTotal(registerer, logger, attestedStateTransitionBlockWriterOutcomeGeneration)

	return &outcomeGenerationMetrics{
		registerer,
		committedSeqNr,
		sentObservationsTotal,
		includedObservationsTotal,
		ledCommittedRoundsTotal,
		attestedBlocksWrittenTotal,
	}
}

func (om *outcomeGenerationMetrics) Close() {
	om.registerer.Unregister(om.committedSeqNr)
	om.registerer.Unregister(om.sentObservationsTotal)
	om.registerer.Unregister(om.includedObservationsTotal)
	om.registerer.Unregister(om.ledCommittedRoundsTotal)
	om.registerer.Unregister(om.attestedBlocksWrittenTotal)
}

type stateSyncMetrics struct {
	registerer                  prometheus.Registerer
	attestedBlocksWrittenTotal  prometheus.Counter
	attestedBlocksReplayedTotal prometheus.Counter
	treeSyncCompletedTotal      prometheus.Counter
}

func newStateSyncMetrics(registerer prometheus.Registerer,
	logger commontypes.Logger) *stateSyncMetrics {

	attestedBlocksWrittenTotal := newAttestedStateTransitionBlocksWrittenTotal(registerer, logger, attestedStateTransitionBlockWriterStateSync)

	attestedBlocksReplayedTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_experimental_attested_state_transition_blocks_replayed_total",
		Help: "Total number of attested state transition blocks replayed",
	})
	metricshelper.RegisterOrLogError(logger, registerer, attestedBlocksReplayedTotal, "ocr3_1_experimental_attested_state_transition_blocks_replayed_total")

	treeSyncCompletedTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_experimental_tree_sync_completed_total",
		Help: "Total number of tree synchronizations to a snapshot successfully completed",
	})
	metricshelper.RegisterOrLogError(logger, registerer, treeSyncCompletedTotal, "ocr3_1_experimental_tree_sync_completed_total")

	return &stateSyncMetrics{
		registerer,
		attestedBlocksWrittenTotal,
		attestedBlocksReplayedTotal,
		treeSyncCompletedTotal,
	}
}

func (sm *stateSyncMetrics) Close() {
	sm.registerer.Unregister(sm.attestedBlocksWrittenTotal)
	sm.registerer.Unregister(sm.attestedBlocksReplayedTotal)
	sm.registerer.Unregister(sm.treeSyncCompletedTotal)
}

type blobExchangeMetrics struct {
	registerer              prometheus.Registerer
	blobsInProgress         prometheus.Gauge
	myBroadcastPayloadBytes prometheus.Histogram
}

// 256 bytes, 1KiB, ..., 256MiB, copied from ragep2p_experimental_peer_message_bytes
var blobPayloadBytesBuckets = []float64{1 << 8, 1 << 10, 1 << 12, 1 << 14, 1 << 16, 1 << 18, 1 << 20, 1 << 22, 1 << 24, 1 << 26, 1 << 28}

func newBlobExchangeMetrics(registerer prometheus.Registerer,
	logger commontypes.Logger) *blobExchangeMetrics {

	blobsInProgress := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_experimental_blobs_in_progress",
		Help: "The number of blobs that are actively being broadcast or fetched",
	})
	metricshelper.RegisterOrLogError(logger, registerer, blobsInProgress, "ocr3_1_experimental_blobs_in_progress")

	myBroadcastPayloadBytes := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ocr3_1_blob_my_broadcast_payload_bytes",
		Help:    "The payload size in bytes of each outbound blob broadcast initiated by this oracle",
		Buckets: blobPayloadBytesBuckets,
	})
	metricshelper.RegisterOrLogError(logger, registerer, myBroadcastPayloadBytes, "ocr3_1_blob_my_broadcast_payload_bytes")

	return &blobExchangeMetrics{
		registerer,
		blobsInProgress,
		myBroadcastPayloadBytes,
	}
}

func (bm *blobExchangeMetrics) Close() {
	bm.registerer.Unregister(bm.blobsInProgress)
	bm.registerer.Unregister(bm.myBroadcastPayloadBytes)
}

type blobOracleMetrics struct {
	registerer prometheus.Registerer

	// Repeat broadcasts after the initial broadcast has been pruned from bex
	// state will be treated as new broadcasts when it comes to metrics.

	// Semantics: Receiving a MessageBlobOffer with the same contents twice
	// increments the counter twice. Broadcaster feels like there's some blob you
	// should accept/reject still.
	// EXPERIMENTAL
	theirBlobOffersTotal prometheus.Counter

	// Semantics: Observe when we add to the append stats for the first time
	// EXPERIMENTAL
	theirBlobAppendedPayloadBytes prometheus.Histogram

	// Semantics: Because of offer resends these are not comparable with
	// theirBlobsOfferedTotal. Counting the number of MesssageBlobOfferResponses
	// with reject=true/false sent. In typical operation and honest broadcaster, this
	// does not overaccount, but might overaccount otherwise.
	//
	// Usage: Should compare with historical trends, instant value is not very
	// meaningful in isolation. e.g., rejectsSentTotal might reasonably be > 0
	// in steady state operation.
	//
	// EXPERIMENTAL
	theirBlobAcceptsSentTotal prometheus.Counter
	// EXPERIMENTAL
	theirBlobRejectsSentTotal prometheus.Counter

	// Usage: Compare with historical trends, instant values not that meaningful.
	// EXPERIMENTAL
	myBlobsAcceptedTotal prometheus.Counter
	// If you are debugging blob broadcast issues, this metric can tell you
	// whether a given recipient oracle is rejecting blobs at an abnormally high
	// rate, which could point to an issue with quotas or rate limits, or
	// expiration being too tight.
	// EXPERIMENTAL
	myBlobsRejectedTotal prometheus.Counter
	// EXPERIMENTAL
	myBlobsUndeterminedTotal prometheus.Counter

	// Semantics: If broadcaster crashes, is adversarial, or re-broadcasts the same
	// blob and we re-reject, we will overaccount this metric. An honest oracle
	// receiving a rejection for a single BroadcastBlob will not ask us again,
	// thus in the typical case there is no overaccounting.
	//
	// Usage: Instant values are also meaningful, e.g.,
	// theirBlobsRejectedDueToOversizePayloadBytesTotal > 0 is already
	// indicative of a bug.
	theirBlobsRejectedDueToOversizePayloadBytesTotal prometheus.Counter
	theirBlobsRejectedDueToQuotaTotal                prometheus.Counter
	theirBlobsRejectedDueToManyInflightTotal         prometheus.Counter
	theirBlobsRejectedDueToExpirationTotal           prometheus.Counter

	// Stats metrics are only meant to be set through SetTheirAppendedStats and
	// SetTheirReapedStats, hence their underscore prefix.

	_statsMetrics *blobOracleStatsMetrics
}

type blobOracleStatsMetrics struct {
	// EXPERIMENTAL
	theirBlobQuotaAppendedCount *metricshelper.LazyGauge
	// EXPERIMENTAL
	theirBlobQuotaAppendedCumulativePayloadBytes *metricshelper.LazyGauge
	// EXPERIMENTAL
	theirBlobQuotaReapedCount *metricshelper.LazyGauge
	// EXPERIMENTAL
	theirBlobQuotaReapedCumulativePayloadBytes *metricshelper.LazyGauge

	theirBlobQuotaUsedCount                  *metricshelper.LazyGauge
	theirBlobQuotaUsedCumulativePayloadBytes *metricshelper.LazyGauge

	theirBlobQuotaFreeCount                  *metricshelper.LazyGauge
	theirBlobQuotaFreeCumulativePayloadBytes *metricshelper.LazyGauge

	// Need for computation of free quota
	limitQuotaStats BlobQuotaStats

	mu       sync.Mutex
	appended *BlobQuotaStats
	reaped   *BlobQuotaStats
}

func newBlobOracleMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
	remote commontypes.OracleID,
	pluginLimits ocr3_1types.ReportingPluginLimits,
) *blobOracleMetrics {
	broadcasterOracleIDLabels := prometheus.Labels{"broadcaster_oracle_id": fmt.Sprint(remote)}
	fetcherOracleIDLabels := prometheus.Labels{"fetcher_oracle_id": fmt.Sprint(remote)}

	newLazyGauge := func(name string, help string, labels prometheus.Labels) *metricshelper.LazyGauge {
		return metricshelper.NewLazyGauge(registerer, logger, prometheus.GaugeOpts{
			Name:        name,
			Help:        help,
			ConstLabels: labels,
		})
	}

	newCounter := func(name string, help string, labels prometheus.Labels) prometheus.Counter {
		c := prometheus.NewCounter(prometheus.CounterOpts{
			Name:        name,
			Help:        help,
			ConstLabels: labels,
		})
		metricshelper.RegisterOrLogError(logger, registerer, c, name)
		return c
	}

	theirBlobOffersTotal := newCounter(
		"ocr3_1_experimental_their_blob_offers_total",
		"Number of blob offers received by the broadcaster oracle. Multiple offers count multiple times, "+
			"including in the case where the oracle is fetching chunks for a previous identical offer, but the same offer is resent.",
		broadcasterOracleIDLabels,
	)

	theirBlobAppendedPayloadBytes := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:        "ocr3_1_experimental_their_blob_appended_payload_bytes",
		Help:        "The payload size in bytes of inbound blob broadcasts received from the broadcaster oracle",
		ConstLabels: broadcasterOracleIDLabels,
		Buckets:     blobPayloadBytesBuckets,
	})
	metricshelper.RegisterOrLogError(logger, registerer, theirBlobAppendedPayloadBytes, "ocr3_1_experimental_their_blob_appended_payload_bytes")

	type theirBroadcastResponse string
	const (
		theirBroadcastResponseAccept theirBroadcastResponse = "accept"
		theirBroadcastResponseReject theirBroadcastResponse = "reject"
	)

	newTheirBlobOfferResponsesTotal := func(t theirBroadcastResponse) prometheus.Counter {
		return newCounter(
			"ocr3_1_experimental_their_blob_offer_responses_total",
			fmt.Sprintf("The total number of inbound blob broadcasts received from the broadcaster oracle, distinguished by response label to indicate what we responded. "+
				"Usage: Should compare with historical trends, instant value is "+
				"not very meaningful in isolation. e.g., %s might reasonably "+
				"be > 0 in steady state operation. "+
				"Valid responses are: %s (i.e., we sent an accepting offer response), "+
				"%s (i.e., we sent a rejecting offer response, either due to outright rejecting the offer or the blob expiring before we could accept)",
				theirBroadcastResponseReject,
				theirBroadcastResponseAccept, theirBroadcastResponseReject),
			util.MapsUnion(
				broadcasterOracleIDLabels,
				prometheus.Labels{"response": string(t)},
			),
		)
	}

	theirBlobAcceptsSentTotal := newTheirBlobOfferResponsesTotal(theirBroadcastResponseAccept)
	theirBlobRejectsSentTotal := newTheirBlobOfferResponsesTotal(theirBroadcastResponseReject)

	type myBroadcastResult string
	const (
		myBroadcastResultAccepted     myBroadcastResult = "accepted"
		myBroadcastResultRejected     myBroadcastResult = "rejected"
		myBroadcastResultUndetermined myBroadcastResult = "undetermined"
	)

	newMyBlobsTotal := func(m myBroadcastResult) prometheus.Counter {
		return newCounter(
			"ocr3_1_experimental_my_blobs_total",
			fmt.Sprintf("The total number of outbound blob broadcasts initiated by us to the fetcher oracle, distinguished by result label. "+
				"Usage: If you are debugging blob broadcast issues, this metric can tell you "+
				"whether a given recipient oracle is rejecting blobs at an abnormally high "+
				"rate, which could point to an issue with quotas or rate limits, or "+
				"expiration being too tight. "+
				"Valid results are: %s (i.e., the fetcher oracle sent us an accepting offer response), "+
				"%s (i.e., the fetcher oracle sent us a rejecting offer response), or "+
				"%s (i.e., the fetcher oracle did not send us an offer response before we gave up on the broadcast). "+
				"Note that once we receive enough accepting or rejecting responses we stop processing responses, "+
				"thus it is expected for a number of oracles that did not respond in time due to being slower to have undetermined result.",
				myBroadcastResultAccepted, myBroadcastResultRejected, myBroadcastResultUndetermined),
			util.MapsUnion(
				fetcherOracleIDLabels,
				prometheus.Labels{"result": string(m)},
			),
		)
	}

	myBlobsAcceptedTotal := newMyBlobsTotal(myBroadcastResultAccepted)
	myBlobsRejectedTotal := newMyBlobsTotal(myBroadcastResultRejected)
	myBlobsUndeterminedTotal := newMyBlobsTotal(myBroadcastResultUndetermined)

	type rejectionReason string
	const (
		rejectionReasonOversizePayloadBytes rejectionReason = "oversize_payload_bytes"
		rejectionReasonQuota                rejectionReason = "quota"
		rejectionReasonTooManyInflight      rejectionReason = "too_many_inflight"
		rejectionReasonExpiration           rejectionReason = "expiration"
	)

	newTheirBlobsRejectedTotal := func(r rejectionReason) prometheus.Counter {
		return newCounter(
			"ocr3_1_their_blobs_rejected_total",
			fmt.Sprintf("The total number of blobs that the broadcaster oracle offered to us, distinguished by the reason label for why we rejected them. "+
				"Usage: Instant values are also meaningful, e.g., "+
				"%s > 0 is already "+
				"indicative of a bug. "+
				"Valid reasons are: %s (i.e., the blob had oversize payload bytes compared to MaxBlobPayloadBytes), "+
				"%s (i.e., the blob went over the defined quota in MaxPerOracleUnexpiredBlobCumulativePayloadBytes, MaxPerOracleUnexpiredBlobCount, see also ocr3_1_blob_quota* metrics), "+
				"%s (i.e., too many blobs were inflight already from this broadcaster oracle), "+
				"%s (i.e., the blob was expired either immediately upon offer receipt or before the time we could fetch all chunks and accept it).",
				rejectionReasonOversizePayloadBytes,
				rejectionReasonOversizePayloadBytes, rejectionReasonQuota, rejectionReasonTooManyInflight, rejectionReasonExpiration,
			),
			util.MapsUnion(
				broadcasterOracleIDLabels,
				prometheus.Labels{"reason": string(r)},
			),
		)
	}

	theirBlobsRejectedDueToOversizePayloadBytesTotal := newTheirBlobsRejectedTotal(rejectionReasonOversizePayloadBytes)
	theirBlobsRejectedDueToQuotaTotal := newTheirBlobsRejectedTotal(rejectionReasonQuota)
	theirBlobsRejectedDueToManyInflightTotal := newTheirBlobsRejectedTotal(rejectionReasonTooManyInflight)
	theirBlobsRejectedDueToExpirationTotal := newTheirBlobsRejectedTotal(rejectionReasonExpiration)

	theirBlobQuotaAppendedCount := newLazyGauge(
		"ocr3_1_experimental_blob_quota_appended_count",
		"The number of blobs we ever started fetching from the broadcaster oracle. "+
			"This metric is not meaningful for broadcaster=self, as we do not count our own broadcasts against our quota.",
		broadcasterOracleIDLabels,
	)
	theirBlobQuotaAppendedCumulativePayloadBytes := newLazyGauge(
		"ocr3_1_experimental_blob_quota_appended_cumulative_payload_bytes",
		"The cumulative payload bytes of blobs we ever started fetching from the broadcaster oracle. "+
			"This metric is not meaningful for broadcaster=self, as we do not count our own broadcasts against our quota.",
		broadcasterOracleIDLabels,
	)
	theirBlobQuotaReapedCount := newLazyGauge(
		"ocr3_1_experimental_blob_quota_reaped_count",
		"The number of blobs that we ever started fetching from the broadcaster oracle, "+
			"that have since been reaped.",
		broadcasterOracleIDLabels,
	)
	theirBlobQuotaReapedCumulativePayloadBytes := newLazyGauge(
		"ocr3_1_experimental_blob_quota_reaped_cumulative_payload_bytes",
		"The cumulative payload bytes of blobs that we ever started fetching from the broadcaster oracle, "+
			"that have since been reaped.",
		broadcasterOracleIDLabels,
	)

	theirBlobQuotaUsedCount := newLazyGauge(
		"ocr3_1_blob_quota_used_count",
		"The number of blobs that we ever started fetching from the broadcaster oracle, "+
			"that have not yet been reaped. Must be less than or equal to MaxPerOracleUnexpiredBlobCount.",
		broadcasterOracleIDLabels,
	)
	theirBlobQuotaUsedCumulativePayloadBytes := newLazyGauge(
		"ocr3_1_blob_quota_used_cumulative_payload_bytes",
		"The cumulative payload bytes of blobs that we ever started fetching from the broadcaster oracle, "+
			"that have not yet been reaped. Must be less than or equal to MaxPerOracleUnexpiredBlobCumulativePayloadBytes.",
		broadcasterOracleIDLabels,
	)

	theirBlobQuotaFreeCount := newLazyGauge(
		"ocr3_1_blob_quota_free_count",
		"The free blob count in the broadcaster oracle's quota. The closer this value is to zero, the more likely it is "+
			"that the broadcaster oracle's offers might be rejected in the future due to insufficient free quota.",
		broadcasterOracleIDLabels,
	)
	theirBlobQuotaFreeCumulativePayloadBytes := newLazyGauge(
		"ocr3_1_blob_quota_free_cumulative_payload_bytes",
		"The free cumulative payload bytes in the broadcaster oracle's quota. The closer this value is to zero, the more likely it is "+
			"that the broadcaster oracle's offers might be rejected in the future due to insufficient free quota.",
		broadcasterOracleIDLabels,
	)

	return &blobOracleMetrics{
		registerer,
		theirBlobOffersTotal,
		theirBlobAppendedPayloadBytes,
		theirBlobAcceptsSentTotal,
		theirBlobRejectsSentTotal,
		myBlobsAcceptedTotal,
		myBlobsRejectedTotal,
		myBlobsUndeterminedTotal,
		theirBlobsRejectedDueToOversizePayloadBytesTotal,
		theirBlobsRejectedDueToQuotaTotal,
		theirBlobsRejectedDueToManyInflightTotal,
		theirBlobsRejectedDueToExpirationTotal,
		&blobOracleStatsMetrics{
			theirBlobQuotaAppendedCount,
			theirBlobQuotaAppendedCumulativePayloadBytes,
			theirBlobQuotaReapedCount,
			theirBlobQuotaReapedCumulativePayloadBytes,
			theirBlobQuotaUsedCount,
			theirBlobQuotaUsedCumulativePayloadBytes,
			theirBlobQuotaFreeCount,
			theirBlobQuotaFreeCumulativePayloadBytes,

			BlobQuotaStats{
				uint64(pluginLimits.MaxPerOracleUnexpiredBlobCount),
				uint64(pluginLimits.MaxPerOracleUnexpiredBlobCumulativePayloadBytes),
			},
			sync.Mutex{},
			nil, nil,
		},
	}
}

func (bosm *blobOracleStatsMetrics) SetTheirAppendedStats(stats BlobQuotaStats) {
	bosm.theirBlobQuotaAppendedCount.Set(float64(stats.Count))
	bosm.theirBlobQuotaAppendedCumulativePayloadBytes.Set(float64(stats.CumulativePayloadLength))

	bosm.mu.Lock()
	defer bosm.mu.Unlock()
	bosm.appended = &stats
	bosm.updateUsedFreeGauges()
}

func (bosm *blobOracleStatsMetrics) SetTheirReapedStats(stats BlobQuotaStats) {
	bosm.theirBlobQuotaReapedCount.Set(float64(stats.Count))
	bosm.theirBlobQuotaReapedCumulativePayloadBytes.Set(float64(stats.CumulativePayloadLength))

	bosm.mu.Lock()
	defer bosm.mu.Unlock()
	bosm.reaped = &stats
	bosm.updateUsedFreeGauges()
}

func (bosm *blobOracleStatsMetrics) updateUsedFreeGauges() {
	if bosm.appended == nil || bosm.reaped == nil {
		return
	}

	used, ok := bosm.appended.Sub(*bosm.reaped)
	if !ok {

		used = BlobQuotaStats{}
	}
	bosm.theirBlobQuotaUsedCount.Set(float64(used.Count))
	bosm.theirBlobQuotaUsedCumulativePayloadBytes.Set(float64(used.CumulativePayloadLength))

	free := bosm.limitQuotaStats.SaturatingSub(used)
	bosm.theirBlobQuotaFreeCount.Set(float64(free.Count))
	bosm.theirBlobQuotaFreeCumulativePayloadBytes.Set(float64(free.CumulativePayloadLength))
}

func (bosm *blobOracleStatsMetrics) Close() {
	bosm.theirBlobQuotaAppendedCount.Unregister()
	bosm.theirBlobQuotaAppendedCumulativePayloadBytes.Unregister()
	bosm.theirBlobQuotaReapedCount.Unregister()
	bosm.theirBlobQuotaReapedCumulativePayloadBytes.Unregister()
	bosm.theirBlobQuotaUsedCount.Unregister()
	bosm.theirBlobQuotaUsedCumulativePayloadBytes.Unregister()
	bosm.theirBlobQuotaFreeCount.Unregister()
	bosm.theirBlobQuotaFreeCumulativePayloadBytes.Unregister()
}

func (bom *blobOracleMetrics) SetTheirAppendedStats(stats BlobQuotaStats) {
	bom._statsMetrics.SetTheirAppendedStats(stats)
}

func (bom *blobOracleMetrics) SetTheirReapedStats(stats BlobQuotaStats) {
	bom._statsMetrics.SetTheirReapedStats(stats)
}

func (bom *blobOracleMetrics) Close() {
	bom.registerer.Unregister(bom.theirBlobOffersTotal)
	bom.registerer.Unregister(bom.theirBlobAppendedPayloadBytes)
	bom.registerer.Unregister(bom.theirBlobAcceptsSentTotal)
	bom.registerer.Unregister(bom.theirBlobRejectsSentTotal)
	bom.registerer.Unregister(bom.myBlobsAcceptedTotal)
	bom.registerer.Unregister(bom.myBlobsRejectedTotal)
	bom.registerer.Unregister(bom.myBlobsUndeterminedTotal)
	bom.registerer.Unregister(bom.theirBlobsRejectedDueToOversizePayloadBytesTotal)
	bom.registerer.Unregister(bom.theirBlobsRejectedDueToQuotaTotal)
	bom.registerer.Unregister(bom.theirBlobsRejectedDueToManyInflightTotal)
	bom.registerer.Unregister(bom.theirBlobsRejectedDueToExpirationTotal)
	bom._statsMetrics.Close()
}

type reportingPluginInfo1Metrics struct {
	registerer prometheus.Registerer
	info       prometheus.Gauge
	limits     *reportingPluginLimitsMetrics
}

func newReportingPluginInfo1Metrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
	reportingPluginInfo ocr3_1types.ReportingPluginInfo1,
) *reportingPluginInfo1Metrics {
	info := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_reporting_plugin_info",
		Help: "Exposes the ReportingPluginInfo1.Name field via labels and has value of 1. See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types#ReportingPluginInfo1 for details",
		ConstLabels: prometheus.Labels{
			"name": reportingPluginInfo.Name,
		},
	})
	info.Set(1)
	metricshelper.RegisterOrLogError(logger, registerer, info, "ocr3_1_reporting_plugin_info")

	return &reportingPluginInfo1Metrics{
		registerer,
		info,
		newReportingPluginLimitsMetrics(registerer, logger, reportingPluginInfo.Limits),
	}
}

func (m *reportingPluginInfo1Metrics) Close() {
	m.registerer.Unregister(m.info)
	m.limits.Close()
}

type reportingPluginLimitsMetrics struct {
	registerer prometheus.Registerer

	maxQueryBytes                                   prometheus.Gauge
	maxObservationBytes                             prometheus.Gauge
	maxReportsPlusPrecursorBytes                    prometheus.Gauge
	maxReportBytes                                  prometheus.Gauge
	maxReportCount                                  prometheus.Gauge
	maxKeyValueModifiedKeys                         prometheus.Gauge
	maxKeyValueModifiedKeysPlusValuesBytes          prometheus.Gauge
	maxBlobPayloadBytes                             prometheus.Gauge
	maxPerOracleUnexpiredBlobCumulativePayloadBytes prometheus.Gauge
	maxPerOracleUnexpiredBlobCount                  prometheus.Gauge
}

func newReportingPluginLimitsMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
	limits ocr3_1types.ReportingPluginLimits,
) *reportingPluginLimitsMetrics {
	newLimitGauge := func(name string, value int) prometheus.Gauge {
		gauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: name,
			Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types#ReportingPluginLimits for details",
		})
		gauge.Set(float64(value))
		metricshelper.RegisterOrLogError(logger, registerer, gauge, name)
		return gauge
	}

	maxQueryBytes := newLimitGauge("ocr3_1_reporting_plugin_limit_max_query_bytes", limits.MaxQueryBytes)
	maxObservationBytes := newLimitGauge("ocr3_1_reporting_plugin_limit_max_observation_bytes", limits.MaxObservationBytes)
	maxReportsPlusPrecursorBytes := newLimitGauge("ocr3_1_reporting_plugin_limit_max_reports_plus_precursor_bytes", limits.MaxReportsPlusPrecursorBytes)
	maxReportBytes := newLimitGauge("ocr3_1_reporting_plugin_limit_max_report_bytes", limits.MaxReportBytes)
	maxReportCount := newLimitGauge("ocr3_1_reporting_plugin_limit_max_report_count", limits.MaxReportCount)
	maxKeyValueModifiedKeys := newLimitGauge("ocr3_1_reporting_plugin_limit_max_key_value_modified_keys", limits.MaxKeyValueModifiedKeys)
	maxKeyValueModifiedKeysPlusValuesBytes := newLimitGauge("ocr3_1_reporting_plugin_limit_max_key_value_modified_keys_plus_values_bytes", limits.MaxKeyValueModifiedKeysPlusValuesBytes)
	maxBlobPayloadBytes := newLimitGauge("ocr3_1_reporting_plugin_limit_max_blob_payload_bytes", limits.MaxBlobPayloadBytes)
	maxPerOracleUnexpiredBlobCumulativePayloadBytes := newLimitGauge("ocr3_1_reporting_plugin_limit_max_per_oracle_unexpired_blob_cumulative_payload_bytes", limits.MaxPerOracleUnexpiredBlobCumulativePayloadBytes)
	maxPerOracleUnexpiredBlobCount := newLimitGauge("ocr3_1_reporting_plugin_limit_max_per_oracle_unexpired_blob_count", limits.MaxPerOracleUnexpiredBlobCount)

	return &reportingPluginLimitsMetrics{
		registerer,
		maxQueryBytes,
		maxObservationBytes,
		maxReportsPlusPrecursorBytes,
		maxReportBytes,
		maxReportCount,
		maxKeyValueModifiedKeys,
		maxKeyValueModifiedKeysPlusValuesBytes,
		maxBlobPayloadBytes,
		maxPerOracleUnexpiredBlobCumulativePayloadBytes,
		maxPerOracleUnexpiredBlobCount,
	}
}

func (m *reportingPluginLimitsMetrics) Close() {
	m.registerer.Unregister(m.maxQueryBytes)
	m.registerer.Unregister(m.maxObservationBytes)
	m.registerer.Unregister(m.maxReportsPlusPrecursorBytes)
	m.registerer.Unregister(m.maxReportBytes)
	m.registerer.Unregister(m.maxReportCount)
	m.registerer.Unregister(m.maxKeyValueModifiedKeys)
	m.registerer.Unregister(m.maxKeyValueModifiedKeysPlusValuesBytes)
	m.registerer.Unregister(m.maxBlobPayloadBytes)
	m.registerer.Unregister(m.maxPerOracleUnexpiredBlobCumulativePayloadBytes)
	m.registerer.Unregister(m.maxPerOracleUnexpiredBlobCount)
}
