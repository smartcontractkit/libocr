package protocol

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/scheduler"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

func RunBlobExchange[RI any](
	ctx context.Context,

	chNetToBlobExchange <-chan MessageToBlobExchangeWithSender[RI],
	chOutcomeGenerationToBlobExchange <-chan EventToBlobExchange[RI],

	chBlobBroadcastRequest <-chan blobBroadcastRequest,
	chBlobBroadcastResponse chan<- blobBroadcastResponse,

	chBlobFetchRequest <-chan blobFetchRequest,
	chBlobFetchResponse chan<- blobFetchResponse,

	config ocr3config.SharedConfig,
	kv KeyValueStore,
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	telemetrySender TelemetrySender,
) {
	missingChunkScheduler := scheduler.NewScheduler[EventMissingBlobChunk[RI]]()
	defer missingChunkScheduler.Close()

	missingCertScheduler := scheduler.NewScheduler[EventMissingBlobCert[RI]]()
	defer missingCertScheduler.Close()

	bex := makeBlobExchangeState[RI](
		ctx, chNetToBlobExchange,
		chOutcomeGenerationToBlobExchange,
		chBlobBroadcastRequest, chBlobBroadcastResponse,
		chBlobFetchRequest, chBlobFetchResponse,
		config, kv,
		id, localConfig, logger, metricsRegisterer, netSender, offchainKeyring,
		telemetrySender,
		missingChunkScheduler, missingCertScheduler,
	)
	bex.run()
}

func makeBlobExchangeState[RI any](
	ctx context.Context,

	chNetToBlobExchange <-chan MessageToBlobExchangeWithSender[RI],
	chOutcomeGenerationToBlobExchange <-chan EventToBlobExchange[RI],

	chBlobBroadcastRequest <-chan blobBroadcastRequest,
	chBlobBroadcastResponse chan<- blobBroadcastResponse,

	chBlobFetchRequest <-chan blobFetchRequest,
	chBlobFetchResponse chan<- blobFetchResponse,

	config ocr3config.SharedConfig,
	kv KeyValueStore,
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	telemetrySender TelemetrySender,
	missingChunkScheduler *scheduler.Scheduler[EventMissingBlobChunk[RI]],
	missingCertScheduler *scheduler.Scheduler[EventMissingBlobCert[RI]],
) blobExchangeState[RI] {
	return blobExchangeState[RI]{
		ctx,
		subprocesses.Subprocesses{},

		make(chan EventToBlobExchange[RI]),
		chNetToBlobExchange,
		chOutcomeGenerationToBlobExchange,

		chBlobBroadcastRequest,
		chBlobBroadcastResponse,

		chBlobFetchRequest,
		chBlobFetchResponse,

		config,
		kv,
		id,
		localConfig,
		logger.MakeUpdated(commontypes.LogFields{"proto": "bex"}),
		netSender,
		offchainKeyring,
		telemetrySender,

		missingChunkScheduler,
		missingCertScheduler,

		make(map[BlobDigest]*blob),
	}
}

const (
	DeltaBlobChunkRequest        = 1 * time.Second
	DeltaBlobChunkRequestTimeout = 5 * time.Second
	DeltaBlobCertRequest         = 10 * time.Second
)

type blobBroadcastRequest struct {
	payload     []byte
	expirySeqNr uint64
}

type blobBroadcastResponse struct {
	chCert chan LightCertifiedBlob
	err    error
}

type blobFetchRequest struct {
	cert LightCertifiedBlob
}

type blobFetchResponse struct {
	chPayload chan []byte
	err       error
}

