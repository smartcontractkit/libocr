package protocol

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
	"sort"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Pacemaker keeps track of the state and message handling for an oracle
// participating in the off-chain reporting protocol
func RunPacemaker[RI any](
	ctx context.Context,

	chNetToPacemaker <-chan MessageToPacemakerWithSender[RI],
	chPacemakerToOutcomeGeneration chan<- EventToOutcomeGeneration[RI],
	chOutcomeGenerationToPacemaker <-chan EventToPacemaker[RI],
	config ocr3config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	telemetrySender TelemetrySender,

	restoredEpoch uint64,
) {
	pace := makePacemakerState[RI](
		ctx, chNetToPacemaker,
		chPacemakerToOutcomeGeneration, chOutcomeGenerationToPacemaker,
		config, database,
		id, localConfig, logger, netSender, offchainKeyring,
		telemetrySender,
	)
	pace.run(restoredEpoch)
}

func makePacemakerState[RI any](
	ctx context.Context,
	chNetToPacemaker <-chan MessageToPacemakerWithSender[RI],
	chPacemakerToOutcomeGeneration chan<- EventToOutcomeGeneration[RI],
	chOutcomeGenerationToPacemaker <-chan EventToPacemaker[RI],
	config ocr3config.SharedConfig,
	database Database, id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	telemetrySender TelemetrySender,
) pacemakerState[RI] {
	return pacemakerState[RI]{
		ctx: ctx,

		chNetToPacemaker:               chNetToPacemaker,
		chPacemakerToOutcomeGeneration: chPacemakerToOutcomeGeneration,
		chOutcomeGenerationToPacemaker: chOutcomeGenerationToPacemaker,
		config:                         config,
		database:                       database,
		id:                             id,
		localConfig:                    localConfig,
		logger:                         logger.MakeUpdated(commontypes.LogFields{"proto": "pacemaker"}),
		netSender:                      netSender,
		offchainKeyring:                offchainKeyring,
		telemetrySender:                telemetrySender,

		newepoch: make([]uint64, config.N()),
	}
}

type pacemakerState[RI any] struct {
	ctx context.Context

	chNetToPacemaker               <-chan MessageToPacemakerWithSender[RI]
	chPacemakerToOutcomeGeneration chan<- EventToOutcomeGeneration[RI]
	chOutcomeGenerationToPacemaker <-chan EventToPacemaker[RI]
	config                         ocr3config.SharedConfig
	database                       Database
	id                             commontypes.OracleID
	localConfig                    types.LocalConfig
	logger                         loghelper.LoggerWithContext
	netSender                      NetworkSender[RI]
	offchainKeyring                types.OffchainKeyring
	telemetrySender                TelemetrySender
	// Test use only: send testBlocker an event to halt the pacemaker event loop,
	// send testUnblocker an event to resume it.
	testBlocker   chan eventTestBlock
	testUnblocker chan eventTestUnblock

	// ne is the highest epoch number this oracle has broadcast in a newepoch
	// message, during the current epoch
	ne uint64

	// e is the number of the current epoch
	e uint64

	// l is the index of the leader for the current epoch
	l commontypes.OracleID

	// newepoch[j] is the highest epoch number oracle j has sent in a newepoch
	// message, during the current epoch.
	newepoch []uint64

	// tResend is a timeout used by the leader-election protocol to
	// periodically resend the latest Newepoch message in order to
	// guard against unreliable network conditions
	tResend <-chan time.Time

	// tProgress is a timeout used by the leader-election protocol to track
	// whether the current leader is making adequate progress.
	tProgress <-chan time.Time

	notifyOutcomeGenerationOfNewEpoch bool
}

func (pace *pacemakerState[RI]) run(restoredEpoch uint64) {
	pace.logger.Info("Running Pacemaker", nil)

	// Initialization

	// rounds start with 1, so let's make epochs also start with 1
	// this also gives us cleaner behavior for the initial epoch, which is otherwise
	// immediately terminated and superseded due to restoreNeFromTransmitter below

	pace.ne = 1
	pace.e = 1
	pace.l = Leader(pace.e, pace.config.N(), pace.config.LeaderSelectionKey())

	if pace.e <= restoredEpoch {
		pace.ne = restoredEpoch + 1
	}

	pace.tProgress = time.After(pace.config.DeltaProgress)

	pace.sendNewepoch(pace.ne)

	pace.notifyOutcomeGenerationOfNewEpoch = true

	// Initialization complete

	// Take a reference to the ctx.Done channel once, here, to avoid taking the
	// context lock below.
	chDone := pace.ctx.Done()

	// Event Loop
	for {
		var nilOrChPacemakerToOutcomeGeneration chan<- EventToOutcomeGeneration[RI]
		if pace.notifyOutcomeGenerationOfNewEpoch {
			nilOrChPacemakerToOutcomeGeneration = pace.chPacemakerToOutcomeGeneration
		} else {
			nilOrChPacemakerToOutcomeGeneration = nil
		}

		select {
		case nilOrChPacemakerToOutcomeGeneration <- EventNewEpochStart[RI]{pace.e}:
			pace.notifyOutcomeGenerationOfNewEpoch = false
		case msg := <-pace.chNetToPacemaker:
			msg.msg.processPacemaker(pace, msg.sender)
		case ev := <-pace.chOutcomeGenerationToPacemaker:
			ev.processPacemaker(pace)
		case <-pace.tResend:
			pace.eventTResendTimeout()
		case <-pace.tProgress:
			pace.eventTProgressTimeout()
		case <-pace.testBlocker:
			<-pace.testUnblocker
		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			pace.logger.Info("Pacemaker: winding down", nil)

			pace.logger.Info("Pacemaker: exiting", nil)
			return
		default:
		}
	}
}

