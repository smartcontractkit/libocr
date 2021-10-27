package tracing

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2/internal/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"
)

// This file contains invariants to test against the data collected by the tracer

// CheckInvariantsIfNoTwinsNoPartitions should only pass when all the oracles in the network are honest.
func CheckInvariantsIfNoTwinsNoPartitions(t *testing.T, traces []Trace) {
	require.NoError(t, NoRoundWithoutAShouldAcceptFinalizedReport(traces))
}

// CheckInvariantsIfPartitionsButNoTwins should pass under normal operation or in the presence of a partition but without twins.
func CheckInvariantsIfPartitionsButNoTwins(t *testing.T, traces []Trace, n int) {
	require.NoError(t, AllStoredPendingTransmissionsMatch(traces, n))
	require.NoError(t, AllOraclesAcceptTheSameFinalisedReport(traces))
}

// CheckInvariantsIfTwinsAndOptionalPartitions should pass when there are twins present - they may even be leaders - as well as partitions.
func CheckInvariantsIfTwinsAndOptionalPartitions(t *testing.T, traces []Trace, f int, rMax uint8, numHonestNodes int) {
	require.NoError(t, ReportContainsAtLeast2FPlus1Observations(traces, f))
	require.NoError(t, OnlyTheLeaderRequestsObservations(traces))
	require.NoError(t, NoDuplicateNetworkMessagesExceptForNewEpoch(traces))
	require.NoError(t, NewEpochAmplificationRule(traces, f))
	require.NoError(t, ReportEpochAndRoundAreMonotonicallyIncreasing(traces))
	require.NoError(t, NetworkMessagesRoundNotLargerThanRMaxPlusOne(traces, rMax))
	require.NoError(t, AllHonestNodesComputeConfigDigestOffchain(traces, numHonestNodes))
}

// ReportContainsAtLeast2FPlus1Observations check that there are 2*f+1 records in a Broadcast message sent on the network.
func ReportContainsAtLeast2FPlus1Observations(traces []Trace, f int) error {
	for _, trace := range traces {
		broadcast, isBroadcast := trace.(*Broadcast)
		if !isBroadcast {
			continue
		}
		reportReq, isReportReq := broadcast.Message.(protocol.MessageReportReq)
		if !isReportReq {
			continue
		}
		if len(reportReq.AttributedSignedObservations) < 2*f+1 {
			return fmt.Errorf("check=ReportContains2FPlus1Observations epoch=%d, round=%d, error=expected at least %d observations but got %d, frame='%s'",
				reportReq.Epoch, reportReq.Round, 2*f+1, len(reportReq.AttributedSignedObservations), trace)
		}
	}
	return nil
}

func OnlyTheLeaderRequestsObservations(traces []Trace) error {
	// It's challenging to figure out who the leader is only from the tracing frames.
	// Instead we collect all the oracles that send a Broadcast->ObserveReq for each (epoch, round).
	//
	// If an oracle that has a twin and is also the leader of the current epoch,
	// then its twin is also acting as a leader, because they share the same id.
	// This leads to both oracles requesting observations, which breaks this test.
	var obsRequesters = make(map[protocol.EpochRound]int)
	for _, trace := range traces {
		broadcast, isBroadcast := trace.(*Broadcast)
		if !isBroadcast {
			continue
		}
		observeReq, isObserveReq := broadcast.Message.(protocol.MessageObserveReq)
		if !isObserveReq {
			continue
		}
		key := protocol.EpochRound{observeReq.Epoch, observeReq.Round}
		if _, epochRoundStarted := obsRequesters[key]; !epochRoundStarted {
			obsRequesters[key] = 0
		}
		obsRequesters[key] += 1
	}
	for key, count := range obsRequesters {
		if count != 1 {
			epoch, round := key.Epoch, key.Round
			return fmt.Errorf("check=OnlyTheLeaderRequestsObservations epoch=%d, round=%d, error=found %d oracles requesting observations", epoch, round, count)
		}
	}
	return nil
}

