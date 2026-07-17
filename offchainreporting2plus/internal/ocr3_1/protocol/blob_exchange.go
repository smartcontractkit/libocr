package protocol

import (
	"context"
	"fmt"
	"iter"
	"slices"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/byzquorum"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/internal/mt"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/scheduler"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol/requestergadget"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

func RunBlobExchange[RI any](
	ctx context.Context,

	chNetToBlobExchange <-chan MessageToBlobExchangeWithSender[RI],

	chBlobBroadcastRequest <-chan blobBroadcastRequest,
	chBlobFetchRequest <-chan blobFetchRequest,

	config ocr3_1config.SharedConfig,
	kv KeyValueDatabase,
	id commontypes.OracleID,
	limits ocr3_1types.ReportingPluginLimits,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	telemetrySender TelemetrySender,

	initialHighestCommittedSeqNr uint64,
) {
	broadcastGraceTimeoutScheduler := scheduler.NewScheduler[EventBlobBroadcastGraceTimeout[RI]]()
	defer broadcastGraceTimeoutScheduler.Close()

	bex := makeBlobExchangeState[RI](
		ctx, chNetToBlobExchange,
		chBlobBroadcastRequest, chBlobFetchRequest,
		config, kv,
		id, limits, localConfig, logger, metricsRegisterer, netSender, offchainKeyring,
		telemetrySender,
		initialHighestCommittedSeqNr,
		broadcastGraceTimeoutScheduler,
	)
	bex.run()
}

func makeBlobExchangeState[RI any](
	ctx context.Context,

	chNetToBlobExchange <-chan MessageToBlobExchangeWithSender[RI],

	chBlobBroadcastRequest <-chan blobBroadcastRequest,
	chBlobFetchRequest <-chan blobFetchRequest,

	config ocr3_1config.SharedConfig,
	kv KeyValueDatabase,
	id commontypes.OracleID,
	limits ocr3_1types.ReportingPluginLimits,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	telemetrySender TelemetrySender,
	initialHighestCommittedSeqNr uint64,

	broadcastGraceTimeoutScheduler *scheduler.Scheduler[EventBlobBroadcastGraceTimeout[RI]],
) *blobExchangeState[RI] {
	chHighestCommittedSeqNrTrackerToBlobExchange := make(chan uint64)
	chBlobExchangeToBlobReap := make(chan uint64)

	var submitters []*blobSubmitter
	var perOracleMetrics []*blobOracleMetrics
	for i := range config.N() {
		submitters = append(submitters, &blobSubmitter{
			loghelper.LogarithmicTaper{},
			loghelper.LogarithmicTaper{},
		})
		perOracleMetrics = append(perOracleMetrics, newBlobOracleMetrics(
			metricsRegisterer,
			logger,
			commontypes.OracleID(i),
			limits,
		))
	}

	bex := &blobExchangeState[RI]{
		ctx,
		subprocesses.Subprocesses{},

		make(chan EventToBlobExchange[RI]),
		chNetToBlobExchange,

		chBlobBroadcastRequest,
		chBlobFetchRequest,

		config,
		kv,
		id,
		limits,
		localConfig,
		logger.MakeUpdated(commontypes.LogFields{"proto": "bex"}),
		newBlobExchangeMetrics(metricsRegisterer, logger),
		netSender,
		offchainKeyring,
		telemetrySender,

		broadcastGraceTimeoutScheduler,
		nil, // must be filled right below

		nil, // must be filled right below
		chHighestCommittedSeqNrTrackerToBlobExchange,
		chBlobExchangeToBlobReap,
		initialHighestCommittedSeqNr,
		true,

		submitters,
		perOracleMetrics,
		make(map[BlobDigest]*blob),
	}

	offerRequesterGadget := requestergadget.NewRequesterGadget[blobOfferItem](
		config.N(),
		config.GetDeltaBlobOfferMinRequestToSameOracleInterval(),
		bex.trySendBlobOffer,
		bex.getPendingBlobOffers,
		bex.getBlobOfferSeeders,
	)
	bex.offerRequesterGadget = offerRequesterGadget

	chunkRequesterGadget := requestergadget.NewRequesterGadget[blobChunkId](
		config.N(),
		config.GetDeltaBlobChunkMinRequestToSameOracleInterval(),
		bex.trySendBlobChunkRequest,
		bex.getPendingBlobChunks,
		bex.getBlobChunkSeeders,
	)
	bex.chunkRequesterGadget = chunkRequesterGadget

	return bex
}

func (bex *blobExchangeState[RI]) trySendBlobChunkRequest(id blobChunkId, seeder commontypes.OracleID) (*requestergadget.RequestInfo, bool) {
	blob, ok := bex.blobs[id.blobDigest]
	if !ok {
		return nil, false
	}

	if blob.fetch == nil {
		return nil, false
	}

	timeout := bex.config.GetDeltaBlobChunkResponseTimeout()

	bex.logger.Debug("sending MessageBlobChunkRequest", commontypes.LogFields{
		"blobDigest": id.blobDigest,
		"chunkIndex": id.chunkIndex,
		"timeout":    timeout,
		"seeder":     seeder,
	})

	requestInfo := &types.RequestInfo{
		time.Now().Add(timeout),
	}
	bex.netSender.SendTo(MessageBlobChunkRequest[RI]{
		types.EmptyRequestHandleForOutboundRequest,
		requestInfo,
		id.blobDigest,
		id.chunkIndex,
	}, seeder)

	return requestInfo, true
}

func (bex *blobExchangeState[RI]) getBlobDigestsOrderedByTimeWhenAdded() []BlobDigest {
	type timedBlobDigest struct {
		blobDigest BlobDigest
		time       time.Time
	}
	timedBlobDigests := make([]timedBlobDigest, 0, len(bex.blobs))
	for blobDigest, blob := range bex.blobs {
		timedBlobDigests = append(timedBlobDigests, timedBlobDigest{blobDigest, blob.timeWhenAdded})
	}

	slices.SortFunc(timedBlobDigests, func(a, b timedBlobDigest) int {
		return a.time.Compare(b.time)
	})

	blobDigests := make([]BlobDigest, 0, len(timedBlobDigests))
	for _, timedBlobDigest := range timedBlobDigests {
		blobDigests = append(blobDigests, timedBlobDigest.blobDigest)
	}
	return blobDigests
}

func (bex *blobExchangeState[RI]) getPendingBlobChunks() []blobChunkId {
	var pending []blobChunkId
	for _, blobDigest := range bex.getBlobDigestsOrderedByTimeWhenAdded() {
		blob := bex.blobs[blobDigest]
		fetch := blob.fetch
		if fetch == nil {
			continue
		}
		if fetch.expired {
			continue
		}
		for chunkIndex := range blob.chunkDigests {
			if blob.chunkHaves[chunkIndex] {
				continue
			}
			pending = append(pending, blobChunkId{blobDigest, uint64(chunkIndex)})
		}
	}
	return pending
}

func (bex *blobExchangeState[RI]) getBlobChunkSeeders(id blobChunkId) map[commontypes.OracleID]struct{} {
	blob, ok := bex.blobs[id.blobDigest]
	if !ok {
		return nil
	}
	if blob.fetch == nil {
		return nil
	}
	return blob.fetch.seeders
}

func (bex *blobExchangeState[RI]) trySendBlobOffer(item blobOfferItem, seeder commontypes.OracleID) (*requestergadget.RequestInfo, bool) {
	blob, ok := bex.blobs[item.blobDigest]
	if !ok {
		return nil, false
	}

	if blob.broadcast == nil {
		return nil, false
	}
	if !blob.broadcast.shouldOfferTo(seeder) {
		return nil, false
	}

	timeout := bex.config.GetDeltaBlobOfferResponseTimeout()

	bex.logger.Trace("sending MessageBlobOffer", commontypes.LogFields{
		"blobDigest":       item.blobDigest,
		"chunkDigestsRoot": blob.chunkDigestsRoot,
		"payloadLength":    blob.payloadLength,
		"expirySeqNr":      blob.expirySeqNr,
		"timeout":          timeout,
		"to":               seeder,
	})

	requestInfo := &types.RequestInfo{
		time.Now().Add(timeout),
	}
	bex.netSender.SendTo(MessageBlobOffer[RI]{
		types.EmptyRequestHandleForOutboundRequest,
		requestInfo,
		blob.chunkDigestsRoot,
		blob.payloadLength,
		blob.expirySeqNr,
	}, seeder)

	return requestInfo, true
}

func (bex *blobExchangeState[RI]) getPendingBlobOffers() []blobOfferItem {
	var pending []blobOfferItem
	for _, blobDigest := range bex.getBlobDigestsOrderedByTimeWhenAdded() {
		blob := bex.blobs[blobDigest]
		if blob.broadcast == nil {
			continue
		}
		if !blob.broadcast.shouldOffer() {
			continue
		}
		for oracleID := range blob.broadcast.oracles {
			if !blob.broadcast.shouldOfferTo(commontypes.OracleID(oracleID)) {
				continue
			}
			pending = append(pending, blobOfferItem{blobDigest, commontypes.OracleID(oracleID)})
		}
	}
	return pending
}

func (bex *blobExchangeState[RI]) getBlobOfferSeeders(item blobOfferItem) map[commontypes.OracleID]struct{} {
	return map[commontypes.OracleID]struct{}{
		item.oracleID: {},
	}
}

const (
	// If we receive an offer from an oracle while we have at least this many
	// concurrent blob fetches with pending offer responses to the same oracle,
	// we will reject the offer.
	maxOwedOfferResponsesPerOracle = 10

	// trackHighestCommittedSeqNrInterval denotes the rough interval with which blob
	// exchange refreshes its view of the highest committed seq nr.
	trackHighestCommittedSeqNrInterval = 1 * time.Second
)

type blobBroadcastRequest struct {
	payload     []byte
	expirySeqNr uint64
	chResponse  chan blobBroadcastResponse
	chDone      <-chan struct{}
}

func (req *blobBroadcastRequest) respond(ctx context.Context, resp blobBroadcastResponse) {
	select {
	case req.chResponse <- resp:
	case <-req.chDone:
	case <-ctx.Done():
	}
}

type blobBroadcastResponse struct {
	cert LightCertifiedBlob
	err  error
}

type blobFetchRequest struct {
	cert       LightCertifiedBlob
	chResponse chan blobFetchResponse
	chDone     <-chan struct{}
}

func (req *blobFetchRequest) respond(ctx context.Context, resp blobFetchResponse) {
	select {
	case req.chResponse <- resp:
	case <-req.chDone:
	case <-ctx.Done():
	}
}

type blobFetchResponse struct {
	payload []byte
	err     error
}

type blobExchangeState[RI any] struct {
	ctx  context.Context
	subs subprocesses.Subprocesses

	chLocalEvent        chan EventToBlobExchange[RI]
	chNetToBlobExchange <-chan MessageToBlobExchangeWithSender[RI]

	chBlobBroadcastRequest <-chan blobBroadcastRequest
	chBlobFetchRequest     <-chan blobFetchRequest

	config          ocr3_1config.SharedConfig
	kv              KeyValueDatabase
	id              commontypes.OracleID
	limits          ocr3_1types.ReportingPluginLimits
	localConfig     types.LocalConfig
	logger          loghelper.LoggerWithContext
	metrics         *blobExchangeMetrics
	netSender       NetworkSender[RI]
	offchainKeyring types.OffchainKeyring
	telemetrySender TelemetrySender

	// blob broadcast
	broadcastGraceTimeoutScheduler *scheduler.Scheduler[EventBlobBroadcastGraceTimeout[RI]]
	offerRequesterGadget           *requestergadget.RequesterGadget[blobOfferItem]

	// blob fetch
	chunkRequesterGadget *requestergadget.RequesterGadget[blobChunkId]

	chHighestCommittedSeqNrTrackerToBlobExchange chan uint64
	chBlobExchangeToBlobReap                     chan uint64
	highestCommittedSeqNr                        uint64
	notifyBlobReapOfHighestCommittedSeqNr        bool

	submitters       []*blobSubmitter
	perOracleMetrics []*blobOracleMetrics
	blobs            map[BlobDigest]*blob
}

type blobSubmitter struct {
	rejectedOfferLogTaper loghelper.LogarithmicTaper
	statsOverflowLogTaper loghelper.LogarithmicTaper
}

type blobOfferItem struct {
	blobDigest BlobDigest
	oracleID   commontypes.OracleID
}

type blobChunkId struct {
	blobDigest BlobDigest
	chunkIndex uint64
}

func numChunks(payloadLength uint64, cfgBlobChunkSize int) uint64 {
	return (payloadLength + uint64(cfgBlobChunkSize) - 1) / uint64(cfgBlobChunkSize)
}

func chunkBytes(payloadLength uint64, chunkIndex uint64, cfgBlobChunkSize int) (uint64, bool) {
	if !(0 <= chunkIndex && chunkIndex < numChunks(payloadLength, cfgBlobChunkSize)) {
		return 0, false
	}
	cfgBlobChunkSizeU64 := uint64(cfgBlobChunkSize)
	remainingAfterPreviousChunks := payloadLength - (chunkIndex * cfgBlobChunkSizeU64)
	return min(remainingAfterPreviousChunks, cfgBlobChunkSizeU64), true
}

func chunkPayload(payload []byte, cfgBlobChunkSize int) iter.Seq2[uint64, []byte] {
	return func(yield func(uint64, []byte) bool) {
		chunkIdx := uint64(0)
		for chunk := range slices.Chunk(payload, cfgBlobChunkSize) {
			if !yield(chunkIdx, chunk) {
				return
			}
			chunkIdx++
		}
	}
}

type blobFetchMeta struct {
	chNotify chan struct{}
	waiters  int
	exchange *blobExchangeMeta
	seeders  map[commontypes.OracleID]struct{}
	expired  bool
}

func (bifm *blobFetchMeta) weServiced() {
	bifm.waiters--
}

func (bifm *blobFetchMeta) prunable() bool {
	if bifm == nil {
		return true
	}
	return bifm.waiters <= 0 && bifm.exchange.prunable()
}

type blobBroadcastPhase string

const (
	blobBroadcastPhaseOffering      blobBroadcastPhase = "offering"
	blobBroadcastPhaseAcceptedGrace blobBroadcastPhase = "acceptedGrace"
	blobBroadcastPhaseAccepted      blobBroadcastPhase = "accepted"
	blobBroadcastPhaseRejected      blobBroadcastPhase = "rejected"
	blobBroadcastPhaseExpired       blobBroadcastPhase = "expired"
)

type blobBroadcastMeta struct {
	chNotify chan struct{}
	waiters  int
	phase    blobBroadcastPhase
	// certOrNil == nil indicates we still have not assembled a cert.
	certOrNil *LightCertifiedBlob
	oracles   []blobBroadcastOracleMeta
}

type blobBroadcastOracleMeta struct {
	weReceivedOfferResponse          bool
	weReceivedOfferResponseAccepting bool
	signature                        BlobAvailabilitySignature
}

func (bibm *blobBroadcastMeta) shouldOffer() bool {
	return bibm.phase == blobBroadcastPhaseOffering || bibm.phase == blobBroadcastPhaseAcceptedGrace
}

func (bibm *blobBroadcastMeta) shouldOfferTo(oracleID commontypes.OracleID) bool {
	return bibm.shouldOffer() && !bibm.oracles[oracleID].weReceivedOfferResponse
}

func (bibm *blobBroadcastMeta) weServiced() {
	bibm.waiters--
}

func (bibm *blobBroadcastMeta) prunable() bool {
	if bibm == nil {
		return true
	}
	return bibm.waiters <= 0
}

type blobExchangeMeta struct {
	weSentOfferResponse      bool
	latestOfferRequestHandle types.RequestHandle
}

func (biem *blobExchangeMeta) weServiced() {
	biem.weSentOfferResponse = true
}

func (biem *blobExchangeMeta) prunable() bool {
	if biem == nil {
		return true
	}
	return biem.weSentOfferResponse || biem.latestOfferRequestHandle == nil
}

type blob struct {
	timeWhenAdded time.Time

	broadcast *blobBroadcastMeta
	fetch     *blobFetchMeta

	chunkDigestsRoot mt.Digest
	chunkDigests     []BlobChunkDigest
	chunkHaves       []bool

	payloadLength uint64
	expirySeqNr   uint64
	submitter     commontypes.OracleID
}

func (b *blob) weOweOfferResponse() bool {
	return b.fetch != nil && b.fetch.exchange != nil && !b.fetch.exchange.weSentOfferResponse
}

func (b *blob) haveAllChunks() bool {
	return haveAllChunks(b.chunkHaves, b.chunkDigests, b.chunkDigestsRoot)
}

func haveAllChunks(chunkHaves []bool, chunkDigests []BlobChunkDigest, chunkDigestsRoot mt.Digest) bool {
	return !slices.Contains(chunkHaves, false) && blobtypes.MakeBlobChunkDigestsRoot(chunkDigests) == chunkDigestsRoot
}

func (b *blob) prunable() bool {
	return b.broadcast.prunable() && b.fetch.prunable()
}

func (bex *blobExchangeState[RI]) run() {
	bex.logger.Info("BlobExchange: running", nil)

	if err := bex.initializeStatsMetricsFromKV(); err != nil {
		bex.logger.Warn("BlobExchange: failed to initialize stats metrics from kv, stats metrics will only be set as updated during operation", commontypes.LogFields{
			"error": err,
		})
	}

	bex.subs.Go(func() {
		RunHighestCommittedSeqNrTracker(
			bex.ctx,
			bex.logger,
			bex.kv,
			trackHighestCommittedSeqNrInterval,
			bex.chHighestCommittedSeqNrTrackerToBlobExchange,
		)
	})

	bex.subs.Go(func() {
		RunBlobReap(bex.ctx, bex.logger, bex.kv, bex.chBlobExchangeToBlobReap, bex.perOracleMetrics)
	})

	bex.subs.Go(func() {
		RunBlobPopulateStaleIndex(bex.ctx, bex.logger, bex.kv)
	})

	// Take a reference to the ctx.Done channel once, here, to avoid taking the
	// context lock below.
	chDone := bex.ctx.Done()

	// Event Loop
	for {
		var nilOrChBlobExchangeToBlobReap chan<- uint64
		if bex.notifyBlobReapOfHighestCommittedSeqNr {
			nilOrChBlobExchangeToBlobReap = bex.chBlobExchangeToBlobReap
		}

		select {
		case nilOrChBlobExchangeToBlobReap <- bex.highestCommittedSeqNr:
			bex.notifyBlobReapOfHighestCommittedSeqNr = false

		case ev := <-bex.chLocalEvent:
			ev.processBlobExchange(bex)

		case msg := <-bex.chNetToBlobExchange:
			msg.msg.processBlobExchange(bex, msg.sender)

		case req := <-bex.chBlobBroadcastRequest:
			bex.processBlobBroadcastRequest(req)
		case req := <-bex.chBlobFetchRequest:
			bex.processBlobFetchRequest(req)

		case ev := <-bex.broadcastGraceTimeoutScheduler.Scheduled():
			ev.processBlobExchange(bex)
		case <-bex.offerRequesterGadget.Ticker():
			bex.offerRequesterGadget.Tick()

		case <-bex.chunkRequesterGadget.Ticker():
			bex.chunkRequesterGadget.Tick()
		case highestCommittedSeqNr := <-bex.chHighestCommittedSeqNrTrackerToBlobExchange:
			bex.eventHighestCommittedSeqNr(highestCommittedSeqNr)

		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			bex.logger.Info("BlobExchange: winding down", nil)
			bex.subs.Wait()
			bex.metrics.Close()
			for _, m := range bex.perOracleMetrics {
				m.Close()
			}
			bex.logger.Info("BlobExchange: exiting", nil)
			return
		default:
		}
	}
}

func (bex *blobExchangeState[RI]) eventHighestCommittedSeqNr(highestCommittedSeqNr uint64) {
	if highestCommittedSeqNr <= bex.highestCommittedSeqNr {
		return
	}
	bex.highestCommittedSeqNr = highestCommittedSeqNr

	bex.logger.Debug("blob exchange learned of even higher highest committed seq nr", commontypes.LogFields{
		"highestCommittedSeqNr": bex.highestCommittedSeqNr,
	})

	for blobDigest, blob := range bex.blobs {
		if !bex.hasBlobExpired(blob.expirySeqNr) {
			continue
		}

		broadcastPending := blob.broadcast != nil && blob.broadcast.shouldOffer()
		fetchPending := blob.fetch != nil && !blob.fetch.expired && !blob.haveAllChunks()

		if !(broadcastPending || fetchPending) {
			continue
		}

		bex.logger.Debug("stopping expired blob broadcast or fetch", commontypes.LogFields{
			"blobDigest":            blobDigest,
			"expirySeqNr":           blob.expirySeqNr,
			"highestCommittedSeqNr": highestCommittedSeqNr,
			"submitter":             blob.submitter,
			"fetchPending":          fetchPending,
			"broadcastPending":      broadcastPending,
		})

		if broadcastPending {
			broadcast := blob.broadcast
			broadcast.phase = blobBroadcastPhaseExpired
			close(broadcast.chNotify)
			bex.incMyBlobsUndeterminedTotalForNonResponders(broadcast)
		}

		if fetchPending {
			fetch := blob.fetch
			if fetch.exchange != nil {
				bex.sendBlobOfferResponseRejecting(blobDigest, blob.submitter, fetch.exchange.latestOfferRequestHandle)
				bex.perOracleMetrics[blob.submitter].theirBlobsRejectedDueToExpirationTotal.Inc()
				fetch.exchange.weServiced()
			}

			fetch.expired = true
			close(fetch.chNotify)
		}

		if blob.prunable() {
			bex.metrics.blobsInProgress.Dec()
			delete(bex.blobs, blobDigest)
		}
	}

	bex.notifyBlobReapOfHighestCommittedSeqNr = true
}

func (bex *blobExchangeState[RI]) hasBlobExpired(expirySeqNr uint64) bool {
	return expirySeqNr <= bex.highestCommittedSeqNr
}

func (bex *blobExchangeState[RI]) allowBlobOfferBasedOnOwedOfferResponsesBudget(sender commontypes.OracleID) error {
	countOwedOfferResponses := 0
	for _, blob := range bex.blobs {
		if blob.submitter != sender {
			continue
		}
		if !blob.weOweOfferResponse() {
			continue
		}
		countOwedOfferResponses++
	}

	if countOwedOfferResponses >= maxOwedOfferResponsesPerOracle {
		return fmt.Errorf("too many pending offer responses for submitter, have %d, max %d", countOwedOfferResponses, maxOwedOfferResponsesPerOracle)
	}

	return nil
}

func (bex *blobExchangeState[RI]) allowBlobOfferBasedOnQuotaStats(offer MessageBlobOffer[RI], submitter commontypes.OracleID) error {
	tx, err := bex.kv.NewReadTransactionUnchecked()
	if err != nil {
		return fmt.Errorf("failed to create read transaction: %w", err)
	}
	defer tx.Discard()

	appendedQuotaStats, err := tx.ReadBlobQuotaStats(BlobQuotaStatsTypeAppended, submitter)
	if err != nil {
		return fmt.Errorf("failed to read blob quota stats: %w", err)
	}

	reapedQuotaStats, err := tx.ReadBlobQuotaStats(BlobQuotaStatsTypeReaped, submitter)
	if err != nil {
		return fmt.Errorf("failed to read blob quota stats: %w", err)
	}

	statsOverflowLogTaper := &bex.submitters[submitter].statsOverflowLogTaper

	totalQuotaStats, ok := appendedQuotaStats.Sub(reapedQuotaStats)
	if !ok {
		statsOverflowLogTaper.Trigger(func(consecutiveOverflows uint64) {
			bex.logger.Debug("overflow when subtracting reaped quota stats from appended quota stats, "+
				"reaper is likely ahead of blob exchange, treating as zero", commontypes.LogFields{
				"consecutiveOverflows": consecutiveOverflows,
				"appendedQuotaStats":   appendedQuotaStats,
				"reapedQuotaStats":     reapedQuotaStats,
				"submitter":            submitter,
			})
		})
		totalQuotaStats = BlobQuotaStats{0, 0}
	} else {
		statsOverflowLogTaper.Reset(func(previouslyConsecutiveOverflows uint64) {
			bex.logger.Debug("stopped seeing overflow when subtracting reaped quota stats from appended quota stats", commontypes.LogFields{
				"previouslyConsecutiveOverflows": previouslyConsecutiveOverflows,
				"appendedQuotaStats":             appendedQuotaStats,
				"reapedQuotaStats":               reapedQuotaStats,
				"submitter":                      submitter,
			})
		})
	}

	totalQuotaStatsIncludingOffer, ok := totalQuotaStats.Add(BlobQuotaStats{
		1,
		offer.PayloadLength,
	})
	if !ok {
		return fmt.Errorf("overflow when considering offer in quota stats")
	}

	maxQuotaStats := BlobQuotaStats{
		uint64(bex.limits.MaxPerOracleUnexpiredBlobCount),
		uint64(bex.limits.MaxPerOracleUnexpiredBlobCumulativePayloadBytes),
	}

	if totalQuotaStatsIncludingOffer.Exceeds(maxQuotaStats) {
		return fmt.Errorf("accepting the offer would exceed our allowed per-oracle unexpired quota, "+
			"have %+v, would have %+v, max %+v", totalQuotaStats, totalQuotaStatsIncludingOffer, maxQuotaStats)
	}

	return nil
}

func (bex *blobExchangeState[RI]) allowBlobOffer(offer MessageBlobOffer[RI], submitter commontypes.OracleID) error {
	if err := bex.allowBlobOfferBasedOnOwedOfferResponsesBudget(submitter); err != nil {
		bex.perOracleMetrics[submitter].theirBlobsRejectedDueToManyInflightTotal.Inc()
		return err
	}

	if err := bex.allowBlobOfferBasedOnQuotaStats(offer, submitter); err != nil {
		bex.perOracleMetrics[submitter].theirBlobsRejectedDueToQuotaTotal.Inc()
		return err
	}

	return nil
}

func (bex *blobExchangeState[RI]) messageBlobOffer(msg MessageBlobOffer[RI], sender commontypes.OracleID) {
	submitter := sender

	bex.perOracleMetrics[submitter].theirBlobOffersTotal.Inc()

	blobDigest := blobtypes.MakeBlobDigest(
		bex.config.ConfigDigest,
		msg.ChunkDigestsRoot,
		msg.PayloadLength,
		msg.ExpirySeqNr,
		submitter,
	)

	rejectedOfferLogTaper := &bex.submitters[sender].rejectedOfferLogTaper

	// Reject if payload length exceeds maximum allowed length
	if msg.PayloadLength > uint64(bex.limits.MaxBlobPayloadBytes) {
		rejectedOfferLogTaper.Trigger(func(consecutiveRejectedOffers uint64) {
			bex.logger.Warn("received MessageBlobOffer with payload length that exceeds maximum allowed length, rejecting", commontypes.LogFields{
				"blobDigest":                blobDigest,
				"submitter":                 submitter,
				"payloadLength":             msg.PayloadLength,
				"maxPayloadLength":          bex.limits.MaxBlobPayloadBytes,
				"consecutiveRejectedOffers": consecutiveRejectedOffers,
			})
		})
		bex.sendBlobOfferResponseRejecting(blobDigest, submitter, msg.RequestHandle)
		bex.perOracleMetrics[submitter].theirBlobsRejectedDueToOversizePayloadBytesTotal.Inc()
		return
	}

	chunkDigests, chunkHaves, err := bex.loadChunkDigestsAndHaves(blobDigest, msg.PayloadLength)
	if err != nil {
		bex.logger.Warn("dropping MessageBlobOffer, failed to check if we already know of it", commontypes.LogFields{
			"blobDigest": blobDigest,
			"sender":     sender,
			"error":      err,
		})
		return
	}

	// check if we maybe already have this blob in full
	if haveAllChunks(chunkHaves, chunkDigests, msg.ChunkDigestsRoot) {
		bex.logger.Debug("received MessageBlobOffer for which we already have the payload", commontypes.LogFields{
			"blobDigest": blobDigest,
			"sender":     sender,
		})
		bex.sendBlobOfferResponseAccepting(blobDigest, submitter, msg.RequestHandle)
		return
	}

	if blob, ok := bex.blobs[blobDigest]; ok {
		bex.logger.Debug("duplicate MessageBlobOffer, updating offer request handle", commontypes.LogFields{
			"blobDigest": blobDigest,
			"sender":     sender,
		})
		if blob.fetch != nil && blob.fetch.exchange != nil {
			blob.fetch.exchange.latestOfferRequestHandle = msg.RequestHandle
		}
		return
	}

	// Reject if blob has already expired
	if bex.hasBlobExpired(msg.ExpirySeqNr) {
		rejectedOfferLogTaper.Trigger(func(consecutiveRejectedOffers uint64) {
			bex.logger.Warn("received MessageBlobOffer for already expired blob, rejecting", commontypes.LogFields{
				"blobDigest":                blobDigest,
				"submitter":                 submitter,
				"expirySeqNr":               msg.ExpirySeqNr,
				"highestCommittedSeqNr":     bex.highestCommittedSeqNr,
				"consecutiveRejectedOffers": consecutiveRejectedOffers,
			})
		})
		bex.sendBlobOfferResponseRejecting(blobDigest, submitter, msg.RequestHandle)
		bex.perOracleMetrics[submitter].theirBlobsRejectedDueToExpirationTotal.Inc()
		return
	}

	if err := bex.allowBlobOffer(msg, submitter); err != nil {

		rejectedOfferLogTaper.Trigger(func(consecutiveRejectedOffers uint64) {
			bex.logger.Info("received MessageBlobOffer that goes over rate limits, rejecting", commontypes.LogFields{
				"blobDigest":                blobDigest,
				"submitter":                 submitter,
				"reason":                    err,
				"consecutiveRejectedOffers": consecutiveRejectedOffers,
			})
		})
		bex.sendBlobOfferResponseRejecting(blobDigest, sender, msg.RequestHandle)
		// bex.perOracleMetrics[submitter].theirBlobsRejectedDueToXTotal
		// incremented inside allowBlobOffer, which knows the exact reason. No
		// need to do anything here.
		return
	}

	rejectedOfferLogTaper.Reset(func(previouslyConsecutiveRejectedOffers uint64) {
		bex.logger.Info("stopped receiving offers that we keep rejecting from submitter", commontypes.LogFields{
			"submitter":                           submitter,
			"previouslyConsecutiveRejectedOffers": previouslyConsecutiveRejectedOffers,
		})
	})

	bex.logger.Info("received MessageBlobOffer for new to us blob that we can accept", commontypes.LogFields{
		"blobDigest":       blobDigest,
		"sender":           sender,
		"chunkDigestsRoot": msg.ChunkDigestsRoot,
		"payloadLength":    msg.PayloadLength,
		"expirySeqNr":      msg.ExpirySeqNr,
	})

	seeders := map[commontypes.OracleID]struct{}{
		submitter: {},
	}

	if err := bex.createOrUpdateBlobMetaAndQuotaStats(blobDigest, msg.PayloadLength, msg.ExpirySeqNr, submitter); err != nil {
		bex.logger.Error("failed to create or update blob meta and quota stats for MessageBlobOffer", commontypes.LogFields{
			"blobDigest": blobDigest,
			"submitter":  submitter,
			"error":      err,
		})
		return
	}

	bex.metrics.blobsInProgress.Inc()
	bex.blobs[blobDigest] = &blob{
		time.Now(),
		nil,
		&blobFetchMeta{
			make(chan struct{}),
			0,
			&blobExchangeMeta{
				false,
				msg.RequestHandle,
			},
			seeders,
			false,
		},

		msg.ChunkDigestsRoot,
		chunkDigests,
		chunkHaves,
		msg.PayloadLength,
		msg.ExpirySeqNr,
		submitter,
	}
	bex.chunkRequesterGadget.PleaseRecheckPendingItems()
}

func (bex *blobExchangeState[RI]) messageBlobOfferResponse(msg MessageBlobOfferResponse[RI], sender commontypes.OracleID) {
	item := blobOfferItem{msg.BlobDigest, sender}
	if !bex.offerRequesterGadget.CheckAndMarkResponse(item, sender) {
		bex.logger.Debug("dropping MessageBlobOfferResponse, not allowed", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	blob, ok := bex.blobs[msg.BlobDigest]
	if !ok {
		bex.logger.Debug("dropping MessageBlobOfferResponse for unknown blob", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	if blob.broadcast == nil {
		bex.logger.Debug("dropping MessageBlobOfferResponse, not broadcasting", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}
	broadcast := blob.broadcast

	if blob.submitter != bex.id {
		bex.logger.Debug("dropping MessageBlobOfferResponse, not the submitter", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"submitter":  blob.submitter,
			"localID":    bex.id,
		})
		return
	}

	if !(broadcast.phase == blobBroadcastPhaseOffering || broadcast.phase == blobBroadcastPhaseAcceptedGrace) {
		bex.logger.Debug("dropping MessageBlobOfferResponse, not in offering or accepted grace phase", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"phase":      broadcast.phase,
		})
		return
	}

	// check if we already have an offer response from this oracle
	if broadcast.oracles[sender].weReceivedOfferResponse {
		bex.logger.Debug("dropping MessageBlobOfferResponse, already have message from oracle", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	// did they accept our offer?
	if !msg.RejectOffer {
		// check signature
		if err := msg.Signature.Verify(msg.BlobDigest, bex.config.OracleIdentities[sender].OffchainPublicKey); err != nil {
			bex.logger.Debug("dropping MessageBlobOfferResponse, invalid signature", commontypes.LogFields{
				"blobDigest": msg.BlobDigest,
				"sender":     sender,
			})
			return
		}

		// save signature for oracle
		broadcast.oracles[sender] = blobBroadcastOracleMeta{
			true,
			true,
			msg.Signature,
		}
		bex.perOracleMetrics[sender].myBlobsAcceptedTotal.Inc()
	} else {
		// save rejection for oracle
		broadcast.oracles[sender] = blobBroadcastOracleMeta{
			true,
			false,
			nil,
		}
		bex.perOracleMetrics[sender].myBlobsRejectedTotal.Inc()
	}

	rejectThreshold := bex.minRejectors()
	acceptThreshold := bex.minCertSigners()

	acceptingOracles, rejectingOracles := 0, 0
	for _, oracle := range broadcast.oracles {
		if oracle.weReceivedOfferResponse {
			if oracle.weReceivedOfferResponseAccepting {
				acceptingOracles++
			} else {
				rejectingOracles++
			}
		}
	}

	bex.logger.Debug("received MessageBlobOfferResponse", commontypes.LogFields{
		"blobDigest":       msg.BlobDigest,
		"sender":           sender,
		"reject":           msg.RejectOffer,
		"acceptingOracles": acceptingOracles,
		"rejectingOracles": rejectingOracles,
		"acceptThreshold":  acceptThreshold,
		"rejectThreshold":  rejectThreshold,
	})

	if broadcast.phase == blobBroadcastPhaseAcceptedGrace {
		return
	}

	if acceptingOracles >= acceptThreshold {
		bex.logger.Debug("minimum number of accepting oracles reached, entering grace period", commontypes.LogFields{
			"acceptingOracles": acceptingOracles,
			"acceptThreshold":  acceptThreshold,
			"blobDigest":       msg.BlobDigest,
			"gracePeriod":      bex.config.GetDeltaBlobBroadcastGrace(),
		})
		broadcast.phase = blobBroadcastPhaseAcceptedGrace
		bex.broadcastGraceTimeoutScheduler.ScheduleDelay(EventBlobBroadcastGraceTimeout[RI]{
			msg.BlobDigest,
		}, bex.config.GetDeltaBlobBroadcastGrace())
		return
	}

	if rejectingOracles >= rejectThreshold {
		bex.logger.Warn("oracle quorum rejected our broadcast", commontypes.LogFields{
			"rejectingOracles": rejectingOracles,
			"rejectThreshold":  rejectThreshold,
			"blobDigest":       msg.BlobDigest,
		})
		broadcast.phase = blobBroadcastPhaseRejected
		close(broadcast.chNotify)
		bex.incMyBlobsUndeterminedTotalForNonResponders(broadcast)

		return
	}
}

// Must be called once when entering terminal phases of broadcast: accepted, rejected, expired
func (bex *blobExchangeState[RI]) incMyBlobsUndeterminedTotalForNonResponders(broadcast *blobBroadcastMeta) {
	for oracleID, oracle := range broadcast.oracles {
		if !oracle.weReceivedOfferResponse {
			bex.perOracleMetrics[oracleID].myBlobsUndeterminedTotal.Inc()
		}
	}
}

func (bex *blobExchangeState[RI]) eventBlobBroadcastGraceTimeout(ev EventBlobBroadcastGraceTimeout[RI]) {
	blob, ok := bex.blobs[ev.BlobDigest]
	if !ok {
		bex.logger.Debug("dropping EventBlobBroadcastGraceTimeout for unknown blob", commontypes.LogFields{
			"blobDigest": ev.BlobDigest,
		})
		return
	}
	broadcast := blob.broadcast
	if broadcast == nil {
		bex.logger.Debug("dropping EventBlobBroadcastGraceTimeout for blob with no broadcast", commontypes.LogFields{
			"blobDigest": ev.BlobDigest,
		})
		return
	}
	if broadcast.phase != blobBroadcastPhaseAcceptedGrace {
		bex.logger.Debug("dropping EventBlobBroadcastGraceTimeout for blob not in accepted grace phase", commontypes.LogFields{
			"blobDigest": ev.BlobDigest,
			"phase":      broadcast.phase,
		})
		return
	}

	maxSigners := bex.maxCertSigners()

	var abass []AttributedBlobAvailabilitySignature
	for oracleID, oracle := range broadcast.oracles {
		if oracle.weReceivedOfferResponse && oracle.weReceivedOfferResponseAccepting && len(abass) < maxSigners {
			abass = append(abass, AttributedBlobAvailabilitySignature{
				oracle.signature,
				commontypes.OracleID(oracleID),
			})
		}
	}

	lcb := LightCertifiedBlob{
		blob.chunkDigestsRoot,
		blob.payloadLength,
		blob.expirySeqNr,
		blob.submitter,
		abass,
	}

	if err := bex.verifyCert(&lcb); err != nil {
		bex.logger.Critical("assumption violation: failed to verify own LightCertifiedBlob", commontypes.LogFields{
			"blobDigest": ev.BlobDigest,
			"error":      err,
		})
		return
	}

	bex.logger.Debug("assembled blob availability certificate", commontypes.LogFields{
		"acceptingOracles": len(abass),
		"blobDigest":       ev.BlobDigest,
	})

	broadcast.certOrNil = &lcb
	broadcast.phase = blobBroadcastPhaseAccepted
	close(broadcast.chNotify)
	bex.incMyBlobsUndeterminedTotalForNonResponders(broadcast)
}

func (bex *blobExchangeState[RI]) sendBlobOfferResponseAccepting(blobDigest BlobDigest, submitter commontypes.OracleID, requestHandle types.RequestHandle) {

	bas, err := blobtypes.MakeBlobAvailabilitySignature(blobDigest, bex.offchainKeyring.OffchainSign)
	if err != nil {
		bex.logger.Error("failed to make blob availability signature", commontypes.LogFields{
			"blobDigest": blobDigest,
			"error":      err,
		})
		return
	}

	bex.logger.Debug("sending accepting MessageBlobOfferResponse", commontypes.LogFields{
		"blobDigest": blobDigest,
		"submitter":  submitter,
	})
	bex.netSender.SendTo(
		MessageBlobOfferResponse[RI]{
			requestHandle,
			blobDigest,
			false,
			bas,
		},
		submitter,
	)

	bex.perOracleMetrics[submitter].theirBlobAcceptsSentTotal.Inc()
}
func (bex *blobExchangeState[RI]) sendBlobOfferResponseRejecting(blobDigest BlobDigest, submitter commontypes.OracleID, requestHandle types.RequestHandle) {
	bex.netSender.SendTo(
		MessageBlobOfferResponse[RI]{
			requestHandle,
			blobDigest,
			true,
			nil,
		},
		submitter,
	)

	bex.perOracleMetrics[submitter].theirBlobRejectsSentTotal.Inc()
}

func (bex *blobExchangeState[RI]) readBlobPayload(blobDigest BlobDigest) ([]byte, error) {
	tx, err := bex.kv.NewReadTransactionUnchecked()
	if err != nil {
		return nil, fmt.Errorf("failed to create read transaction: %w", err)
	}
	defer tx.Discard()

	payload, err := tx.ReadBlobPayload(blobDigest)
	if err != nil {
		return nil, fmt.Errorf("failed to read blob payload: %w", err)
	}
	return payload, nil
}

func (bex *blobExchangeState[RI]) messageBlobChunkRequest(msg MessageBlobChunkRequest[RI], sender commontypes.OracleID) {
	chunkIndex := msg.ChunkIndex

	bex.logger.Trace("received MessageBlobChunkRequest", commontypes.LogFields{
		"blobDigest": msg.BlobDigest,
		"sender":     sender,
		"chunkIndex": chunkIndex,
	})

	tx, err := bex.kv.NewReadTransactionUnchecked()
	if err != nil {
		bex.logger.Error("failed to create read transaction for MessageBlobChunkRequest", commontypes.LogFields{
			"error": err,
		})
		return
	}
	defer tx.Discard()

	meta, err := tx.ReadBlobMeta(msg.BlobDigest)
	if err != nil {
		bex.logger.Error("failed to read blob meta for MessageBlobChunkRequest", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"error":      err,
		})
		return
	}

	if meta == nil {
		bex.logger.Debug("dropping MessageBlobChunkRequest, blob meta does not exist", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
		})
		return
	}

	if slices.Contains(meta.ChunkHaves, false) {
		bex.logger.Debug("dropping MessageBlobChunkRequest, we do not have all chunks", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
		})
		return
	}

	chunk, err := tx.ReadBlobChunk(msg.BlobDigest, chunkIndex)
	if err != nil {
		bex.logger.Error("failed to read blob chunk for MessageBlobChunkRequest", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
			"error":      err,
		})
		return
	}

	goAway := chunk == nil

	proof, err := blobtypes.ProveBlobChunkDigest(meta.ChunkDigests, chunkIndex)
	if err != nil {
		bex.logger.Error("failed to prove blob chunk digest for MessageBlobChunkRequest", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
			"error":      err,
		})
		return
	}

	bex.logger.Debug("sending MessageBlobChunkResponse", commontypes.LogFields{
		"blobDigest": msg.BlobDigest,
		"chunkIndex": chunkIndex,
		"goAway":     goAway,
		"to":         sender,
	})

	bex.netSender.SendTo(
		MessageBlobChunkResponse[RI]{
			msg.RequestHandle,
			msg.BlobDigest,
			chunkIndex,
			goAway,
			chunk,
			proof,
		},
		sender,
	)
}

func (bex *blobExchangeState[RI]) messageBlobChunkResponse(msg MessageBlobChunkResponse[RI], sender commontypes.OracleID) {
	bcid := blobChunkId{msg.BlobDigest, msg.ChunkIndex}
	if !bex.chunkRequesterGadget.CheckAndMarkResponse(bcid, sender) {
		bex.logger.Debug("dropping MessageBlobChunkResponse, not allowed", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	if msg.GoAway {
		bex.logger.Debug("dropping MessageBlobChunkResponse, go away", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		bex.chunkRequesterGadget.MarkGoAwayResponse(bcid, sender)
		return
	}

	blob, ok := bex.blobs[msg.BlobDigest]
	if !ok {
		bex.logger.Debug("dropping MessageBlobChunkResponse for unknown blob", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	fetch := blob.fetch
	if fetch == nil {
		bex.logger.Debug("dropping MessageBlobChunkResponse for blob with no fetch", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}
	if fetch.expired {
		bex.logger.Debug("dropping MessageBlobChunkResponse for expired blob", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	chunkIndex := msg.ChunkIndex

	if !(0 <= chunkIndex && chunkIndex < uint64(len(blob.chunkDigests))) {
		bex.logger.Warn("dropping MessageBlobChunkResponse, chunk index out of range", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
			"chunkCount": len(blob.chunkDigests),
		})
		bex.chunkRequesterGadget.MarkBadResponse(bcid, sender)
		return
	}

	expectedChunkBytes, ok := chunkBytes(blob.payloadLength, chunkIndex, bex.config.GetBlobChunkBytes())

	if !ok || uint64(len(msg.Chunk)) != expectedChunkBytes {
		bex.logger.Warn("dropping MessageBlobChunkResponse, incorrectly sized chunk", commontypes.LogFields{
			"blobDigest":         msg.BlobDigest,
			"sender":             sender,
			"chunkIndex":         chunkIndex,
			"chunkCount":         len(blob.chunkDigests),
			"payloadLength":      blob.payloadLength,
			"expectedChunkBytes": expectedChunkBytes,
			"actualChunkBytes":   len(msg.Chunk),
		})
		bex.chunkRequesterGadget.MarkBadResponse(bcid, sender)
		return
	}

	if blob.chunkHaves[chunkIndex] {
		bex.logger.Debug("dropping MessageBlobChunkResponse, already have chunk", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
		})
		return
	}

	actualChunkDigest := blobtypes.MakeBlobChunkDigest(msg.Chunk)
	if err := blobtypes.VerifyBlobChunkDigest(blob.chunkDigestsRoot, chunkIndex, actualChunkDigest, msg.Proof); err != nil {
		bex.logger.Debug("dropping MessageBlobChunkResponse, chunk digest verification failed", commontypes.LogFields{
			"blobDigest":       msg.BlobDigest,
			"sender":           sender,
			"chunkIndex":       chunkIndex,
			"actualDigest":     actualChunkDigest,
			"chunkDigestsRoot": blob.chunkDigestsRoot,
			"error":            err,
		})
		bex.chunkRequesterGadget.MarkBadResponse(bcid, sender)
		return
	}

	bex.chunkRequesterGadget.MarkGoodResponse(bcid, sender)

	bex.logger.Debug("received MessageBlobChunkResponse", commontypes.LogFields{
		"blobDigest":    msg.BlobDigest,
		"sender":        sender,
		"chunkIndex":    chunkIndex,
		"payloadLength": blob.payloadLength,
	})

	tx, err := bex.kv.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		bex.logger.Error("failed to create read-write transaction for MessageBlobChunkResponse", commontypes.LogFields{
			"error": err,
		})
		return
	}
	defer tx.Discard()

	err = tx.WriteBlobChunk(msg.BlobDigest, chunkIndex, msg.Chunk)
	if err != nil {
		bex.logger.Error("failed to write blob chunk for MessageBlobChunkResponse", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
			"error":      err,
		})
		return
	}

	chunkHaves := slices.Clone(blob.chunkHaves)
	chunkHaves[chunkIndex] = true
	chunkDigests := slices.Clone(blob.chunkDigests)
	chunkDigests[chunkIndex] = actualChunkDigest
	blobMeta := BlobMeta{
		blob.payloadLength,
		chunkHaves,
		chunkDigests,
		blob.expirySeqNr,
		blob.submitter,
	}
	err = tx.WriteBlobMeta(msg.BlobDigest, blobMeta)
	if err != nil {
		bex.logger.Error("failed to write blob meta for MessageBlobChunkResponse", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
			"error":      err,
		})
		return
	}

	err = tx.WriteStaleBlobIndex(staleBlob(blob.expirySeqNr, msg.BlobDigest))
	if err != nil {
		bex.logger.Error("failed to write stale blob index for MessageBlobChunkResponse", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"error":      err,
		})
		return
	}

	err = tx.Commit()
	if err != nil {
		bex.logger.Error("failed to commit transaction for MessageBlobChunkResponse", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
			"error":      err,
		})
		return
	}

	blob.chunkHaves[chunkIndex] = true
	blob.chunkDigests[chunkIndex] = actualChunkDigest

	if !blob.haveAllChunks() {
		return
	}

	bex.logger.Debug("blob fully received", commontypes.LogFields{
		"blobDigest":    msg.BlobDigest,
		"sender":        sender,
		"payloadLength": blob.payloadLength,
	})
	close(fetch.chNotify)
	if fetch.exchange != nil {
		bex.sendBlobOfferResponseAccepting(msg.BlobDigest, blob.submitter, fetch.exchange.latestOfferRequestHandle)
		fetch.exchange.weServiced()
	}
	if blob.prunable() {
		bex.metrics.blobsInProgress.Dec()
		delete(bex.blobs, msg.BlobDigest)
	}
}

func (bex *blobExchangeState[RI]) processBlobBroadcastRequest(req blobBroadcastRequest) {
	if len(req.payload) > bex.limits.MaxBlobPayloadBytes {
		req.respond(bex.ctx, blobBroadcastResponse{
			LightCertifiedBlob{},
			fmt.Errorf("blob payload length %d exceeds maximum allowed length %d",
				len(req.payload), bex.limits.MaxBlobPayloadBytes),
		})
		return
	}
	expirySeqNr := req.expirySeqNr
	if bex.hasBlobExpired(expirySeqNr) {
		req.respond(bex.ctx, blobBroadcastResponse{
			LightCertifiedBlob{},
			fmt.Errorf("blob expiry seq nr %d is not above highest committed seq nr %d", expirySeqNr, bex.highestCommittedSeqNr),
		})
		return
	}

	payload := req.payload
	payloadLength := uint64(len(payload))
	cfgBlobChunkSize := bex.config.GetBlobChunkBytes()

	chunkDigests := make([]BlobChunkDigest, 0, numChunks(payloadLength, cfgBlobChunkSize))
	chunkHaves := make([]bool, 0, numChunks(payloadLength, cfgBlobChunkSize))

	for _, payloadChunk := range chunkPayload(payload, cfgBlobChunkSize) {
		// prepare for offer
		chunkDigest := blobtypes.MakeBlobChunkDigest(payloadChunk)
		chunkDigests = append(chunkDigests, chunkDigest)

		// for local accounting
		chunkHaves = append(chunkHaves, true)
	}

	submitter := bex.id

	chunkDigestsRoot := blobtypes.MakeBlobChunkDigestsRoot(chunkDigests)

	blobDigest := blobtypes.MakeBlobDigest(
		bex.config.ConfigDigest,
		chunkDigestsRoot,
		payloadLength,
		expirySeqNr,
		submitter,
	)

	bex.logger.Debug("processing BlobBroadcastRequest", commontypes.LogFields{"blobDigest": blobDigest})

	var chNotifyCertAvailable chan struct{}
	startedNewBroadcast := false
	if existingBlob, ok := bex.blobs[blobDigest]; ok {
		if existingBlob.broadcast == nil {
			existingBlob.broadcast = &blobBroadcastMeta{
				make(chan struct{}),
				1,
				blobBroadcastPhaseOffering,
				nil,
				make([]blobBroadcastOracleMeta, bex.config.N()),
			}
			startedNewBroadcast = true
		} else {
			existingBlob.broadcast.waiters++
		}
		chNotifyCertAvailable = existingBlob.broadcast.chNotify
	} else {
		// if we haven't written the chunks to kv, we can't serve requests

		if err := bex.writeBlobBeforeBroadcast(blobDigest, payloadLength, payload, expirySeqNr); err != nil {
			req.respond(bex.ctx, blobBroadcastResponse{
				LightCertifiedBlob{},
				fmt.Errorf("failed to write blob: %w", err),
			})
			return
		}

		// write in-memory state
		chNotifyCertAvailable = make(chan struct{})

		bex.metrics.blobsInProgress.Inc()
		bex.blobs[blobDigest] = &blob{
			time.Now(),
			&blobBroadcastMeta{
				chNotifyCertAvailable,
				1,
				blobBroadcastPhaseOffering,
				nil,
				make([]blobBroadcastOracleMeta, bex.config.N()),
			},
			nil,

			chunkDigestsRoot,
			chunkDigests,
			chunkHaves,
			payloadLength,
			expirySeqNr,
			submitter,
		}
		startedNewBroadcast = true
	}

	if startedNewBroadcast {
		bex.metrics.myBroadcastPayloadBytes.Observe(float64(payloadLength))
	}

	bex.offerRequesterGadget.PleaseRecheckPendingItems()

	chDone := bex.ctx.Done()

	bex.subs.Go(func() {
		select {
		case <-req.chDone:
		case <-chDone:
			return
		}

		select {
		case bex.chLocalEvent <- EventBlobBroadcastRequestDone[RI]{blobDigest}:
		case <-chDone:
		}
	})

	bex.subs.Go(func() {
		select {
		case <-chNotifyCertAvailable:
		case <-chDone:
			return
		case <-req.chDone:
			return
		}

		select {
		case bex.chLocalEvent <- EventBlobBroadcastRequestRespond[RI]{blobDigest, req}:
		case <-req.chDone:
		case <-chDone:
		}
	})
}

func (bex *blobExchangeState[RI]) getCert(blobDigest BlobDigest) (LightCertifiedBlob, error) {
	blob, ok := bex.blobs[blobDigest]
	if !ok {
		return LightCertifiedBlob{}, fmt.Errorf("no such blob, unexpected")
	}
	if blob.broadcast == nil {
		return LightCertifiedBlob{}, fmt.Errorf("no broadcast metadata available, unexpected")
	}
	switch blob.broadcast.phase {
	case blobBroadcastPhaseOffering:
		return LightCertifiedBlob{}, fmt.Errorf("blob still being offered, unexpected")
	case blobBroadcastPhaseAcceptedGrace:
		return LightCertifiedBlob{}, fmt.Errorf("blob still in grace period, unexpected")
	case blobBroadcastPhaseRejected:
		return LightCertifiedBlob{}, fmt.Errorf("blob broadcast rejected by quorum")
	case blobBroadcastPhaseExpired:
		return LightCertifiedBlob{}, fmt.Errorf("blob broadcast expired")
	case blobBroadcastPhaseAccepted:
		if blob.broadcast.certOrNil == nil {
			return LightCertifiedBlob{}, fmt.Errorf("blob was accepted but cert is nil, unexpected")
		}
		return *blob.broadcast.certOrNil, nil
	}
	panic("unreachable")
}

func (bex *blobExchangeState[RI]) eventBlobBroadcastRequestRespond(ev EventBlobBroadcastRequestRespond[RI]) {
	cert, err := bex.getCert(ev.BlobDigest)
	ev.Request.respond(bex.ctx, blobBroadcastResponse{cert, err})
}

func (bex *blobExchangeState[RI]) eventBlobBroadcastRequestDone(ev EventBlobBroadcastRequestDone[RI]) {
	blob, ok := bex.blobs[ev.BlobDigest]
	if !ok {
		return
	}
	broadcast := blob.broadcast
	if broadcast != nil {
		broadcast.weServiced()
	}
	if blob.prunable() {
		bex.metrics.blobsInProgress.Dec()
		delete(bex.blobs, ev.BlobDigest)
	}
}

func (bex *blobExchangeState[RI]) writeBlobBeforeBroadcast(blobDigest BlobDigest, payloadLength uint64, payload []byte, expirySeqNr uint64) error {
	cfgBlobChunkSize := bex.config.GetBlobChunkBytes()
	tx, err := bex.kv.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return fmt.Errorf("failed to create read/write transaction: %w", err)
	}
	defer tx.Discard()

	numChunks := numChunks(payloadLength, cfgBlobChunkSize)

	chunkHaves := make([]bool, numChunks)
	chunkDigests := make([]BlobChunkDigest, numChunks)
	for chunkIdx, payloadChunk := range chunkPayload(payload, cfgBlobChunkSize) {
		if err := tx.WriteBlobChunk(blobDigest, chunkIdx, payloadChunk); err != nil {
			return fmt.Errorf("failed to write local blob chunk: %w", err)
		}
		chunkDigests[chunkIdx] = blobtypes.MakeBlobChunkDigest(payloadChunk)
		chunkHaves[chunkIdx] = true // mark all chunks as present since we're writing the full blob
	}

	blobMeta := BlobMeta{
		payloadLength,
		chunkHaves,
		chunkDigests,
		expirySeqNr,
		bex.id,
	}
	if err := tx.WriteBlobMeta(blobDigest, blobMeta); err != nil {
		return fmt.Errorf("failed to write local blob meta: %w", err)
	}
	if err := tx.WriteStaleBlobIndex(staleBlob(expirySeqNr, blobDigest)); err != nil {
		return fmt.Errorf("failed to write stale blob index: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit kv transaction: %w", err)
	}
	return nil
}

func (bex *blobExchangeState[RI]) processBlobFetchRequest(req blobFetchRequest) {
	chDone := bex.ctx.Done()

	cert := req.cert
	err := bex.verifyCert(&cert)
	if err != nil {
		req.respond(bex.ctx, blobFetchResponse{nil, fmt.Errorf("invalid cert")})
		return
	}

	if bex.hasBlobExpired(cert.ExpirySeqNr) {
		req.respond(bex.ctx, blobFetchResponse{nil, fmt.Errorf("blob expired")})
		return
	}

	blobDigest := blobtypes.MakeBlobDigest(
		bex.config.ConfigDigest,
		cert.ChunkDigestsRoot,
		cert.PayloadLength,
		cert.ExpirySeqNr,
		cert.Submitter,
	)

	bex.logger.Debug("processing BlobFetchRequest", commontypes.LogFields{"blobDigest": blobDigest})

	seeders := make(map[commontypes.OracleID]struct{}, len(cert.AttributedBlobAvailabilitySignatures))
	for _, abs := range cert.AttributedBlobAvailabilitySignatures {
		seeders[abs.Signer] = struct{}{}
	}

	var chNotifyPayloadAvailable chan struct{}

	if existingBlob, ok := bex.blobs[blobDigest]; ok {
		if existingBlob.fetch == nil {
			chNotifyPayloadAvailable = make(chan struct{})

			existingBlob.fetch = &blobFetchMeta{
				chNotifyPayloadAvailable,
				1,
				nil,
				seeders,
				false,
			}
			if existingBlob.haveAllChunks() {
				close(chNotifyPayloadAvailable)
			}
		} else {
			for seeder := range seeders {
				existingBlob.fetch.seeders[seeder] = struct{}{} // broaden seeders per cert
			}

			existingBlob.fetch.waiters++

			chNotifyPayloadAvailable = existingBlob.fetch.chNotify
		}
	} else {
		chNotifyPayloadAvailable = make(chan struct{})

		if err := bex.createOrUpdateBlobMetaAndQuotaStats(blobDigest, cert.PayloadLength, cert.ExpirySeqNr, cert.Submitter); err != nil {
			req.respond(bex.ctx, blobFetchResponse{nil, fmt.Errorf("failed to create or update blob meta and quota stats: %w", err)})
			return
		}

		chunkDigests, chunkHaves, err := bex.loadChunkDigestsAndHaves(blobDigest, cert.PayloadLength)
		if err != nil {
			req.respond(bex.ctx, blobFetchResponse{nil, fmt.Errorf("failed to import blob chunk haves from disk: %w", err)})
			return
		}

		newBlob := &blob{
			time.Now(),
			nil,
			&blobFetchMeta{
				chNotifyPayloadAvailable,
				1,
				nil,
				seeders,
				false,
			},

			cert.ChunkDigestsRoot,
			chunkDigests,
			chunkHaves,
			cert.PayloadLength,
			cert.ExpirySeqNr,
			cert.Submitter,
		}

		bex.metrics.blobsInProgress.Inc()
		bex.blobs[blobDigest] = newBlob

		if newBlob.haveAllChunks() {
			close(chNotifyPayloadAvailable)
		}
	}

	bex.chunkRequesterGadget.PleaseRecheckPendingItems()

	bex.subs.Go(func() {
		select {
		case <-req.chDone:
		case <-chDone:
			return
		}

		select {
		case bex.chLocalEvent <- EventBlobFetchRequestDone[RI]{blobDigest}:
		case <-chDone:
		}
	})

	bex.subs.Go(func() {
		select {
		case <-chNotifyPayloadAvailable:
		case <-req.chDone:
			return
		case <-chDone:
			return
		}

		select {
		case bex.chLocalEvent <- EventBlobFetchRequestRespond[RI]{blobDigest, req}:
		case <-req.chDone:
		case <-chDone:
		}
	})
}

func (bex *blobExchangeState[RI]) initializeStatsMetricsFromKV() error {
	tx, err := bex.kv.NewReadTransactionUnchecked()
	if err != nil {
		return fmt.Errorf("failed to create read transaction: %w", err)
	}
	defer tx.Discard()

	for i, m := range bex.perOracleMetrics {
		appended, err := tx.ReadBlobQuotaStats(BlobQuotaStatsTypeAppended, commontypes.OracleID(i))
		if err != nil {
			return fmt.Errorf("failed to read appended blob quota stats for submitter %d: %w", i, err)
		}
		reaped, err := tx.ReadBlobQuotaStats(BlobQuotaStatsTypeReaped, commontypes.OracleID(i))
		if err != nil {
			return fmt.Errorf("failed to read reaped blob quota stats for submitter %d: %w", i, err)
		}
		m.SetTheirAppendedStats(appended)
		m.SetTheirReapedStats(reaped)
	}
	return nil
}

func (bex *blobExchangeState[RI]) createOrUpdateBlobMetaAndQuotaStats(
	blobDigest BlobDigest,
	payloadLength uint64,
	expirySeqNr uint64,
	submitter commontypes.OracleID,
) error {
	tx, err := bex.kv.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return fmt.Errorf("failed to create read/write transaction: %w", err)
	}
	defer tx.Discard()

	existingMeta, err := tx.ReadBlobMeta(blobDigest)
	if err != nil {
		return fmt.Errorf("failed to read blob meta: %w", err)
	}

	if existingMeta != nil {
		return nil // meta already exists, which means the quota stats already include this blob
	}

	// we must bump the quota stats
	existingQuotaStats, err := tx.ReadBlobQuotaStats(BlobQuotaStatsTypeAppended, submitter)
	if err != nil {
		return fmt.Errorf("failed to read blob quota stats: %w", err)
	}

	updatedQuotaStats, ok := existingQuotaStats.Add(BlobQuotaStats{
		1,
		payloadLength,
	})
	if !ok {
		return fmt.Errorf("quotaStats overflow")
	}

	err = tx.WriteBlobQuotaStats(BlobQuotaStatsTypeAppended, submitter, updatedQuotaStats)
	if err != nil {
		return fmt.Errorf("failed to write blob quota stats: %w", err)
	}

	cfgBlobChunkSize := bex.config.GetBlobChunkBytes()

	// and now write the meta
	blobMeta := BlobMeta{
		payloadLength,
		make([]bool, numChunks(payloadLength, cfgBlobChunkSize)),
		make([]BlobChunkDigest, numChunks(payloadLength, cfgBlobChunkSize)),
		expirySeqNr,
		submitter,
	}
	err = tx.WriteBlobMeta(blobDigest, blobMeta)
	if err != nil {
		return fmt.Errorf("failed to write blob meta: %w", err)
	}
	err = tx.WriteStaleBlobIndex(staleBlob(expirySeqNr, blobDigest))
	if err != nil {
		return fmt.Errorf("failed to write stale blob index: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	bex.perOracleMetrics[submitter].SetTheirAppendedStats(updatedQuotaStats)
	bex.perOracleMetrics[submitter].theirBlobAppendedPayloadBytes.Observe(float64(payloadLength))

	return nil
}

func (bex *blobExchangeState[RI]) eventBlobFetchRequestRespond(ev EventBlobFetchRequestRespond[RI]) {
	var (
		payload []byte
		err     error
	)
	blob, ok := bex.blobs[ev.BlobDigest]
	if ok && blob != nil && blob.fetch != nil && blob.fetch.expired {
		err = fmt.Errorf("blob expired during fetching")
	} else {
		payload, err = bex.readBlobPayload(ev.BlobDigest)
		if payload == nil && err == nil {
			err = fmt.Errorf("blob payload is unexpectedly nil")
		}
	}
	ev.Request.respond(bex.ctx, blobFetchResponse{payload, err})
}

func (bex *blobExchangeState[RI]) eventBlobFetchRequestDone(ev EventBlobFetchRequestDone[RI]) {
	blob, ok := bex.blobs[ev.BlobDigest]
	if !ok {
		return
	}
	fetch := blob.fetch
	if fetch != nil {
		fetch.weServiced()
	}
	if blob.prunable() {
		bex.metrics.blobsInProgress.Dec()
		delete(bex.blobs, ev.BlobDigest)
	}
}

func (bex *blobExchangeState[RI]) loadChunkDigestsAndHaves(blobDigest BlobDigest, payloadLength uint64) ([]BlobChunkDigest, []bool, error) {
	tx, err := bex.kv.NewReadTransactionUnchecked()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create read transaction")
	}
	defer tx.Discard()
	blobMeta, err := tx.ReadBlobMeta(blobDigest)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read blob meta: %w", err)
	}
	if blobMeta == nil {
		numChunks := numChunks(payloadLength, bex.config.GetBlobChunkBytes())
		return make([]BlobChunkDigest, numChunks), make([]bool, numChunks), nil
	}
	if blobMeta.PayloadLength != payloadLength {
		return nil, nil, fmt.Errorf("payload length mismatch: disk %d != mem %d", blobMeta.PayloadLength, payloadLength)
	}
	return blobMeta.ChunkDigests, blobMeta.ChunkHaves, nil
}

func (bex *blobExchangeState[RI]) minRejectors() int {
	return bex.config.F + 1
}

func (bex *blobExchangeState[RI]) minCertSigners() int {
	return byzquorum.Size(bex.config.N(), bex.config.F)
}

func (bex *blobExchangeState[RI]) maxCertSigners() int {
	return bex.config.N()
}

func (bex *blobExchangeState[RI]) verifyCert(cert *LightCertifiedBlob) error {
	return cert.Verify(bex.config.ConfigDigest, bex.config.OracleIdentities, bex.minCertSigners(), bex.maxCertSigners())
}

func staleBlob(expirySeqNr uint64, blobDigest BlobDigest) StaleBlob {
	return StaleBlob{expirySeqNr, blobDigest}
}
