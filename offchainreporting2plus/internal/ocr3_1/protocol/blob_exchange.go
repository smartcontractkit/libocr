package protocol

import (
	"context"
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/byzquorum"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common/scheduler"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/protocol/requestergadget"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/subprocesses"
	"github.com/prometheus/client_golang/prometheus"
)

func RunBlobExchange[RI any](
	ctx context.Context,

	chNetToBlobExchange <-chan MessageToBlobExchangeWithSender[RI],
	chOutcomeGenerationToBlobExchange <-chan EventToBlobExchange[RI],

	chBlobBroadcastRequest <-chan blobBroadcastRequest,
	chBlobFetchRequest <-chan blobFetchRequest,

	config ocr3config.SharedConfig,
	kv KeyValueDatabase,
	id commontypes.OracleID,
	limits ocr3_1types.ReportingPluginLimits,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	telemetrySender TelemetrySender,
) {
	broadcastGraceTimeoutScheduler := scheduler.NewScheduler[EventBlobBroadcastGraceTimeout[RI]]()
	defer broadcastGraceTimeoutScheduler.Close()

	bex := makeBlobExchangeState[RI](
		ctx, chNetToBlobExchange,
		chOutcomeGenerationToBlobExchange,
		chBlobBroadcastRequest, chBlobFetchRequest,
		config, kv,
		id, limits, localConfig, logger, metricsRegisterer, netSender, offchainKeyring,
		telemetrySender,
		broadcastGraceTimeoutScheduler,
	)
	bex.run()
}