type blobExchangeState[RI any] struct {
	ctx  context.Context
	subs subprocesses.Subprocesses

	chLocalEvent                      chan EventToBlobExchange[RI]
	chNetToBlobExchange               <-chan MessageToBlobExchangeWithSender[RI]
	chOutcomeGenerationToBlobExchange <-chan EventToBlobExchange[RI]

	chBlobBroadcastRequest  <-chan blobBroadcastRequest
	chBlobBroadcastResponse chan<- blobBroadcastResponse

	chBlobFetchRequest  <-chan blobFetchRequest
	chBlobFetchResponse chan<- blobFetchResponse

	config          ocr3config.SharedConfig
	kv              KeyValueStore
	id              commontypes.OracleID
	localConfig     types.LocalConfig
	logger          loghelper.LoggerWithContext
	netSender       NetworkSender[RI]
	offchainKeyring types.OffchainKeyring
	telemetrySender TelemetrySender

	missingChunkScheduler *scheduler.Scheduler[EventMissingBlobChunk[RI]]
	missingCertScheduler  *scheduler.Scheduler[EventMissingBlobCert[RI]]
	blobs                 map[BlobDigest]*blob
}

const BlobChunkSize = 1 << 22 // 4MiB

type blob struct {
	chNotifyCertAvailable    chan struct{}
	chNotifyPayloadAvailable chan struct{}

	// certOrNil == nil indicates that we're the submitter and want to collect a
	// cert.
	certOrNil *LightCertifiedBlob

	chunks  []chunk
	oracles []blobOracle

	payloadLength uint64
	expirySeqNr   uint64
	submitter     commontypes.OracleID
}

type blobOracle struct {
	signature BlobAvailabilitySignature
}

type chunk struct {
	have   bool
	digest BlobChunkDigest
}

