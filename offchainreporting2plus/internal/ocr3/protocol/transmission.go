package protocol

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol/persist"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/permutation"
	"github.com/smartcontractkit/libocr/subprocesses"
	"golang.org/x/crypto/sha3"
)

const ContractTransmitterTimeoutWarningGracePeriod = 50 * time.Millisecond

const chPersistCapacityTransmission = 16

// TransmissionProtocol tracks the local oracle process's role in the transmission of a
// report to the on-chain oracle contract.
//
// Note: The transmission protocol doesn't clean up pending transmissions
// when it is terminated. This is by design, but means that old pending
// transmissions may accumulate in the database. They should be garbage
// collected once in a while.
func RunTransmission[RI any](
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	chReportFinalizationToTransmission <-chan EventToTransmission[RI],
	config ocr3config.SharedConfig,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	database ocr3types.Database,
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	reportingPlugin ocr3types.OCR3Plugin[RI],
) {
	t := transmissionState[RI]{
		ctx:          ctx,
		subprocesses: subprocesses,

		chReportFinalizationToTransmission: chReportFinalizationToTransmission,
		config:                             config,
		contractTransmitter:                contractTransmitter,
		database:                           database,
		id:                                 id,
		localConfig:                        localConfig,
		logger:                             logger,
		reportingPlugin:                    reportingPlugin,
	}
	t.run()
}

type transmissionState[RI any] struct {
	ctx          context.Context
	subprocesses *subprocesses.Subprocesses

	chReportFinalizationToTransmission <-chan EventToTransmission[RI]
	config                             ocr3config.SharedConfig
	contractTransmitter                ocr3types.ContractTransmitter[RI]
	database                           ocr3types.Database
	id                                 commontypes.OracleID
	localConfig                        types.LocalConfig
	logger                             loghelper.LoggerWithContext
	reportingPlugin                    ocr3types.OCR3Plugin[RI]

	chPersist chan<- persist.TransmissionDBUpdate
	times     MinHeapTimeToPendingTransmission[RI]
	tTransmit <-chan time.Time
}

// run runs the event loop for the local transmission protocol
func (t *transmissionState[RI]) run() {
	// t.restoreFromDatabase()

	chPersist := make(chan persist.TransmissionDBUpdate, chPersistCapacityTransmission)
	t.chPersist = chPersist
	t.subprocesses.Go(func() {
		persist.PersistTransmission(
			t.ctx,
			chPersist,
			t.database,
			t.localConfig.DatabaseTimeout,
			t.logger,
		)
	})

	chDone := t.ctx.Done()
	for {
		select {
		case ev := <-t.chReportFinalizationToTransmission:
			ev.processTransmission(t)
		case <-t.tTransmit:
			t.eventTTransmitTimeout()
		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			t.logger.Info("Transmission: exiting", nil)
			return
		default:
		}
	}
}

// func (t *transmissionState[RI]) restoreFromDatabase() {
//
// 	// childCtx, childCancel := context.WithTimeout(t.ctx, t.localConfig.DatabaseTimeout)
// 	// defer childCancel()
// 	// pending, err := t.database.PendingTransmissionsWithConfigDigest(childCtx, t.config.ConfigDigest)
// 	// if err != nil {
// 	// 	t.logger.ErrorIfNotCanceled("Transmission: error fetching pending transmissions from database", childCtx, commontypes.LogFields{"error": err})
// 	// 	return
// 	// }

// 	// now := time.Now()

// 	// // insert non-expired transmissions into queue
// 	// for key, trans := range pending {
// 	// 	if now.Before(trans.Time) {
// 	// 		t.times.Push(MinHeapTimeToPendingTransmissionItem{
// 	// 			key,
// 	// 			trans,
// 	// 		})
// 	// 	}
// 	// }

// 	// // if queue isn't empty, set tTransmit to expire at next transmission time
// 	// if t.times.Len() != 0 {
// 	// 	next := t.times.Peek()
// 	// 	t.tTransmit = time.After(time.Until(next.Time))
// 	// }
// }

// eventTransmit is called when the local process sends a transmit event
func (t *transmissionState[RI]) eventTransmit(ev EventTransmit[RI]) {
	t.logger.Debug("Received transmit event", commontypes.LogFields{
		"seqNr":     ev.SeqNr,
		"index":     ev.Index,
		"reportLen": len(ev.AttestedReport.ReportWithInfo.Report),
	})

	{
		ctx, cancel := context.WithTimeout(t.ctx, t.config.MaxDurationShouldAcceptFinalizedReport)
		defer cancel()

		ins := loghelper.NewIfNotStopped(
			t.config.MaxDurationShouldAcceptFinalizedReport+ReportingPluginTimeoutWarningGracePeriod,
			func() {
				t.logger.Error("Transmission: OCR3Plugin.ShouldAcceptFinalizedReport is taking too long", commontypes.LogFields{
					"seqNr":       ev.SeqNr,
					"index":       ev.Index,
					"maxDuration": t.config.MaxDurationShouldAcceptFinalizedReport,
				})
			},
		)

		shouldAccept, err := t.reportingPlugin.ShouldAcceptFinalizedReport(
			ctx,
			ev.SeqNr,
			ev.AttestedReport.ReportWithInfo,
		)

		ins.Stop()

		if err != nil {
			t.logger.Error("eventTransmit(ev): error in ReportingPlugin.ShouldAcceptFinalizedReport", commontypes.LogFields{
				"error": err,
				"seqNr": ev.SeqNr,
				"index": ev.Index,
			})
			return
		}

		if !shouldAccept {
			t.logger.Debug("eventTransmit(ev): ReportingPlugin.ShouldAcceptFinalizedReport returned false", commontypes.LogFields{
				"seqNr": ev.SeqNr,
				"index": ev.Index,
			})
			return
		}
	}

	now := time.Now()
	delayMaybe := t.transmitDelay(ev.SeqNr, ev.Index)
	if delayMaybe == nil {
		return
	}
	delay := *delayMaybe

	// transmission := types.PendingTransmission{
	// 	now.Add(delay),
	// 	ev.H,
	// 	ev.AttestedReport.Report,
	// 	ev.AttestedReport.AttributedSignatures,
	// }

	// select {
	// case t.chPersist <- persist.TransmissionDBUpdate{ts, &transmission}:
	// default:
	// 	t.logger.Warn("eventTransmit: chPersist is overflowing", nil)
	// }

	t.times.Push(MinHeapTimeToPendingTransmissionItem[RI]{
		now.Add(delay),
		ev.SeqNr,
		ev.Index,
		ev.AttestedReport,
	})

	next := t.times.Peek()
	if ev.SeqNr == next.SeqNr && ev.Index == next.Index {
		t.tTransmit = time.After(delay)
	}
}