func makeBlobExchangeState[RI any](
	ctx context.Context,

	chNetToBlobExchange <-chan MessageToBlobExchangeWithSender[RI],
	chOutcomeGenerationToBlobExchange <-chan EventToBlobExchange[RI],

	chBlobBroadcastRequest <-chan blobBroadcastRequest,
	chBlobFetchRequest <-chan blobFetchRequest,

	config ocr3config.SharedConfig,
	kv KeyValueDatabase,
	id commontypes.OracleID,
	limits ocr3_1types.ReportingPluginLimits,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	telemetrySender TelemetrySender,

	broadcastGraceTimeoutScheduler *scheduler.Scheduler[EventBlobBroadcastGraceTimeout[RI]],
) *blobExchangeState[RI] {
	tStopExpiredBlobFetches := time.After(DeltaStopExpiredBlobFetches)

	bex := &blobExchangeState[RI]{
		ctx,
		subprocesses.Subprocesses{},

		make(chan EventToBlobExchange[RI]),
		chNetToBlobExchange,
		chOutcomeGenerationToBlobExchange,

		chBlobBroadcastRequest,
		chBlobFetchRequest,

		config,
		kv,
		id,
		limits,
		localConfig,
		logger.MakeUpdated(commontypes.LogFields{"proto": "bex"}),
		netSender,
		offchainKeyring,
		telemetrySender,

		broadcastGraceTimeoutScheduler,
		nil, // must be filled right below

		nil, // must be filled right below
		tStopExpiredBlobFetches,

		make(map[BlobDigest]*blob),
	}

	offerRequesterGadget := requestergadget.NewRequesterGadget[blobOfferItem](
		config.N(),
		DeltaBlobOfferBroadcast,
		bex.trySendBlobOffer,
		bex.getPendingBlobOffers,
		bex.getBlobOfferSeeders,
	)
	bex.offerRequesterGadget = offerRequesterGadget

	chunkRequesterGadget := requestergadget.NewRequesterGadget[blobChunkId](
		config.N(),
		DeltaBlobChunkRequest,
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

	bex.logger.Debug("sending MessageBlobChunkRequest", commontypes.LogFields{
		"blobDigest": id.blobDigest,
		"chunkIndex": id.chunkIndex,
		"seeder":     seeder,
	})

	chunkSize := blob.getBlobChunkSize(id.chunkIndex)
	expiryTimestamp := time.Now().Add(blobChunkRequestExpiration(chunkSize))
	bex.netSender.SendTo(MessageBlobChunkRequest[RI]{
		nil,
		&MessageBlobChunkRequestInfo{
			expiryTimestamp,
		},
		id.blobDigest,
		id.chunkIndex,
	}, seeder)

	return &requestergadget.RequestInfo{
		expiryTimestamp,
	}, true
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

	bex.logger.Trace("sending MessageBlobOffer", commontypes.LogFields{
		"blobDigest":    item.blobDigest,
		"chunkDigests":  blob.chunkDigests,
		"payloadLength": blob.payloadLength,
		"expirySeqNr":   blob.expirySeqNr,
		"to":            seeder,
	})

	expiryTimestamp := time.Now().Add(blobOfferBroadcastExpiration(blob.payloadLength))

	bex.netSender.SendTo(MessageBlobOffer[RI]{
		nil,
		&MessageBlobOfferRequestInfo{
			expiryTimestamp,
		},
		blob.chunkDigests,
		blob.payloadLength,
		blob.expirySeqNr,
	}, seeder)

	return &requestergadget.RequestInfo{
		expiryTimestamp,
	}, true
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
	rateBytesPerSecond = 10 * 1024 * 1024 // 10 MiB/s
	latencyOverhead    = 1 * time.Second

	// DeltaBlobChunkRequest denotes the minimum duration between sending two
	// MessageBlobChunkRequest messages to a particular oracle.
	DeltaBlobChunkRequest = 10 * time.Millisecond

	// DeltaBlobOfferBroadcast denotes the minimum duration between sending two
	// MessageBlobOffer messages to a particular oracle.
	DeltaBlobOfferBroadcast = 10 * time.Millisecond

	// DeltaBlobBroadcastGrace denotes the duration that we will wait after
	// receiving minSigners valid accepting MessageBlobOfferResponse messages,
	// to give a last chance to straggling oracles to send us a
	// MessageBlobOfferResponse.
	DeltaBlobBroadcastGrace = 100 * time.Millisecond

	// DeltaStopExpiredBlobFetches denotes the interval with which we check for
	// in-progress blob fetches for blobs that might have expired, and mark them
	// as expired and/or send reject MessageBlobOfferResponse to the submitter
	// if appropriate.
	DeltaStopExpiredBlobFetches = 5 * time.Second
)

func transmitDataDuration(rateBytesPerSecond int, size uint64) time.Duration {
	secs := float64(size) / float64(rateBytesPerSecond)
	return time.Duration(secs * float64(time.Second))
}

func blobChunkRequestExpiration(chunkSize uint64) time.Duration {
	return latencyOverhead + transmitDataDuration(rateBytesPerSecond, chunkSize)
}

func blobOfferBroadcastExpiration(payloadLength uint64) time.Duration {
	const latencyPerChunk = DeltaBlobChunkRequest

	expiration := latencyOverhead
	for i := uint64(0); i < payloadLength; i += BlobChunkSize {
		chunkSize := min(BlobChunkSize, payloadLength-i)
		expiration += latencyPerChunk
		expiration += transmitDataDuration(rateBytesPerSecond, chunkSize)
	}
	return expiration
}

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

	chLocalEvent                      chan EventToBlobExchange[RI]
	chNetToBlobExchange               <-chan MessageToBlobExchangeWithSender[RI]
	chOutcomeGenerationToBlobExchange <-chan EventToBlobExchange[RI]

	chBlobBroadcastRequest <-chan blobBroadcastRequest
	chBlobFetchRequest     <-chan blobFetchRequest

	config          ocr3config.SharedConfig
	kv              KeyValueDatabase
	id              commontypes.OracleID
	limits          ocr3_1types.ReportingPluginLimits
	localConfig     types.LocalConfig
	logger          loghelper.LoggerWithContext
	netSender       NetworkSender[RI]
	offchainKeyring types.OffchainKeyring
	telemetrySender TelemetrySender

	// blob broadcast
	broadcastGraceTimeoutScheduler *scheduler.Scheduler[EventBlobBroadcastGraceTimeout[RI]]
	offerRequesterGadget           *requestergadget.RequesterGadget[blobOfferItem]

	// blob fetch
	chunkRequesterGadget    *requestergadget.RequesterGadget[blobChunkId]
	tStopExpiredBlobFetches <-chan time.Time

	blobs map[BlobDigest]*blob
}

