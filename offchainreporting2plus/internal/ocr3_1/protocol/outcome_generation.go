package protocol

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/pool"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// Identifies an instance of the outcome generation protocol
type OutcomeGenerationID struct {
	ConfigDigest types.ConfigDigest
	Epoch        uint64
}

const futureMessageBufferSize = 10 // big enough for a couple of full rounds of outgen protocol
const poolSize = 3

func RunOutcomeGeneration[RI any](
	ctx context.Context,

	chNetToOutcomeGeneration <-chan MessageToOutcomeGenerationWithSender[RI],
	chPacemakerToOutcomeGeneration <-chan EventToOutcomeGeneration[RI],
	chOutcomeGenerationToPacemaker chan<- EventToPacemaker[RI],
	chOutcomeGenerationToReportAttestation chan<- EventToReportAttestation[RI],
	chOutcomeGenerationToStateSync chan<- EventToStateSync[RI],
	blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher,
	config ocr3_1config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	kvDb KeyValueDatabase,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	reportingPlugin ocr3_1types.ReportingPlugin[RI],
	telemetrySender TelemetrySender,

	restoredCert CertifiedPrepareOrCommit,
) {

	outgen := outcomeGenerationState[RI]{
		ctx:  ctx,
		subs: subprocesses.Subprocesses{},

		chLocalEvent:                           make(chan EventToOutcomeGeneration[RI]),
		chNetToOutcomeGeneration:               chNetToOutcomeGeneration,
		chPacemakerToOutcomeGeneration:         chPacemakerToOutcomeGeneration,
		chOutcomeGenerationToPacemaker:         chOutcomeGenerationToPacemaker,
		chOutcomeGenerationToReportAttestation: chOutcomeGenerationToReportAttestation,
		chOutcomeGenerationToStateSync:         chOutcomeGenerationToStateSync,
		blobBroadcastFetcher:                   blobBroadcastFetcher,
		config:                                 config,
		database:                               database,
		id:                                     id,
		kvDb:                                   kvDb,
		localConfig:                            localConfig,
		logger:                                 logger.MakeUpdated(commontypes.LogFields{"proto": "outgen"}),
		metrics:                                newOutcomeGenerationMetrics(metricsRegisterer, logger),
		netSender:                              netSender,
		offchainKeyring:                        offchainKeyring,
		reportingPlugin:                        reportingPlugin,
		telemetrySender:                        telemetrySender,
	}
	outgen.run(restoredCert)
}

type outcomeGenerationState[RI any] struct {
	ctx  context.Context
	subs subprocesses.Subprocesses

	chLocalEvent                           chan EventToOutcomeGeneration[RI]
	chNetToOutcomeGeneration               <-chan MessageToOutcomeGenerationWithSender[RI]
	chPacemakerToOutcomeGeneration         <-chan EventToOutcomeGeneration[RI]
	chOutcomeGenerationToPacemaker         chan<- EventToPacemaker[RI]
	chOutcomeGenerationToReportAttestation chan<- EventToReportAttestation[RI]
	chOutcomeGenerationToStateSync         chan<- EventToStateSync[RI]
	blobBroadcastFetcher                   ocr3_1types.BlobBroadcastFetcher
	config                                 ocr3_1config.SharedConfig
	database                               Database
	id                                     commontypes.OracleID
	kvDb                                   KeyValueDatabase
	localConfig                            types.LocalConfig
	logger                                 loghelper.LoggerWithContext
	metrics                                *outcomeGenerationMetrics
	netSender                              NetworkSender[RI]
	offchainKeyring                        types.OffchainKeyring
	reportingPlugin                        ocr3_1types.ReportingPlugin[RI]
	telemetrySender                        TelemetrySender

	epochCtx         context.Context
	epochCtxCancel   context.CancelFunc
	bufferedMessages []*MessageBuffer[RI]
	leaderState      leaderState[RI]
	followerState    followerState[RI]
	sharedState      sharedState
}