// NoRoundWithoutAShouldAcceptFinalizedReport for each (epoch, round), sample a ShouldAcceptFinalizedReport trace.
// There should be at least one in each timestamp.
// Note: this checker is configured to skip one timestamp. That's because when the test
// context is canceled and the oracles terminate, it might be the case that some network messages don't get through.
func NoRoundWithoutAShouldAcceptFinalizedReport(traces []Trace) error {
	var shouldAcceptEpochRounds = map[protocol.EpochRound]struct{}{}
	for _, trace := range traces {
		shouldAccept, isShouldAccept := trace.(*ShouldAcceptFinalizedReport)
		if !isShouldAccept {
			continue
		}
		shouldAcceptEpochRounds[protocol.EpochRound{
			shouldAccept.Timestamp.Epoch,
			shouldAccept.Timestamp.Round,
		}] = struct{}{}
	}
	epochRounds := []protocol.EpochRound{}
	for er := range shouldAcceptEpochRounds {
		epochRounds = append(epochRounds, er)
	}
	sort.Slice(epochRounds, func(i, j int) bool {
		a, b := epochRounds[i], epochRounds[j]
		return a.Less(b)
	})
	for i := 0; i < len(epochRounds)-1; i++ {
		a, b := epochRounds[i], epochRounds[i+1]
		if !((a.Epoch == b.Epoch && a.Round == b.Round) || // same timestamp
			(a.Epoch == b.Epoch && a.Round+1 == b.Round) || // same epoch but next round
			(a.Epoch+1 == b.Epoch && b.Round == 1) || // next epoch, first round
			(a.Epoch == b.Epoch && a.Round+2 == b.Round) || // same epoch but skip a round
			(a.Epoch+2 == b.Epoch && b.Round == 1)) { // skip two epochs, first round
			return fmt.Errorf("check=NoRoundWithoutAShouldAcceptFinalizedReport epoch=%d round=%d error=gap detected nextEpoch=%d nextRound=%d",
				a.Epoch, a.Round, b.Epoch, b.Round)
		}
	}
	return nil
}

func AllStoredPendingTransmissionsMatch(traces []Trace, n int) error {
	// All calls of StorePendingTransmissions in a (epoch, round) contain the same Transmission structure.
	var transmissions = make(map[protocol.EpochRound]*StorePendingTransmission)
	for _, trace := range traces {
		storedPendingTransmission, isStoredPendingTransmission := trace.(*StorePendingTransmission)
		if !isStoredPendingTransmission {
			continue
		}
		key := protocol.EpochRound{storedPendingTransmission.Timestamp.Epoch, storedPendingTransmission.Timestamp.Round}
		if _, oneExists := transmissions[key]; !oneExists {
			transmissions[key] = storedPendingTransmission
			continue
		}
		// TODO: make sure equality doesn't check timestamps
		if !storedPendingTransmission.EqualIgnoringTimestamp(transmissions[key]) {
			fmt.Printf(">>>>>>>>>%#v\n%#v\n", storedPendingTransmission, transmissions[key])
			return fmt.Errorf("epoch=%d round=%d found two differing stored pending transmissions t1=%s, t2=%s",
				storedPendingTransmission.Timestamp.Epoch, storedPendingTransmission.Timestamp.Round,
				storedPendingTransmission, transmissions[key])
		}
	}
	return nil
}

func NoDuplicateNetworkMessagesExceptForNewEpoch(traces []Trace) error {
	// Build a set of the string representation of each sent message (either point-to-point or broadcast) in the trace.
	// Note! We need to remove the timestamps from that string representation otherwise we will miss duplicates.
	// Note! NewEpoch is "special" in that nodes can produce this event more than once per epoch.
	var uniqueEvents = make(map[string]struct{}) // the key is the serialised message
	for _, trace := range traces {
		switch typed := trace.(type) {
		case *SendTo:
			if _, isNewEpoch := typed.Message.(protocol.MessageNewEpoch); isNewEpoch {
				continue
			}
		case *Broadcast:
			if _, isNewEpoch := typed.Message.(protocol.MessageNewEpoch); isNewEpoch {
				continue
			}
		default:
			continue // we're not interested in any other trace types
		}
		parts := strings.SplitAfterN(trace.String(), "]", 2) // serialised traces have the format "[timestamp] (originator) <eventType>...
		key := parts[1]
		if _, alreadySeen := uniqueEvents[key]; alreadySeen {
			return fmt.Errorf("check=NoDuplicateNetworkMessagesExceptForNewEpoch error=found duplicate network message '%s'", key)
		}
		uniqueEvents[key] = struct{}{}
	}
	return nil
}