type blobOfferItem struct {
	blobDigest BlobDigest
	oracleID   commontypes.OracleID
}

type blobChunkId struct {
	blobDigest BlobDigest
	chunkIndex uint64
}

const BlobChunkSize = 1 << 22 // 4MiB

func numChunks(payloadLength uint64) uint64 {
	return (payloadLength + BlobChunkSize - 1) / BlobChunkSize
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

	chunkDigests []BlobChunkDigest
	chunkHaves   []bool

	payloadLength uint64
	expirySeqNr   uint64
	submitter     commontypes.OracleID
}

func (b *blob) getBlobChunkSize(chunkIndex uint64) uint64 {
	if chunkIndex == uint64(len(b.chunkDigests))-1 {
		return b.payloadLength % BlobChunkSize
	}
	return BlobChunkSize
}

func (b *blob) haveAllChunks() bool {
	return !slices.Contains(b.chunkHaves, false)
}

func (b *blob) prunable() bool {
	return b.broadcast.prunable() && b.fetch.prunable()
}

func (bex *blobExchangeState[RI]) run() {
	bex.logger.Info("BlobExchange: running", nil)

	bex.subs.Go(func() {
		RunBlobReap(bex.ctx, bex.logger, bex.kv)
	})

	// Take a reference to the ctx.Done channel once, here, to avoid taking the
	// context lock below.
	chDone := bex.ctx.Done()

	// Event Loop
	for {
		select {
		case ev := <-bex.chLocalEvent:
			ev.processBlobExchange(bex)

		case msg := <-bex.chNetToBlobExchange:
			msg.msg.processBlobExchange(bex, msg.sender)
		case ev := <-bex.chOutcomeGenerationToBlobExchange:
			ev.processBlobExchange(bex)

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
		case <-bex.tStopExpiredBlobFetches:
			bex.eventTStopExpiredBlobFetches()

		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			bex.logger.Info("BlobExchange: winding down", nil)
			bex.subs.Wait()
			// bex.metrics.Close()
			bex.logger.Info("BlobExchange: exiting", nil)
			return
		default:
		}
	}
}

func (bex *blobExchangeState[RI]) eventTStopExpiredBlobFetches() {
	defer func() {
		bex.tStopExpiredBlobFetches = time.After(DeltaStopExpiredBlobFetches)
	}()

	tx, err := bex.kv.NewReadTransactionUnchecked()
	if err != nil {
		bex.logger.Error("failed to create read transaction for eventTStopExpiredBlobFetches", commontypes.LogFields{
			"error": err,
		})
		return
	}
	defer tx.Discard()

	highestCommittedSeqNr, err := tx.ReadHighestCommittedSeqNr()
	if err != nil {
		bex.logger.Error("failed to read highest committed seq nr for eventTStopExpiredBlobFetches", commontypes.LogFields{
			"error": err,
		})
		return
	}

	for blobDigest, blob := range bex.blobs {
		fetch := blob.fetch
		if fetch == nil {
			continue
		}
		if fetch.expired || blob.haveAllChunks() {
			continue
		}

		if !hasBlobExpired(blob.expirySeqNr, highestCommittedSeqNr) {
			continue
		}

		bex.logger.Debug("stopping expired blob fetch", commontypes.LogFields{
			"blobDigest":            blobDigest,
			"expirySeqNr":           blob.expirySeqNr,
			"highestCommittedSeqNr": highestCommittedSeqNr,
			"submitter":             blob.submitter,
		})

		if fetch.exchange != nil {
			bex.sendBlobOfferResponseRejecting(blobDigest, blob.submitter, fetch.exchange.latestOfferRequestHandle)
			fetch.exchange.weServiced()
		}

		fetch.expired = true
		close(fetch.chNotify)

		if blob.prunable() {
			delete(bex.blobs, blobDigest)
		}
	}
}