type leaderState[RI any] struct {
	phase outgenLeaderPhase

	epochStartRequests map[commontypes.OracleID]*epochStartRequest[RI]

	readyToStartRound bool
	tRound            <-chan time.Time

	query           types.Query
	observationPool *pool.Pool[SignedObservation]
	tGrace          <-chan time.Time
}

type epochStartRequest[RI any] struct {
	message MessageEpochStartRequest[RI]
	bad     bool
}

type followerState[RI any] struct {
	phase outgenFollowerPhase

	tInitial <-chan time.Time

	leaderAbdicated bool

	roundStartPool *pool.Pool[MessageRoundStart[RI]]

	query *types.Query

	proposalPool *pool.Pool[MessageProposal[RI]]

	stateTransitionInfo stateTransitionInfo

	openKVTxn KeyValueDatabaseReadWriteTransaction

	// lock

	cert CertifiedPrepareOrCommit

	preparePool *pool.Pool[PrepareSignature]
	commitPool  *pool.Pool[CommitSignature]
}

//go-sumtype:decl stateTransitionInfo

type stateTransitionInfo interface {
	isStateTransitionInfo()
	digests() stateTransitionInfoDigests
}

type stateTransitionInfoDigests struct {
	InputsDigest               StateTransitionInputsDigest
	OutputDigest               StateTransitionOutputDigest
	StateRootDigest            StateRootDigest
	ReportsPlusPrecursorDigest ReportsPlusPrecursorDigest
}

func (stateTransitionInfoDigests) isStateTransitionInfo() {}

func (stid stateTransitionInfoDigests) digests() stateTransitionInfoDigests {
	return stid
}

type stateTransitionInfoDigestsAndPreimages struct {
	stateTransitionInfoDigests
	Outputs              StateTransitionOutputs
	ReportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor
}

func (stateTransitionInfoDigestsAndPreimages) isStateTransitionInfo() {}

func (stid stateTransitionInfoDigestsAndPreimages) digests() stateTransitionInfoDigests {
	return stid.stateTransitionInfoDigests
}

type sharedState struct {
	e uint64               // Current epoch number
	l commontypes.OracleID // Current leader number

	firstSeqNrOfEpoch uint64
	seqNr             uint64
	observationQuorum *int

	committedSeqNr uint64
}

func (outgen *outcomeGenerationState[RI]) run(restoredCert CertifiedPrepareOrCommit) {
	var restoredCommitedSeqNr uint64
	if restoredCert != nil {
		if commitQC, ok := restoredCert.(*CertifiedCommit); ok {
			restoredCommitedSeqNr = commitQC.SeqNr()
		} else if prepareQc, ok := restoredCert.(*CertifiedPrepare); ok {
			if prepareQc.SeqNr() > 1 {
				restoredCommitedSeqNr = prepareQc.SeqNr() - 1
			}
		}
	}

	outgen.logger.Info("OutcomeGeneration: running", commontypes.LogFields{
		"restoredCommittedSeqNr": restoredCommitedSeqNr,
	})

	// Initialization
	outgen.epochCtx, outgen.epochCtxCancel = context.WithCancel(outgen.ctx)

	for i := 0; i < outgen.config.N(); i++ {
		outgen.bufferedMessages = append(outgen.bufferedMessages, NewMessageBuffer[RI](futureMessageBufferSize))
	}

	outgen.leaderState = leaderState[RI]{
		outgenLeaderPhaseUnknown,
		map[commontypes.OracleID]*epochStartRequest[RI]{},
		false,
		nil,
		nil,
		nil,
		nil,
	}

	outgen.followerState = followerState[RI]{
		outgenFollowerPhaseUnknown,
		nil,
		false,
		nil,
		nil,
		nil,
		stateTransitionInfoDigests{},
		nil,
		restoredCert,
		nil,
		nil,
	}

	outgen.sharedState = sharedState{
		0,
		0,

		0,
		restoredCommitedSeqNr,
		nil,
		restoredCommitedSeqNr,
	}

	outgen.subs.Go(func() {
		RunOutcomeGenerationReap(outgen.ctx, outgen.logger, outgen.kvDb)
	})

	// Event Loop
	chDone := outgen.ctx.Done()
	for {
		select {
		case ev := <-outgen.chLocalEvent:
			ev.processOutcomeGeneration(outgen)
		case msg := <-outgen.chNetToOutcomeGeneration:
			outgen.messageToOutcomeGeneration(msg)
		case ev := <-outgen.chPacemakerToOutcomeGeneration:
			ev.processOutcomeGeneration(outgen)
		case <-outgen.followerState.tInitial:
			outgen.eventTInitialTimeout()
		case <-outgen.leaderState.tGrace:
			outgen.eventTGraceTimeout()
		case <-outgen.leaderState.tRound:
			outgen.eventTRoundTimeout()
		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			outgen.logger.Info("OutcomeGeneration: winding down", commontypes.LogFields{
				"e": outgen.sharedState.e,
				"l": outgen.sharedState.l,
			})
			outgen.subs.Wait()
			outgen.metrics.Close()
			outgen.logger.Info("OutcomeGeneration: exiting", commontypes.LogFields{
				"e": outgen.sharedState.e,
				"l": outgen.sharedState.l,
			})
			return
		default:
		}
	}
}

