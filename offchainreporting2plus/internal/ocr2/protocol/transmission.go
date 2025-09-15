package protocol

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr2config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr2/protocol/persist"
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
func RunTransmission(
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	config ocr2config.SharedConfig,
	chReportFinalizationToTransmission <-chan EventToTransmission,
	database types.Database,
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	reportingPlugin types.ReportingPlugin,
	telemetrySender TelemetrySender,
	transmitter types.ContractTransmitter,
) {
	t := transmissionState{
		ctx:          ctx,
		subprocesses: subprocesses,

		config:                             config,
		chReportFinalizationToTransmission: chReportFinalizationToTransmission,
		database:                           database,
		id:                                 id,
		localConfig:                        localConfig,
		logger:                             logger,
		reportingPlugin:                    reportingPlugin,
		telemetrySender:                    telemetrySender,
		transmitter:                        transmitter,
	}
	t.run()
}

type transmissionState struct {
	ctx          context.Context
	subprocesses *subprocesses.Subprocesses

	config                             ocr2config.SharedConfig
	chReportFinalizationToTransmission <-chan EventToTransmission
	database                           types.Database
	id                                 commontypes.OracleID
	localConfig                        types.LocalConfig
	logger                             loghelper.LoggerWithContext
	reportingPlugin                    types.ReportingPlugin
	telemetrySender                    TelemetrySender
	transmitter                        types.ContractTransmitter

	chPersist chan<- persist.TransmissionDBUpdate
	times     MinHeapTimeToPendingTransmission
	tTransmit <-chan time.Time
}

// run runs the event loop for the local transmission protocol
func (t *transmissionState) run() {
	t.restoreFromDatabase()

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

func (t *transmissionState) restoreFromDatabase() {
	childCtx, childCancel := context.WithTimeout(t.ctx, t.localConfig.DatabaseTimeout)
	defer childCancel()
	pending, err := t.database.PendingTransmissionsWithConfigDigest(childCtx, t.config.ConfigDigest)
	if err != nil {
		t.logger.ErrorIfNotCanceled("Transmission: error fetching pending transmissions from database", childCtx, commontypes.LogFields{"error": err})
		return
	}

	now := time.Now()

	// insert non-expired transmissions into queue
	for key, trans := range pending {
		if now.Before(trans.Time) {
			t.times.Push(MinHeapTimeToPendingTransmissionItem{
				key,
				trans,
			})
		}
	}

	// if queue isn't empty, set tTransmit to expire at next transmission time
	if t.times.Len() != 0 {
		next := t.times.Peek()
		t.tTransmit = time.After(time.Until(next.Time))
	}
}

// eventTransmit is called when the local process sends a transmit event
func (t *transmissionState) eventTransmit(ev EventTransmit) {
	t.logger.Debug("Received transmit event", commontypes.LogFields{
		"epoch": ev.Epoch,
		"round": ev.Round,
	})

	ts := types.ReportTimestamp{t.config.ConfigDigest, ev.Epoch, ev.Round}

	{
		ctx, cancel := context.WithTimeout(t.ctx, t.config.MaxDurationShouldAcceptFinalizedReport)
		defer cancel()

		ins := loghelper.NewIfNotStopped(
			t.config.MaxDurationShouldAcceptFinalizedReport+ReportingPluginTimeoutWarningGracePeriod,
			func() {
				t.logger.Error("Transmission: ReportingPlugin.ShouldAcceptFinalizedReport is taking too long", commontypes.LogFields{
					"epoch": ev.Epoch, "round": ev.Round, "maxDuration": t.config.MaxDurationShouldAcceptFinalizedReport.String(),
				})
			},
		)

		shouldAccept, err := t.reportingPlugin.ShouldAcceptFinalizedReport(
			ctx,
			ts,
			ev.AttestedReport.Report,
		)

		ins.Stop()

		if err != nil {
			t.logger.Error("eventTransmit(ev): error in ReportingPlugin.ShouldAcceptFinalizedReport", commontypes.LogFields{
				"error": err,
				"epoch": ev.Epoch,
				"round": ev.Round,
			})
			t.telemetrySender.TransmissionShouldAcceptFinalizedReportComputed(ts, false, false)
			return
		}

		if !shouldAccept {
			t.logger.Debug("eventTransmit(ev): ReportingPlugin.ShouldAcceptFinalizedReport returned false", commontypes.LogFields{
				"epoch": ev.Epoch,
				"round": ev.Round,
			})
			t.telemetrySender.TransmissionShouldAcceptFinalizedReportComputed(ts, false, true)
			return
		}
	}

	t.telemetrySender.TransmissionShouldAcceptFinalizedReportComputed(ts, true, true)

	now := time.Now()
	delays := t.transmitDelays(ev.Epoch, ev.Round)
	t.telemetrySender.TransmissionScheduleComputed(ts, now, delays)
	delay, ok := delays[t.id]
	if !ok {
		return
	}

	transmission := types.PendingTransmission{
		now.Add(delay),
		ev.H,
		ev.AttestedReport.Report,
		ev.AttestedReport.AttributedSignatures,
	}

	select {
	case t.chPersist <- persist.TransmissionDBUpdate{ts, &transmission}:
	default:
		t.logger.Warn("eventTransmit: chPersist is overflowing", nil)
	}

	t.times.Push(MinHeapTimeToPendingTransmissionItem{ts, transmission})

	next := t.times.Peek()
	if (EpochRound{ev.Epoch, ev.Round}) == (EpochRound{next.Epoch, next.Round}) {
		t.tTransmit = time.After(delay)
	}
}

func (t *transmissionState) eventTTransmitTimeout() {
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

	select {
	case t.chPersist <- persist.TransmissionDBUpdate{
		types.ReportTimestamp{
			t.config.ConfigDigest,
			item.Epoch,
			item.Round,
		},
		nil,
	}:
	default:
		t.logger.Warn("eventTTransmitTimeout: chPersist is overflowing", nil)
	}

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
					"item": item, "maxDuration": t.config.MaxDurationShouldTransmitAcceptedReport.String(),
				})
			},
		)

		shouldTransmit, err := t.reportingPlugin.ShouldTransmitAcceptedReport(
			ctx,
			item.ReportTimestamp,
			item.Report,
		)

		ins.Stop()

		if err != nil {
			t.logger.Error("eventTTransmitTimeout: ReportingPlugin.ShouldTransmitAcceptedReport error", commontypes.LogFields{"error": err})
			t.telemetrySender.TransmissionShouldTransmitAcceptedReportComputed(item.ReportTimestamp, false, false)
			return
		}

		if !shouldTransmit {
			t.logger.Info("eventTTransmitTimeout: ReportingPlugin.ShouldTransmitAcceptedReport returned false", nil)
			t.telemetrySender.TransmissionShouldTransmitAcceptedReportComputed(item.ReportTimestamp, false, true)
			return
		}
	}

	t.logger.Info("eventTTransmitTimeout: Transmitting", commontypes.LogFields{
		"epoch": item.Epoch,
		"round": item.Round,
	})
	t.telemetrySender.TransmissionShouldTransmitAcceptedReportComputed(item.ReportTimestamp, true, true)

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
					"item": item, "maxDuration": t.localConfig.ContractTransmitterTransmitTimeout.String(),
				})
			},
		)

		err := t.transmitter.Transmit(
			ctx,
			types.ReportContext{
				item.ReportTimestamp,
				item.ExtraHash,
			},
			item.Report,
			item.AttributedSignatures,
		)

		ins.Stop()

		if err != nil {
			t.logger.Error("eventTTransmitTimeout: ContractTransmitter.Transmit error", commontypes.LogFields{"error": err})
			return
		}

	}

	t.logger.Info("eventTTransmitTimeout:❗️successfully called Transmit", commontypes.LogFields{
		"epoch": item.Epoch,
		"round": item.Round,
	})
}