func (bex *blobExchangeState[RI]) messageBlobOffer(msg MessageBlobOffer[RI], sender commontypes.OracleID) {
	submitter := sender

	blobDigest := blobtypes.MakeBlobDigest(
		bex.config.ConfigDigest,
		msg.ChunkDigests,
		msg.PayloadLength,
		msg.ExpirySeqNr,
		submitter,
	)

	chunkHaves, err := bex.loadChunkHaves(blobDigest, msg.PayloadLength)
	if err != nil {
		bex.logger.Warn("dropping MessageBlobOffer, failed to check if we already have the payload", commontypes.LogFields{
			"blobDigest": blobDigest,
			"sender":     sender,
		})
		return
	}

	// check if we maybe already have this blob in full
	if !slices.Contains(chunkHaves, false) {
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

	// Reject if payload length exceeds maximum allowed length
	if msg.PayloadLength > uint64(bex.limits.MaxBlobPayloadLength) {
		bex.logger.Debug("received MessageBlobOffer with payload length that exceeds maximum allowed length, rejecting", commontypes.LogFields{
			"blobDigest":       blobDigest,
			"submitter":        submitter,
			"payloadLength":    msg.PayloadLength,
			"maxPayloadLength": bex.limits.MaxBlobPayloadLength,
		})
		bex.sendBlobOfferResponseRejecting(blobDigest, submitter, msg.RequestHandle)
		return
	}

	// Reject if blob has already expired
	committedSeqNr, err := bex.kv.HighestCommittedSeqNr()
	if err != nil {
		bex.logger.Error("failed to read highest committed seq nr for MessageBlobOffer", commontypes.LogFields{
			"error": err,
		})
		return
	}
	if hasBlobExpired(msg.ExpirySeqNr, committedSeqNr) {
		bex.logger.Debug("received MessageBlobOffer for already expired blob, rejecting", commontypes.LogFields{
			"blobDigest":     blobDigest,
			"submitter":      submitter,
			"expirySeqNr":    msg.ExpirySeqNr,
			"committedSeqNr": committedSeqNr,
		})
		bex.sendBlobOfferResponseRejecting(blobDigest, submitter, msg.RequestHandle)
		return
	}

	// TODO: enforce rate limit based on sender / length

	bex.logger.Debug("received MessageBlobOffer", commontypes.LogFields{
		"blobDigest":    blobDigest,
		"sender":        sender,
		"chunkDigests":  msg.ChunkDigests,
		"payloadLength": msg.PayloadLength,
		"expirySeqNr":   msg.ExpirySeqNr,
	})

	seeders := map[commontypes.OracleID]struct{}{
		submitter: {},
	}

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

		msg.ChunkDigests,
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
	} else {
		// save rejection for oracle
		broadcast.oracles[sender] = blobBroadcastOracleMeta{
			true,
			false,
			nil,
		}
	}

	threshold := bex.minCertSigners()

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
		"threshold":        threshold,
	})

	if broadcast.phase == blobBroadcastPhaseAcceptedGrace {
		return
	}

	if acceptingOracles >= threshold {
		bex.logger.Debug("minimum number of accepting oracles reached, entering grace period", commontypes.LogFields{
			"acceptingOracles": acceptingOracles,
			"threshold":        threshold,
			"blobDigest":       msg.BlobDigest,
			"gracePeriod":      DeltaBlobBroadcastGrace,
		})
		broadcast.phase = blobBroadcastPhaseAcceptedGrace
		bex.broadcastGraceTimeoutScheduler.ScheduleDelay(EventBlobBroadcastGraceTimeout[RI]{
			msg.BlobDigest,
		}, DeltaBlobBroadcastGrace)
		return
	}

	if rejectingOracles >= threshold {
		bex.logger.Warn("oracle quorum rejected our broadcast", commontypes.LogFields{
			"rejectingOracles": rejectingOracles,
			"blobDigest":       msg.BlobDigest,
		})
		broadcast.phase = blobBroadcastPhaseRejected
		close(broadcast.chNotify)

		return
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

	shuffledOracles := make([]commontypes.OracleID, 0, bex.config.N())
	for i := range bex.config.N() {
		shuffledOracles = append(shuffledOracles, commontypes.OracleID(i))
	}

	rand.Shuffle(len(shuffledOracles), func(i, j int) {
		shuffledOracles[i], shuffledOracles[j] = shuffledOracles[j], shuffledOracles[i]
	})

	var abass []AttributedBlobAvailabilitySignature
	for _, oracleID := range shuffledOracles {
		oracle := broadcast.oracles[oracleID]
		if oracle.weReceivedOfferResponse && oracle.weReceivedOfferResponseAccepting && len(abass) < maxSigners {
			abass = append(abass, AttributedBlobAvailabilitySignature{
				oracle.signature,
				oracleID,
			})
		}
	}

	lcb := LightCertifiedBlob{
		blob.chunkDigests,
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
	defer tx.Discard()
	if err != nil {
		bex.logger.Error("failed to create read transaction for MessageBlobChunkRequest", commontypes.LogFields{
			"error": err,
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
		bex.chunkRequesterGadget.MarkBadResponse(bcid, sender)
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

	if blob.chunkHaves[chunkIndex] {
		bex.logger.Debug("dropping MessageBlobChunkResponse, already have chunk", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
		})
		bex.chunkRequesterGadget.MarkBadResponse(bcid, sender)
		return
	}

	expectedChunkDigest := blob.chunkDigests[chunkIndex]
	actualChunkDigest := blobtypes.MakeBlobChunkDigest(msg.Chunk)
	if expectedChunkDigest != actualChunkDigest {
		bex.logger.Debug("dropping MessageBlobChunkResponse, chunk digest mismatch", commontypes.LogFields{
			"blobDigest":     msg.BlobDigest,
			"sender":         sender,
			"chunkIndex":     chunkIndex,
			"expectedDigest": expectedChunkDigest,
			"actualDigest":   actualChunkDigest,
		})
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
	blobMeta := BlobMeta{
		blob.payloadLength,
		chunkHaves,
		blob.expirySeqNr,
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
		delete(bex.blobs, msg.BlobDigest)
	}
}

func (bex *blobExchangeState[RI]) processBlobBroadcastRequest(req blobBroadcastRequest) {
	if len(req.payload) > bex.limits.MaxBlobPayloadLength {
		req.respond(bex.ctx, blobBroadcastResponse{
			LightCertifiedBlob{},
			fmt.Errorf("blob payload length %d exceeds maximum allowed length %d",
				len(req.payload), bex.limits.MaxBlobPayloadLength),
		})
		return
	}

	payload := req.payload
	payloadLength := uint64(len(payload))

	chunkDigests := make([]BlobChunkDigest, 0, numChunks(payloadLength))
	chunkHaves := make([]bool, 0, numChunks(payloadLength))

	for i, chunkIdx := 0, 0; i < len(payload); i, chunkIdx = i+BlobChunkSize, chunkIdx+1 {
		payloadChunk := payload[i:min(i+BlobChunkSize, len(payload))]

		// prepare for offer
		chunkDigest := blobtypes.MakeBlobChunkDigest(payloadChunk)
		chunkDigests = append(chunkDigests, chunkDigest)

		// for local accounting
		chunkHaves = append(chunkHaves, true)
	}

	expirySeqNr := req.expirySeqNr
	submitter := bex.id

	blobDigest := blobtypes.MakeBlobDigest(
		bex.config.ConfigDigest,
		chunkDigests,
		payloadLength,
		expirySeqNr,
		submitter,
	)

	bex.logger.Debug("processing BlobBroadcastRequest", commontypes.LogFields{"blobDigest": blobDigest})

	var chNotifyCertAvailable chan struct{}
	if existingBlob, ok := bex.blobs[blobDigest]; ok {
		if existingBlob.broadcast == nil {
			existingBlob.broadcast = &blobBroadcastMeta{
				make(chan struct{}),
				1,
				blobBroadcastPhaseOffering,
				nil,
				make([]blobBroadcastOracleMeta, bex.config.N()),
			}
		}
		chNotifyCertAvailable = existingBlob.broadcast.chNotify
	} else {
		// if we haven't written the chunks to kv, we can't serve requests

		if err := bex.writeBlob(blobDigest, payloadLength, payload, expirySeqNr); err != nil {
			req.respond(bex.ctx, blobBroadcastResponse{
				LightCertifiedBlob{},
				fmt.Errorf("failed to write blob: %w", err),
			})
			return
		}

		// write in-memory state
		chNotifyCertAvailable = make(chan struct{})

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

			chunkDigests,
			chunkHaves,
			payloadLength,
			expirySeqNr,
			submitter,
		}

		bex.offerRequesterGadget.PleaseRecheckPendingItems()
	}

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
		delete(bex.blobs, ev.BlobDigest)
	}
}

func (bex *blobExchangeState[RI]) writeBlob(blobDigest BlobDigest, payloadLength uint64, payload []byte, expirySeqNr uint64) error {
	tx, err := bex.kv.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return fmt.Errorf("failed to create read/write transaction: %w", err)
	}
	defer tx.Discard()
	for i, chunkIdx := 0, uint64(0); i < len(payload); i, chunkIdx = i+BlobChunkSize, chunkIdx+1 {
		payloadChunk := payload[i:min(i+BlobChunkSize, len(payload))]
		if err := tx.WriteBlobChunk(blobDigest, chunkIdx, payloadChunk); err != nil {
			return fmt.Errorf("failed to write local blob chunk: %w", err)
		}
	}

	chunksHave := make([]bool, numChunks(payloadLength))
	for i := range chunksHave {
		chunksHave[i] = true // mark all chunks as present since we're writing the full blob
	}
	blobMeta := BlobMeta{
		payloadLength,
		chunksHave,
		expirySeqNr,
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

	blobDigest := blobtypes.MakeBlobDigest(
		bex.config.ConfigDigest,
		cert.ChunkDigests,
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

		if !existingBlob.haveAllChunks() {
			bex.chunkRequesterGadget.PleaseRecheckPendingItems()
		}
	} else {
		chNotifyPayloadAvailable = make(chan struct{})

		chunkHaves, err := bex.loadChunkHaves(blobDigest, cert.PayloadLength)
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

			cert.ChunkDigests,
			chunkHaves,
			cert.PayloadLength,
			cert.ExpirySeqNr,
			cert.Submitter,
		}

		bex.blobs[blobDigest] = newBlob

		if newBlob.haveAllChunks() {
			close(chNotifyPayloadAvailable)
		} else {
			bex.chunkRequesterGadget.PleaseRecheckPendingItems()
		}
	}

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
		delete(bex.blobs, ev.BlobDigest)
	}
}