func (outgen *outcomeGenerationState[RI]) messageToOutcomeGeneration(msg MessageToOutcomeGenerationWithSender[RI]) {
	msgEpoch := msg.msg.epoch()
	if msgEpoch < outgen.sharedState.e {
		// drop
		outgen.logger.Debug("dropping message for past epoch", commontypes.LogFields{
			"epoch":    outgen.sharedState.e,
			"msgEpoch": msgEpoch,
			"sender":   msg.sender,
		})
	} else if msgEpoch == outgen.sharedState.e {
		msg.msg.processOutcomeGeneration(outgen, msg.sender)
	} else {
		outgen.bufferedMessages[msg.sender].Push(msg.msg)
		outgen.logger.Trace("buffering message for future epoch", commontypes.LogFields{
			"msgEpoch": msgEpoch,
			"sender":   msg.sender,
		})
	}
}

func (outgen *outcomeGenerationState[RI]) unbufferMessages() {
	outgen.logger.Trace("getting messages for new epoch", nil)
	for i, buffer := range outgen.bufferedMessages {
		sender := commontypes.OracleID(i)
		for buffer.Length() > 0 {
			msg := buffer.Peek()
			msgEpoch := msg.epoch()
			if msgEpoch < outgen.sharedState.e {
				buffer.Pop()
				outgen.logger.Debug("unbuffered and dropped message", commontypes.LogFields{
					"msgEpoch": msgEpoch,
					"sender":   sender,
				})
			} else if msgEpoch == outgen.sharedState.e {
				buffer.Pop()
				outgen.logger.Trace("unbuffered message for new epoch", commontypes.LogFields{
					"msgEpoch": msgEpoch,
					"sender":   sender,
				})
				msg.processOutcomeGeneration(outgen, sender)
			} else { // msgEpoch > e
				// this and all subsequent messages are for future epochs
				// leave them in the buffer
				break
			}
		}
	}
	outgen.logger.Trace("done unbuffering messages for new epoch", nil)
}