// eventProgress is called when a "progress" event is emitted by the reporting
// prototol. It resets the timer which will trigger the oracle to broadcast a
// "newepoch" message, if it runs out.
func (pace *pacemakerState[RI]) eventProgress() {
	pace.tProgress = time.After(pace.config.DeltaProgress)
}

func (pace *pacemakerState[RI]) sendNewepoch(newEpoch uint64) {
	pace.netSender.Broadcast(MessageNewEpoch[RI]{newEpoch})
	if pace.ne != newEpoch {
		pace.ne = newEpoch
	}
	pace.tResend = time.After(pace.config.DeltaResend)
}

func (pace *pacemakerState[RI]) eventTResendTimeout() {
	pace.sendNewepoch(pace.ne)
}

func (pace *pacemakerState[RI]) eventTProgressTimeout() {
	pace.logger.Debug("Pacemaker: TProgress expired", nil)
	pace.eventChangeLeader()
}

func (pace *pacemakerState[RI]) eventChangeLeader() {
	pace.tProgress = nil
	sendEpoch := pace.ne
	epochPlusOne := pace.e + 1
	if epochPlusOne <= pace.e {
		pace.logger.Error("Pacemaker: epoch overflows, cannot change leader", nil)
		return
	}

	if sendEpoch < epochPlusOne {
		sendEpoch = epochPlusOne
	}
	pace.sendNewepoch(sendEpoch)
}

func (pace *pacemakerState[RI]) messageNewepoch(msg MessageNewEpoch[RI], sender commontypes.OracleID) {
	if pace.newepoch[sender] < msg.Epoch {
		pace.newepoch[sender] = msg.Epoch
	} else {
		// neither of the following two "upon" handlers can be triggered
		return
	}

	// upon |{p_j ∈ P | newepoch[j] > ne}| > f do
	{
		candidateEpochs := sortedGreaterThan(pace.newepoch, pace.ne)
		if len(candidateEpochs) > pace.config.F {
			// ē ← max {e' | {p_j ∈ P | newepoch[j] ≥ e' } > f}
			newEpoch := candidateEpochs[len(candidateEpochs)-(pace.config.F+1)]
			pace.sendNewepoch(newEpoch)
		}
	}

	// upon |{p_j ∈ P | newepoch[j] > e}| > 2f do
	{
		candidateEpochs := sortedGreaterThan(pace.newepoch, pace.e)
		if len(candidateEpochs) > 2*pace.config.F {
			// ē ← max {e' | {p_j ∈ P | newepoch[j] ≥ e' } > 2f}
			//
			// since candidateEpochs contains, in increasing order, the epochs from
			// the received newepoch messages, this value of newEpoch was sent by at
			// least 2F+1 processes
			newEpoch := candidateEpochs[len(candidateEpochs)-(2*pace.config.F+1)]
			pace.logger.Debug("Moving to epoch, based on candidateEpochs", commontypes.LogFields{
				"newEpoch":        newEpoch,
				"candidateEpochs": candidateEpochs,
			})
			l := Leader(newEpoch, pace.config.N(), pace.config.LeaderSelectionKey())
			pace.e, pace.l = newEpoch, l // (e, l) ← (ē, leader(ē))
			if pace.ne < pace.e {        // ne ← max{ne, e}
				pace.ne = pace.e
			}

			pace.notifyOutcomeGenerationOfNewEpoch = true

			pace.tProgress = time.After(pace.config.DeltaProgress) // restart timer T_{progress}
		}
	}
}

// sortedGreaterThan returns the *sorted* elements of xs which are greater than y
func sortedGreaterThan(xs []uint64, y uint64) (rv []uint64) {
	for _, x := range xs {
		if x > y {
			rv = append(rv, x)
		}
	}
	sort.Slice(rv, func(i, j int) bool { return rv[i] < rv[j] })
	return rv
}

// Leader will produce an oracle id for the given epoch.
func Leader(epoch uint64, n int, key [16]byte) (leader commontypes.OracleID) {
	mac := hmac.New(sha256.New, key[:])
	_ = binary.Write(mac, binary.BigEndian, epoch)

	r := big.NewInt(0).SetBytes(mac.Sum(nil))
	// This is biased, but we don't care because the prob of us hitting the bias are
	// less than 2**5/2**256 = 2**-251.
	r.Mod(r, big.NewInt(int64(n)))
	return commontypes.OracleID(r.Int64())
}

type eventTestBlock struct{}
type eventTestUnblock struct{}