func (bex *blobExchangeState[RI]) loadChunkHaves(blobDigest BlobDigest, payloadLength uint64) ([]bool, error) {
	tx, err := bex.kv.NewReadTransactionUnchecked()
	if err != nil {
		return nil, fmt.Errorf("failed to create read transaction")
	}
	defer tx.Discard()
	blobMeta, err := tx.ReadBlobMeta(blobDigest)
	if err != nil {
		return nil, fmt.Errorf("failed to read blob meta: %w", err)
	}
	if blobMeta == nil {
		return make([]bool, numChunks(payloadLength)), nil
	}
	if blobMeta.PayloadLength != payloadLength {
		return nil, fmt.Errorf("payload length mismatch: disk %d != mem %d", blobMeta.PayloadLength, payloadLength)
	}
	return blobMeta.ChunksHave, nil
}

func (bex *blobExchangeState[RI]) minCertSigners() int {
	return bex.config.F + 1
}

func (bex *blobExchangeState[RI]) maxCertSigners() int {

	return byzquorum.Size(bex.config.N(), bex.config.F)
}

func (bex *blobExchangeState[RI]) verifyCert(cert *LightCertifiedBlob) error {
	return cert.Verify(bex.config.ConfigDigest, bex.config.OracleIdentities, bex.minCertSigners(), bex.maxCertSigners())
}

func staleBlob(expirySeqNr uint64, blobDigest BlobDigest) StaleBlob {
	return StaleBlob{expirySeqNr + 1, blobDigest}
}

func hasBlobExpired(expirySeqNr uint64, committedSeqNr uint64) bool {
	return expirySeqNr < committedSeqNr
}