func (t *transmissionState[RI]) eventTTransmitTimeout() {
	defer func() {
		// if queue isn't empty, set tTransmit to expire at next transmission time
		if t.times.Len() != 0 {
			next := t.times.Peek()
			t.tTransmit = time.After(time.Until(next.Time))
		}
	}()

	if t.times.Len() == 0 {
		return
	}
	item := t.times.Pop()

	// select {
	// case t.chPersist <- persist.TransmissionDBUpdate{
	// 	types.ReportTimestamp{
	// 		t.config.ConfigDigest,
	// 		item.Epoch,
	// 		item.Round,
	// 	},
	// 	nil,
	// }:
	// default:
	// 	t.logger.Warn("eventTTransmitTimeout: chPersist is overflowing", nil)
	// }

	{
		ctx, cancel := context.WithTimeout(
			t.ctx,
			t.config.MaxDurationShouldTransmitAcceptedReport,
		)
		defer cancel()

		ins := loghelper.NewIfNotStopped(
			t.config.MaxDurationShouldTransmitAcceptedReport+ReportingPluginTimeoutWarningGracePeriod,
			func() {
				t.logger.Error("Transmission: ReportingPlugin.ShouldTransmitAcceptedReport is taking too long", commontypes.LogFields{
					"maxDuration": t.config.MaxDurationShouldTransmitAcceptedReport,
					"seqNr":       item.SeqNr,
					"index":       item.Index,
				})
			},
		)

		shouldTransmit, err := t.reportingPlugin.ShouldTransmitAcceptedReport(
			ctx,
			item.SeqNr,
			item.AttestedReport.ReportWithInfo,
		)

		ins.Stop()

		if err != nil {
			t.logger.Error("eventTTransmitTimeout: ReportingPlugin.ShouldTransmitAcceptedReport error", commontypes.LogFields{
				"error": err,
				"seqNr": item.SeqNr,
				"index": item.Index,
			})
			return
		}

		if !shouldTransmit {
			t.logger.Info("eventTTransmitTimeout: ReportingPlugin.ShouldTransmitAcceptedReport returned false", commontypes.LogFields{
				"seqNr": item.SeqNr,
				"index": item.Index,
			})
			return
		}
	}

	t.logger.Info("eventTTransmitTimeout: Transmitting", commontypes.LogFields{
		"seqNr": item.SeqNr,
		"index": item.Index,
	})

	{
		ctx, cancel := context.WithTimeout(
			t.ctx,
			t.localConfig.ContractTransmitterTransmitTimeout,
		)
		defer cancel()

		ins := loghelper.NewIfNotStopped(
			t.localConfig.ContractTransmitterTransmitTimeout+ContractTransmitterTimeoutWarningGracePeriod,
			func() {
				t.logger.Error("Transmission: ContractTransmitter.Transmit is taking too long", commontypes.LogFields{
					"maxDuration": t.localConfig.ContractTransmitterTransmitTimeout,
					"seqNr":       item.SeqNr,
					"index":       item.Index,
				})
			},
		)

		err := t.contractTransmitter.Transmit(
			ctx,
			t.config.ConfigDigest,
			item.AttestedReport.ReportWithInfo,
			item.AttestedReport.AttributedSignatures,
		)

		ins.Stop()

		if err != nil {
			t.logger.Error("eventTTransmitTimeout: ContractTransmitter.Transmit error", commontypes.LogFields{"error": err})
			return
		}

	}

	t.logger.Info("eventTTransmitTimeout:❗️successfully invoked ContractTransmitter.Transmit", commontypes.LogFields{
		"seqNr": item.SeqNr,
		"index": item.Index,
	})
}

func (t *transmissionState[RI]) transmitDelay(seqNr uint64, index int) *time.Duration {
	// No need for HMAC. Since we use Keccak256, prepending
	// with key gives us a PRF already.
	hash := sha3.NewLegacyKeccak256()
	transmissionOrderKey := t.config.TransmissionOrderKey()
	hash.Write(transmissionOrderKey[:])
	hash.Write(t.config.ConfigDigest[:])
	temp := make([]byte, 8)
	binary.LittleEndian.PutUint64(temp, seqNr)
	hash.Write(temp)
	binary.LittleEndian.PutUint64(temp, uint64(index))
	hash.Write(temp)

	var key [16]byte
	copy(key[:], hash.Sum(nil))
	pi := permutation.Permutation(t.config.N(), key)

	sum := 0
	for i, s := range t.config.S {
		sum += s
		if pi[t.id] < sum {
			result := time.Duration(i) * t.config.DeltaStage
			return &result
		}
	}
	return nil
}