func NewEpochAmplificationRule(traces []Trace, f int) error {
	// When a node receives f NewEpoch messages with a higher epoch number,
	// it will update their internal sent epoch accordingly.
	// Also, when a node's internale Tprogress expires without a progress event,
	// it will advance to the next epoch and broadcast that to other nodes.
	// In this test, we filter all the Receive->NewEpoch, Broadcast->NewEpoch and WriteState
	// and we check that any oracle that received at least f NewEpochs OR broadcasted
	// a NewEpoch itself has written the state to db.
	var countReceiveNewEpochs = make(map[string]int)     // key is "<oracleID>_<newepoch>"
	var countBroadcastedNewEpochs = make(map[string]int) // key is again "<oracleID>_<newepoch>"
	for _, trace := range traces {
		switch typed := trace.(type) {
		case *Receive:
			newEpoch, isNewEpoch := typed.Message.(protocol.MessageNewEpoch)
			if !isNewEpoch {
				continue
			}
			key := fmt.Sprintf("%s/%d", typed.Dst, newEpoch.Epoch)
			if _, found := countReceiveNewEpochs[key]; !found {
				countReceiveNewEpochs[key] = 0
			}
			countReceiveNewEpochs[key] += 1
		case *Broadcast:
			newEpoch, isNewEpoch := typed.Message.(protocol.MessageNewEpoch)
			if !isNewEpoch {
				continue
			}
			key := fmt.Sprintf("%s/%d", typed.Originator, newEpoch.Epoch)
			if _, found := countBroadcastedNewEpochs[key]; !found {
				countBroadcastedNewEpochs[key] = 0
			}
			countBroadcastedNewEpochs[key] += 1
		case *WriteState:
			highestSentEpoch := typed.State.HighestSentEpoch
			key := fmt.Sprintf("%s/%d", typed.Originator, highestSentEpoch)
			countReceives := countReceiveNewEpochs[key]
			countBroadcasts := countBroadcastedNewEpochs[key]
			if countReceives < f && countBroadcasts == 0 {
				return fmt.Errorf("check=NewEpochAmplificationRule error=node (%s) received insufficient NewEpochs %d (%d needed) OR no broadcasted NewEpochs %d, frame='%s'",
					typed.Originator, countReceives, f, countBroadcasts, trace)
			}
		default:
			continue
		}
	}
	return nil
}

func ReportEpochAndRoundAreMonotonicallyIncreasing(traces []Trace) error {
	// Check that consecutive calls to the reporting plugin's Report() method
	// have monotonically increasing timestamps for each oracle.
	prevEpochRoundByOracle := map[OracleID]protocol.EpochRound{}
	for _, trace := range traces {
		report, isReport := trace.(*Report)
		if !isReport {
			continue
		}
		currentER := protocol.EpochRound{
			report.Timestamp.Epoch,
			report.Timestamp.Round,
		}
		previousER, found := prevEpochRoundByOracle[report.Common.Originator]
		if !found {
			prevEpochRoundByOracle[report.Common.Originator] = currentER
			continue
		}
		if !previousER.Less(currentER) {
			return fmt.Errorf("check=ReportEpochAndRoundAreMonotonicallyIncreasing epoch=%d round=%d error=incompatible previous epoch=%d and round=%d, frame='%s'",
				currentER.Epoch, currentER.Round, previousER.Epoch, previousER.Round, trace)
		}
	}
	return nil
}