// Computes a map from oracle ids to to transmission delays. This is
// deterministic across all oracles. The result is derived pseudorandomly
// uniformly and independently per epoch and round.
func (t *transmissionState) transmitDelays(epoch uint32, round uint8) map[commontypes.OracleID]time.Duration {
	// No need for HMAC. Since we use Keccak256, prepending
	// with key gives us a PRF already.
	hash := sha3.NewLegacyKeccak256()
	transmissionOrderKey := t.config.TransmissionOrderKey()
	hash.Write(transmissionOrderKey[:])
	hash.Write(t.config.ConfigDigest[:])
	temp := make([]byte, 8)
	binary.LittleEndian.PutUint64(temp, uint64(epoch))
	hash.Write(temp)
	binary.LittleEndian.PutUint64(temp, uint64(round))
	hash.Write(temp)

	var key [16]byte
	copy(key[:], hash.Sum(nil))

	// Permutation from transmission order index to oracle id
	piInv := make([]int, t.config.N())
	{
		// Permutation from oracle id to transmission order index. The
		// permutations are structured in this "inverted" way for historical
		// compatibility
		pi := permutation.Permutation(t.config.N(), key)
		for i := range pi {
			piInv[pi[i]] = i
		}
	}

	result := make(map[commontypes.OracleID]time.Duration, t.config.N())

	accumulatedStageSize := 0
	for stageIdx, stageSize := range t.config.S {
		// i is the index of the oracle sorted by transmission order
		for i := accumulatedStageSize; i < accumulatedStageSize+stageSize; i++ {
			if i >= len(piInv) {
				// Index is larger than index of the last oracle. This happens
				// when sum(S) > N.
				break
			}
			oracleId := commontypes.OracleID(piInv[i])
			result[oracleId] = time.Duration(stageIdx) * t.config.DeltaStage
		}

		accumulatedStageSize += stageSize
	}

	return result
}