func (bex *blobExchangeState[RI]) run() {
	bex.logger.Info("BlobExchange: running", nil)

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

		case ev := <-bex.missingCertScheduler.Scheduled():
			ev.processBlobExchange(bex)
		case ev := <-bex.missingChunkScheduler.Scheduled():
			ev.processBlobExchange(bex)

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

func (bex *blobExchangeState[RI]) messageBlobOffer(msg MessageBlobOffer[RI], sender commontypes.OracleID) {
	if msg.PayloadLength == 0 {
		bex.logger.Debug("dropping MessageBlobOffer with zero payload length", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	blobDigest := blobtypes.MakeBlobDigest(
		bex.config.ConfigDigest,
		msg.ChunkDigests,
		msg.PayloadLength,
		msg.ExpirySeqNr,
		msg.Submitter,
	)

	// check if we maybe already have this blob in full
	{
		payload, err := bex.readBlobPayload(blobDigest)
		if err != nil {
			bex.logger.Warn("dropping MessageBlobOffer, failed to check if we already have the payload", commontypes.LogFields{
				"blobDigest": blobDigest,
				"sender":     sender,
			})
			return
		}
		if payload != nil {
			bex.logger.Debug("received MessageBlobOffer for which we already have the payload", commontypes.LogFields{
				"blobDigest": blobDigest,
				"sender":     sender,
			})

			bex.sendAvailabilitySignature(blobDigest, msg.Submitter)
			return
		}
	}

	if _, ok := bex.blobs[blobDigest]; ok {
		bex.logger.Debug("dropping duplicate MessageBlobOffer", commontypes.LogFields{
			"blobDigest": blobDigest,
			"sender":     sender,
		})
		return
	}

	// TODO: enforce rate limit based on sender / length
	// TODO: check payload length against Max
	// TODO: check Max against MaxMax (in plugin config)

	bex.logger.Debug("received MessageBlobOffer", commontypes.LogFields{
		"blobDigest":    blobDigest,
		"sender":        sender,
		"chunkDigests":  msg.ChunkDigests,
		"payloadLength": msg.PayloadLength,
		"expirySeqNr":   msg.ExpirySeqNr,
	})

	chunks := make([]chunk, len(msg.ChunkDigests))
	for i, chunkDigest := range msg.ChunkDigests {
		chunks[i] = chunk{
			false,
			chunkDigest,
		}
	}

	bex.blobs[blobDigest] = &blob{
		make(chan struct{}),
		make(chan struct{}),

		nil,

		chunks,
		make([]blobOracle, bex.config.N()),
		msg.PayloadLength,
		msg.ExpirySeqNr,
		msg.Submitter,
	}

	bex.missingChunkScheduler.ScheduleDelay(EventMissingBlobChunk[RI]{
		blobDigest,
	}, 0)
}

func (bex *blobExchangeState[RI]) readBlobPayload(blobDigest BlobDigest) ([]byte, error) {
	tx, err := bex.kv.NewReadTransactionUnchecked()
	if err != nil {
		return nil, fmt.Errorf("failed to create read transaction")
	}
	defer tx.Discard()
	return tx.ReadBlob(blobDigest)
}

func (bex *blobExchangeState[RI]) messageBlobChunkRequest(msg MessageBlobChunkRequest[RI], sender commontypes.OracleID) {
	blob, ok := bex.blobs[msg.BlobDigest]
	if !ok {
		bex.logger.Debug("dropping MessageBlobChunkRequest for unknown blob", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	chunkIndex := msg.ChunkIndex

	bex.logger.Debug("received MessageBlobChunkRequest", commontypes.LogFields{
		"blobDigest":    msg.BlobDigest,
		"sender":        sender,
		"chunkIndex":    chunkIndex,
		"payloadLength": blob.payloadLength,
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
	if chunk == nil {
		bex.logger.Debug("dropping MessageBlobChunkRequest, do not have chunk", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
		})
		return
	}

	chunkDigest := blobtypes.MakeBlobChunkDigest(chunk)
	expectedChunkDigest := blob.chunks[chunkIndex].digest
	if chunkDigest != expectedChunkDigest {
		bex.logger.Critical("assumption violation: chunk digest mismatch while preparing MessageBlobChunkResponse", commontypes.LogFields{
			"blobDigest":          msg.BlobDigest,
			"expectedChunkDigest": expectedChunkDigest,
			"actualChunkDigest":   chunkDigest,
		})
		return
	}

	bex.logger.Debug("sending MessageBlobChunkResponse", commontypes.LogFields{
		"blobDigest": msg.BlobDigest,
		"chunkIndex": chunkIndex,
		"to":         sender,
	})

	bex.netSender.SendTo(
		MessageBlobChunkResponse[RI]{
			msg.RequestHandle,
			msg.BlobDigest,
			chunkIndex,
			chunk,
		},
		sender,
	)
}

func (bex *blobExchangeState[RI]) messageBlobChunkResponse(msg MessageBlobChunkResponse[RI], sender commontypes.OracleID) {
	blob, ok := bex.blobs[msg.BlobDigest]
	if !ok {
		bex.logger.Debug("dropping MessageBlobChunkResponse for unknown blob", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	chunkIndex := msg.ChunkIndex

	if !(0 <= chunkIndex && chunkIndex < uint64(len(blob.chunks))) {
		bex.logger.Warn("dropping MessageBlobChunkResponse, chunk index out of range", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
			"chunkCount": len(blob.chunks),
		})
		return
	}

	if blob.chunks[chunkIndex].have {
		bex.logger.Debug("dropping MessageBlobChunkResponse, already have chunk", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
		})
		return
	}

	expectedChunkDigest := blob.chunks[chunkIndex].digest
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

	bex.logger.Debug("received MessageBlobChunkResponse", commontypes.LogFields{
		"blobDigest":    msg.BlobDigest,
		"sender":        sender,
		"chunkIndex":    chunkIndex,
		"payloadLength": blob.payloadLength,
	})

	tx, err := bex.kv.NewReadWriteTransactionUnchecked()
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

	err = tx.WriteBlobMeta(msg.BlobDigest, blob.payloadLength)
	if err != nil {
		bex.logger.Error("failed to write blob meta for MessageBlobChunkResponse", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"chunkIndex": chunkIndex,
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

	blob.chunks[chunkIndex].have = true

	for _, chunk := range blob.chunks {
		if !chunk.have {
			return
		}
	}

	bex.logger.Debug("blob fully received", commontypes.LogFields{
		"blobDigest":    msg.BlobDigest,
		"sender":        sender,
		"payloadLength": blob.payloadLength,
	})
	close(blob.chNotifyPayloadAvailable)
	bex.sendAvailabilitySignature(msg.BlobDigest, blob.submitter)
}

func (bex *blobExchangeState[RI]) messageBlobAvailable(msg MessageBlobAvailable[RI], sender commontypes.OracleID) {
	blob, ok := bex.blobs[msg.BlobDigest]
	if !ok {
		bex.logger.Debug("dropping MessageBlobAvailable for unknown blob", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	if blob.submitter != bex.id {
		bex.logger.Debug("dropping MessageBlobAvailable, not the submitter", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
			"submitter":  blob.submitter,
			"localID":    bex.id,
		})
		return
	}

	// check if we already have signature from oracle
	if len(blob.oracles[sender].signature) != 0 {
		bex.logger.Debug("dropping MessageBlobAvailable, already have signature from oracle", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	// check if we already have a certificate
	if blob.certOrNil != nil {
		bex.logger.Debug("dropping MessageBlobAvailable, already have certificate", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	// check signature
	if err := msg.Signature.Verify(msg.BlobDigest, bex.config.OracleIdentities[sender].OffchainPublicKey); err != nil {
		bex.logger.Debug("dropping MessageBlobAvailable, invalid signature", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
			"sender":     sender,
		})
		return
	}

	// save signature for oracle
	blob.oracles[sender].signature = msg.Signature

	// when we obtain certificate, notify upcall
	threshold := bex.config.N() - bex.config.F
	signers := 0
	for _, oracle := range blob.oracles {
		if len(oracle.signature) != 0 {
			signers++
		}
	}

	bex.logger.Debug("received MessageBlobAvailable", commontypes.LogFields{
		"blobDigest": msg.BlobDigest,
		"sender":     sender,
		"signers":    signers,
		"threshold":  threshold,
	})

	if signers < threshold {
		return
	}

	{
		var abass []AttributedBlobAvailabilitySignature
		for i, oracle := range blob.oracles {
			if len(oracle.signature) != 0 {
				abass = append(abass, AttributedBlobAvailabilitySignature{
					oracle.signature,
					commontypes.OracleID(i),
				})
			}
		}

		var chunkDigests []BlobChunkDigest
		for _, chunk := range blob.chunks {
			chunkDigests = append(chunkDigests, chunk.digest)
		}

		lcb := LightCertifiedBlob{
			chunkDigests,
			blob.payloadLength,
			blob.expirySeqNr,
			blob.submitter,
			abass,
		}

		if err := lcb.Verify(bex.config.ConfigDigest, bex.config.OracleIdentities, threshold); err != nil {
			bex.logger.Critical("assumption violation: failed to verify own LightCertifiedBlob", commontypes.LogFields{
				"blobDigest": msg.BlobDigest,
				"error":      err,
			})
			return
		}

		bex.logger.Debug("assembled blob availability certificate", commontypes.LogFields{
			"blobDigest": msg.BlobDigest,
		})

		blob.certOrNil = &lcb
		close(blob.chNotifyCertAvailable)
	}
}

func (bex *blobExchangeState[RI]) eventMissingChunk(ev EventMissingBlobChunk[RI]) {
	blob, ok := bex.blobs[ev.BlobDigest]
	if !ok {
		bex.logger.Debug("dropping EventMissingBlobChunk, unknown blob", commontypes.LogFields{
			"blobDigest": ev.BlobDigest,
		})
		return
	}

	for i, chunk := range blob.chunks {
		if !chunk.have {
			chunkIndex := uint64(i)

			bex.logger.Debug("broadcasting MessageBlobChunkRequest", commontypes.LogFields{
				"blobDigest":    ev.BlobDigest,
				"chunkIndex":    chunkIndex,
				"payloadLength": blob.payloadLength,
			})
			bex.netSender.Broadcast(
				MessageBlobChunkRequest[RI]{
					nil,
					ev.BlobDigest,
					chunkIndex,
				},
			)

			bex.missingChunkScheduler.ScheduleDelay(EventMissingBlobChunk[RI]{
				ev.BlobDigest,
			}, DeltaBlobChunkRequest)
			return
		}
	}
}

func (bex *blobExchangeState[RI]) eventMissingCert(ev EventMissingBlobCert[RI]) {
	blob, ok := bex.blobs[ev.BlobDigest]
	if !ok {
		bex.logger.Debug("dropping EventMissingBlobCert, unknown blob", commontypes.LogFields{
			"blobDigest": ev.BlobDigest,
		})
		return
	}

	if blob.certOrNil != nil {
		return
	}

	var chunkDigests []BlobChunkDigest
	for _, chnk := range blob.chunks {
		chunkDigests = append(chunkDigests, chnk.digest)
	}

	bex.netSender.Broadcast(MessageBlobOffer[RI]{
		chunkDigests,
		blob.payloadLength,
		blob.expirySeqNr,
		blob.submitter,
	})

	bex.missingCertScheduler.ScheduleDelay(EventMissingBlobCert[RI]{
		ev.BlobDigest,
	}, DeltaBlobCertRequest)
}

func (bex *blobExchangeState[RI]) sendAvailabilitySignature(blobDigest BlobDigest, submitter commontypes.OracleID) {

	// delete(bex.blobs, blobDigest)

	bas, err := blobtypes.MakeBlobAvailabilitySignature(blobDigest, bex.offchainKeyring.OffchainSign)
	if err != nil {
		bex.logger.Error("failed to make blob availability signature", commontypes.LogFields{
			"blobDigest": blobDigest,
			"error":      err,
		})
		return
	}

	bex.logger.Debug("sending MessageBlobAvailable", commontypes.LogFields{
		"blobDigest": blobDigest,
		"submitter":  submitter,
	})
	bex.netSender.SendTo(
		MessageBlobAvailable[RI]{
			blobDigest,
			bas,
		},
		submitter,
	)
}

func (bex *blobExchangeState[RI]) processBlobBroadcastRequest(req blobBroadcastRequest) {
	var chunkDigests []BlobChunkDigest
	var chunks []chunk
	payload := req.payload
	payloadLength := uint64(len(payload))

	for i, chunkIdx := 0, 0; i < len(payload); i, chunkIdx = i+BlobChunkSize, chunkIdx+1 {
		payloadChunk := payload[i:min(i+BlobChunkSize, len(payload))]

		// prepare for offer
		chunkDigest := blobtypes.MakeBlobChunkDigest(payloadChunk)
		chunkDigests = append(chunkDigests, chunkDigest)

		// for local accounting
		chunk := chunk{true, chunkDigest}
		chunks = append(chunks, chunk)
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

	var chNotifyCertAvailable chan struct{}
	if existingBlob, ok := bex.blobs[blobDigest]; ok {
		chNotifyCertAvailable = existingBlob.chNotifyCertAvailable
	} else {
		// if we haven't written the chunks to kv, we can't serve requests
		if err := bex.writeBlob(blobDigest, payloadLength, payload); err != nil {
			bex.chBlobBroadcastResponse <- blobBroadcastResponse{
				nil,
				fmt.Errorf("failed to write blob: %w", err),
			}
			return
		}

		// write in-memory state
		chNotifyCertAvailable = make(chan struct{})

		chNotifyPayloadAvailable := make(chan struct{})
		close(chNotifyPayloadAvailable)

		bex.blobs[blobDigest] = &blob{
			chNotifyCertAvailable,
			chNotifyPayloadAvailable,

			nil,

			chunks,
			make([]blobOracle, bex.config.N()),
			payloadLength,
			expirySeqNr,
			submitter,
		}

		bex.missingCertScheduler.ScheduleDelay(EventMissingBlobCert[RI]{
			blobDigest,
		}, 0)
	}

	chCert := make(chan LightCertifiedBlob)
	bex.chBlobBroadcastResponse <- blobBroadcastResponse{chCert, nil}

	bex.subs.Go(func() {
		chDone := bex.ctx.Done()

		select {
		case <-chNotifyCertAvailable:
			select {
			case bex.chLocalEvent <- EventRespondWithBlobCert[RI]{blobDigest, chCert}:
			case <-chDone:
			}
		case <-chDone:
		}
	})
}

func (bex *blobExchangeState[RI]) eventRespondWithBlobCert(ev EventRespondWithBlobCert[RI]) {
	blob, ok := bex.blobs[ev.BlobDigest]
	if !ok {
		bex.logger.Warn("dropping EventRespondWithBlobCert, no such blob", commontypes.LogFields{
			"blobDigest": ev.BlobDigest,
		})
		return
	}

	if blob.certOrNil == nil {
		close(ev.Channel)
		bex.logger.Critical("assumption violation: dropping EventRespondWithBlobCert, no cert available", commontypes.LogFields{
			"blobDigest": ev.BlobDigest,
		})
		return
	}

	select {
	case ev.Channel <- *blob.certOrNil:
	case <-bex.ctx.Done():
	}
}

func (bex *blobExchangeState[RI]) writeBlob(blobDigest BlobDigest, payloadLength uint64, payload []byte) error {
	tx, err := bex.kv.NewReadWriteTransactionUnchecked()
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

	if err := tx.WriteBlobMeta(blobDigest, payloadLength); err != nil {
		return fmt.Errorf("failed to write local blob meta: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit kv transaction: %w", err)
	}
	return nil
}

func (bex *blobExchangeState[RI]) processBlobFetchRequest(req blobFetchRequest) {
	threshold := bex.config.N() - bex.config.F
	cert := req.cert
	err := cert.Verify(
		bex.config.ConfigDigest,
		bex.config.OracleIdentities,
		threshold,
	)
	if err != nil {
		bex.chBlobFetchResponse <- blobFetchResponse{nil, fmt.Errorf("invalid cert")}
	}

	blobDigest := blobtypes.MakeBlobDigest(
		bex.config.ConfigDigest,
		cert.ChunkDigests,
		cert.PayloadLength,
		cert.ExpirySeqNr,
		cert.Submitter,
	)

	var chNotifyPayloadAvailable chan struct{}
	if existingBlob, ok := bex.blobs[blobDigest]; ok {
		chNotifyPayloadAvailable = existingBlob.chNotifyPayloadAvailable
	} else {
		chNotifyPayloadAvailable = make(chan struct{})
		chNotifyCertAvailable := make(chan struct{})

		var chunks []chunk
		for _, chunkDigest := range cert.ChunkDigests {
			chunks = append(chunks, chunk{false, chunkDigest})
		}

		bex.blobs[blobDigest] = &blob{
			chNotifyCertAvailable,
			chNotifyPayloadAvailable,

			&req.cert,

			chunks,
			make([]blobOracle, bex.config.N()),
			cert.PayloadLength,
			cert.ExpirySeqNr,
			cert.Submitter,
		}

		bex.missingChunkScheduler.ScheduleDelay(EventMissingBlobChunk[RI]{
			blobDigest,
		}, 0)
	}

	chPayload := make(chan []byte)
	bex.chBlobFetchResponse <- blobFetchResponse{chPayload, nil}

	bex.subs.Go(func() {
		chDone := bex.ctx.Done()

		select {
		case <-chNotifyPayloadAvailable:
			select {
			case bex.chLocalEvent <- EventRespondWithBlobPayload[RI]{blobDigest, chPayload}:
			case <-chDone:
			}
		case <-chDone:
		}
	})
}

func (bex *blobExchangeState[RI]) eventRespondWithBlobPayload(ev EventRespondWithBlobPayload[RI]) {
	payload, err := bex.readBlobPayload(ev.BlobDigest)
	if err != nil {
		close(ev.Channel)
		bex.logger.Warn("dropping EventRespondWithBlobPayload, failed to read payload", commontypes.LogFields{
			"blobDigest": ev.BlobDigest,
		})
		return
	}

	select {
	case ev.Channel <- payload:
	case <-bex.ctx.Done():
	}
}