func NetworkMessagesRoundNotLargerThanRMaxPlusOne(traces []Trace, RMax uint8) error {
	// Checking all the network messages (either sent or broadcast) except for ObserveReq
	// have round number no bigger than RMax
	handleMessage := func(msg protocol.Message) error {
		switch typedMsg := msg.(type) {
		case protocol.MessageObserve:
			if typedMsg.Round > RMax+1 {
				return fmt.Errorf("check=NetworkMessagesRoundNotLargerThanRMax epoch=%d, round=%d, error=found MessageReport with epoch larger than RMax, RMax=%d, originalMessage=%s",
					typedMsg.Epoch, typedMsg.Round, RMax, messageAsStr(typedMsg))
			}
		case protocol.MessageReportReq:
			if typedMsg.Round > RMax+1 {
				return fmt.Errorf("check=NetworkMessagesRoundNotLargerThanRMax epoch=%d, round=%d, error=found MessageReportReq with epoch larger than RMax, RMax=%d, originalMessage=%s",
					typedMsg.Epoch, typedMsg.Round, RMax, messageAsStr(typedMsg))
			}
		case protocol.MessageReport:
			if typedMsg.Round > RMax+1 {
				return fmt.Errorf("check=NetworkMessagesRoundNotLargerThanRMax epoch=%d, round=%d, error=found MessageReport with epoch larger than RMax, RMax=%d, originalMessage=%s",
					typedMsg.Epoch, typedMsg.Round, RMax, messageAsStr(typedMsg))
			}
		case protocol.MessageFinal:
			if typedMsg.Round > RMax+1 {
				return fmt.Errorf("check=NetworkMessagesRoundNotLargerThanRMax epoch=%d, round=%d, error=found MessageFinal with epoch larger than RMax, RMax=%d, originalMessage=%s",
					typedMsg.Epoch, typedMsg.Round, RMax, messageAsStr(typedMsg))
			}
		case protocol.MessageFinalEcho:
			if typedMsg.Round > RMax+1 {
				return fmt.Errorf("check=NetworkMessagesRoundNotLargerThanRMax epoch=%d, round=%d, error=found MessageReport with epoch larger than RMax, RMax=%d, originalMessage=%s",
					typedMsg.Epoch, typedMsg.Round, RMax, messageAsStr(typedMsg))
			}
		} // Oh golang?!
		return nil
	}
	for _, trace := range traces {
		switch typedTrace := trace.(type) {
		case *SendTo:
			if err := handleMessage(typedTrace.Message); err != nil {
				return err
			}
		case *Broadcast:
			if err := handleMessage(typedTrace.Message); err != nil {
				return err
			}
		}
	}
	return nil
}

func AllOraclesAcceptTheSameFinalisedReport(traces []Trace) error {
	// Record the types.Report accepted at each timestamp in ShouldAcceptFinalizedReport
	// They should all be the same.
	mapping := map[types.ReportTimestamp]types.Report{}
	for _, trace := range traces {
		shouldAcceptFinalisedReport, isShouldAcceptFinalisedReport := trace.(*ShouldAcceptFinalizedReport)
		if !isShouldAcceptFinalisedReport {
			continue
		}
		key := shouldAcceptFinalisedReport.Timestamp
		expectedReport, isNotFirst := mapping[key]
		if !isNotFirst {
			mapping[key] = shouldAcceptFinalisedReport.Report
			expectedReport = shouldAcceptFinalisedReport.Report
		}
		if !bytes.Equal(expectedReport, shouldAcceptFinalisedReport.Report) {
			return fmt.Errorf("check=AllHonestOraclesAcceptTheSameFinalisedReport expected report for epoch=%d, round=%d to be %s but got %s from oracle id=%s",
				key.Epoch, key.Round, reportAsStr(expectedReport), reportAsStr(shouldAcceptFinalisedReport.Report),
				shouldAcceptFinalisedReport.Common.Originator)
		}
	}
	return nil
}

func AllHonestNodesComputeConfigDigestOffchain(traces []Trace, numHonestNodes int) error {
	// Count the number of calls to ConfigDigest each healthy oracle is making.
	// They should all call ConfigDigest the same number of times.
	counts := map[OracleID]int{}
	for _, trace := range traces {
		configDigest, isConfigDigest := trace.(*ConfigDigest)
		if !isConfigDigest {
			continue
		}
		if _, ok := counts[configDigest.Common.Originator]; !ok {
			counts[configDigest.Common.Originator] = 0
		}
		counts[configDigest.Common.Originator] += 1
	}
	targetCount := 0
	var targetOracleID OracleID
	for oid, count := range counts {
		if targetCount == 0 {
			targetCount = count
			targetOracleID = oid
		}
		if count != targetCount {
			return fmt.Errorf("check=AllHonestNodesComputeConfigDigestOffchain different count of calls to ConfigDigest, oracleID=%s had %d calls and oracleID=%s had %d calls",
				targetOracleID, targetCount, oid, count)
		}
	}
	return nil
}