func (outgen *outcomeGenerationState[RI]) eventNewEpochStart(ev EventNewEpochStart[RI]) {
	// Initialization
	outgen.logger.Info("starting new epoch", commontypes.LogFields{
		"epoch": ev.Epoch,
	})

	outgen.epochCtxCancel()
	outgen.epochCtx, outgen.epochCtxCancel = context.WithCancel(outgen.ctx)

	outgen.sharedState.e = ev.Epoch
	outgen.sharedState.l = Leader(outgen.sharedState.e, outgen.config.N(), outgen.config.LeaderSelectionKey())

	outgen.logger = outgen.logger.MakeUpdated(commontypes.LogFields{
		"e": outgen.sharedState.e,
		"l": outgen.sharedState.l,
	})

	outgen.sharedState.firstSeqNrOfEpoch = 0
	outgen.sharedState.seqNr = 0

	outgen.followerState.phase = outgenFollowerPhaseNewEpoch
	outgen.followerState.tInitial = time.After(outgen.config.GetDeltaInitial())
	outgen.followerState.leaderAbdicated = false
	outgen.followerState.stateTransitionInfo = stateTransitionInfoDigests{}

	outgen.followerState.roundStartPool = pool.NewPool[MessageRoundStart[RI]](poolSize)
	outgen.followerState.proposalPool = pool.NewPool[MessageProposal[RI]](poolSize)
	outgen.followerState.preparePool = pool.NewPool[PrepareSignature](poolSize)
	outgen.followerState.commitPool = pool.NewPool[CommitSignature](poolSize)

	outgen.leaderState.phase = outgenLeaderPhaseNewEpoch
	outgen.leaderState.epochStartRequests = map[commontypes.OracleID]*epochStartRequest[RI]{}
	outgen.leaderState.readyToStartRound = false
	outgen.leaderState.observationPool = pool.NewPool[SignedObservation](1) // only one observation per sender & round, and we do not need to worry about observations from the future
	outgen.leaderState.tGrace = nil

	outgen.refreshCommittedSeqNrAndCert()
	outgen.sendStateSyncRequestFromCertifiedPrepareOrCommit(outgen.followerState.cert)

	var highestCertified CertifiedPrepareOrCommit
	var highestCertifiedTimestamp HighestCertifiedTimestamp
	highestCertified = outgen.followerState.cert
	highestCertifiedTimestamp = outgen.followerState.cert.Timestamp()

	signedHighestCertifiedTimestamp, err := MakeSignedHighestCertifiedTimestamp(
		outgen.ID(),
		highestCertifiedTimestamp,
		outgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		outgen.logger.Error("error signing timestamp", commontypes.LogFields{
			"error": err,
		})
		return
	}

	outgen.logger.Info("sending MessageEpochStartRequest to leader", commontypes.LogFields{
		"highestCertifiedTimestamp": highestCertifiedTimestamp,
	})
	outgen.netSender.SendTo(MessageEpochStartRequest[RI]{
		outgen.sharedState.e,
		highestCertified,
		signedHighestCertifiedTimestamp,
	}, outgen.sharedState.l)

	if outgen.id == outgen.sharedState.l {
		outgen.leaderState.tRound = time.After(outgen.config.DeltaRound)
	}

	outgen.unbufferMessages()
}

func (outgen *outcomeGenerationState[RI]) ID() OutcomeGenerationID {
	return OutcomeGenerationID{outgen.config.ConfigDigest, outgen.sharedState.e}
}

func (outgen *outcomeGenerationState[RI]) RoundCtx(seqNr uint64) RoundContext {
	if seqNr != outgen.sharedState.committedSeqNr+1 {
		outgen.logger.Critical("assumption violation, seqNr isn't successor to committedSeqNr", commontypes.LogFields{
			"seqNr":          seqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
		panic("")
	}
	return RoundContext{
		seqNr,
		outgen.sharedState.e,
		seqNr - outgen.sharedState.firstSeqNrOfEpoch + 1,
	}
}

func callPluginFromOutcomeGenerationBackground[T any](
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	name string,
	recommendedMaxDuration time.Duration,
	roundCtx RoundContext,
	f func(context.Context, RoundContext) (T, error),
) (T, bool) {
	return common.CallPluginFromBackground[T](
		ctx,
		logger,
		commontypes.LogFields{
			"seqNr": roundCtx.SeqNr,
			"round": roundCtx.Round, // nolint: staticcheck
		},
		name,
		recommendedMaxDuration,
		func(ctx context.Context) (T, error) {
			return f(ctx, roundCtx)
		},
	)
}

func (outgen *outcomeGenerationState[RI]) sendStateSyncRequestFromCertifiedPrepareOrCommit(cert CertifiedPrepareOrCommit) {
	var seqNr uint64
	switch cert := cert.(type) {
	case *CertifiedPrepare:
		seqNr = cert.PrepareSeqNr - 1
	case *CertifiedCommit:
		seqNr = cert.CommitSeqNr
	}

	select {
	case outgen.chOutcomeGenerationToStateSync <- EventStateSyncRequest[RI]{
		seqNr,
	}:
	case <-outgen.ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) tryToMoveCertAndKVStateToCommitQC(commitQC *CertifiedCommit) {
	ok := outgen.commit(*commitQC)
	if !ok {
		outgen.logger.Error("commit() failed", commontypes.LogFields{
			"commitQCSeqNr": commitQC.CommitSeqNr,
		})
		// We intentionally fall through to give a chance to advance the KV
		// state regardless of the cert persistence failure.
	}

	// Early return to avoid unnecessary error logs
	{
		committedKVSeqNr, err := outgen.committedKVSeqNr()
		if err != nil {
			outgen.logger.Error("failed to read highest committed kv seq nr", commontypes.LogFields{
				"error": err,
			})
			return
		}
		if committedKVSeqNr == commitQC.CommitSeqNr {
			// already where we want to be
			return
		}
		if committedKVSeqNr+1 != commitQC.CommitSeqNr {
			outgen.logger.Debug("commit qc seq nr quite ahead from committed kv seq nr, can't do anything", commontypes.LogFields{
				"committedKVSeqNr": committedKVSeqNr,
				"commitQCSeqNr":    commitQC.CommitSeqNr,
			})
			return
		}
	}

	tx, err := outgen.kvDb.NewSerializedReadWriteTransaction(commitQC.CommitSeqNr)
	if err != nil {

		outgen.logger.Error("failed to create serialized transaction", commontypes.LogFields{
			"error": err,
		})
		return
	}
	defer tx.Discard()

	committedKVSeqNr, err := tx.ReadHighestCommittedSeqNr()
	if err != nil {
		outgen.logger.Error("failed to read highest committed kv seq nr", commontypes.LogFields{
			"error": err,
		})
		return
	}

	if commitQC.CommitSeqNr == committedKVSeqNr {
		outgen.logger.Debug("kv state already at commit qc seq nr, nothing to do", commontypes.LogFields{
			"committedKVSeqNr": committedKVSeqNr,
			"commitQCSeqNr":    commitQC.CommitSeqNr,
		})
		return
	}

	if commitQC.CommitSeqNr != committedKVSeqNr+1 {
		outgen.logger.Debug("commit qc seq nr quite ahead from committed kv seq nr, can't do anything", commontypes.LogFields{
			"committedKVSeqNr": committedKVSeqNr,
			"commitQCSeqNr":    commitQC.CommitSeqNr,
		})
		return
	}

	stb, err := tx.ReadUnattestedStateTransitionBlock(commitQC.CommitSeqNr, commitQC.StateTransitionInputsDigest)
	if err != nil {
		outgen.logger.Error("error during ReadUnattestedStateTransitionBlock", commontypes.LogFields{
			"commitQCSeqNr": commitQC.CommitSeqNr,
			"error":         err,
		})
		return
	}
	if stb == nil {
		outgen.logger.Debug("unattested state transition block not found, can't move kv state", commontypes.LogFields{
			"commitQCSeqNr": commitQC.CommitSeqNr,
		})
		return
	}
	if err := outgen.isCompatibleUnattestedStateTransitionBlockSanityCheck(commitQC, *stb); err != nil {
		outgen.logger.Critical("sanity check of unattested state transition block failed, very surprising!", commontypes.LogFields{
			"commitQCSeqNr": commitQC.CommitSeqNr,
			"error":         err,
		})
		return
	}

	astb := AttestedStateTransitionBlock{
		*stb,
		commitQC.CommitQuorumCertificate,
	}

	// write astb
	err = tx.WriteAttestedStateTransitionBlock(commitQC.CommitSeqNr, astb)
	if err != nil {
		outgen.logger.Error("error writing attested state transition block", commontypes.LogFields{
			"commitQCSeqNr": commitQC.CommitSeqNr,
			"error":         err,
		})
		return
	}

	// apply write set
	stateRootDigest, err := tx.ApplyWriteSet(stb.StateTransitionOutputs.WriteSet)
	if err != nil {
		outgen.logger.Error("error applying write set", commontypes.LogFields{
			"commitQCSeqNr": commitQC.CommitSeqNr,
			"error":         err,
		})
		return
	}

	if stateRootDigest != stb.StateRootDigest {
		outgen.logger.Error("state root digest mismatch from write set application", commontypes.LogFields{
			"commitQCSeqNr": commitQC.CommitSeqNr,
			"expected":      stb.StateRootDigest,
			"actual":        stateRootDigest,
		})
		return
	}

	err = tx.Commit()
	if err != nil {
		outgen.logger.Error("error committing transaction", commontypes.LogFields{
			"commitQCSeqNr": commitQC.CommitSeqNr,
			"error":         err,
		})
		return
	}

	outgen.logger.Debug("successfully moved kv state to commit qc", commontypes.LogFields{
		"oldCommittedKVSeqNr": committedKVSeqNr,
		"newCommittedKVSeqNr": commitQC.CommitSeqNr,
	})
}

func (outgen *outcomeGenerationState[RI]) persistUnattestedStateTransitionBlockAndReportsPlusPrecursor(stb StateTransitionBlock, stateTransitionInputsDigest StateTransitionInputsDigest, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) error {
	kvTxn, err := outgen.kvDb.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return err
	}
	defer kvTxn.Discard()
	seqNr := stb.BlockSeqNr
	err = kvTxn.WriteUnattestedStateTransitionBlock(seqNr, stateTransitionInputsDigest, stb)
	if err != nil {
		return fmt.Errorf("failed to write unattested state transition block: %w", err)
	}
	err = kvTxn.WriteReportsPlusPrecursor(seqNr, stb.ReportsPlusPrecursorDigest, reportsPlusPrecursor)
	if err != nil {
		return fmt.Errorf("failed to write reports plus precursor: %w", err)
	}
	err = kvTxn.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}

func (outgen *outcomeGenerationState[RI]) isCompatibleUnattestedStateTransitionBlockSanityCheck(commitQC *CertifiedCommit, stb StateTransitionBlock) error {
	stbStateTransitionOutputsDigest := MakeStateTransitionOutputDigest(
		outgen.config.ConfigDigest,
		stb.BlockSeqNr,
		stb.StateTransitionOutputs.WriteSet,
	)

	if stbStateTransitionOutputsDigest != commitQC.StateTransitionOutputsDigest {
		return fmt.Errorf("local state transition block outputs digest does not match commitQC: expected %s but got %s", commitQC.StateTransitionOutputsDigest, stbStateTransitionOutputsDigest)
	}
	if stb.StateRootDigest != commitQC.StateRootDigest {
		return fmt.Errorf("local state transition block state root digest does not match commitQC: expected %s but got %s", commitQC.StateRootDigest, stb.StateRootDigest)
	}
	if stb.ReportsPlusPrecursorDigest != commitQC.ReportsPlusPrecursorDigest {
		return fmt.Errorf("local state transition block reportsPlusPrecursor digest does not match commitQC: expected %s but got %s", commitQC.ReportsPlusPrecursorDigest, stb.ReportsPlusPrecursorDigest)
	}
	return nil
}

func (outgen *outcomeGenerationState[RI]) committedKVSeqNr() (uint64, error) {
	return committedKVSeqNr(outgen.kvDb)
}

func committedKVSeqNr(kvDb KeyValueDatabase) (uint64, error) {
	tx, err := kvDb.NewReadTransactionUnchecked()
	if err != nil {
		return 0, fmt.Errorf("failed to create read transaction: %w", err)
	}
	defer tx.Discard()
	return tx.ReadHighestCommittedSeqNr()
}
