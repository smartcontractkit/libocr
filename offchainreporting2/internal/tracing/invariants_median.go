package tracing

import (
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/internal/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func CheckInvariantsForMedian(
	t *testing.T,
	allTraces []Trace,
	alphaPPB uint64,
) {
	require.NoError(t, ConsecutiveReportedValuesAreDifferent(allTraces, alphaPPB))
	require.NoError(t, ReportedValuesAreObservedValues(allTraces))
	require.NoError(t, ReportedValueIsMedianOfObservations(allTraces))
}

func ConsecutiveReportedValuesAreDifferent(traces []Trace, alphaPPB uint64) error {
	// We traverse the traces and check each consecutive pair of StorePendingTransmission.
	// If they have the same (round, epoch), the medians should be identical.
	// If (round, epoch) are consecutive, the medians shouldn't be less than alphaPPB.
	var prevEpoch uint32
	var prevRound uint8
	var prevTimestamp time.Time
	var prevMedian *big.Int
	for _, trace := range traces {
		storePendingTr, isStorePendingTr := trace.(*StorePendingTransmission)
		if !isStorePendingTr {
			continue
		}
		if storePendingTr.Err != nil {
			continue
		}
		curEpoch, curRound := storePendingTr.Timestamp.Epoch, storePendingTr.Timestamp.Round
		curTimestamp := storePendingTr.Transmission.Time
		curMedian, err := evmreportcodec.ReportCodec{}.MedianFromReport(storePendingTr.Transmission.Report)
		if err != nil {
			return err
		}
		if prevRound == 0 && prevEpoch == 0 {
			prevEpoch, prevRound = curEpoch, curRound
			prevMedian = curMedian
			continue
		}
		if prevEpoch == curEpoch && prevRound == curRound {
			continue
		}
		if ((prevEpoch == curEpoch && prevRound+1 == curRound) ||
			(prevEpoch+1 == curEpoch && curRound == 1)) &&
			// Sanity check: if no more than 100ms have passed between two consecutive transmission with "similar" medians, then throw an error.
			// That's because the system generates two consecutive reports with the same value if sufficient time passes between them
			prevTimestamp.Add(100*time.Millisecond).After(curTimestamp) &&
			median.Deviates(alphaPPB, curMedian, prevMedian) {
			return fmt.Errorf("epoch=%d, round=%d, error=previous observation %s is too close to current observation %s, prevEpoch=%d, prevRound=%d, frame=%s",
				curEpoch, curRound, prevMedian, curMedian, prevEpoch, prevRound, trace)
		}
	}
	return nil
}

func ReportedValuesAreObservedValues(traces []Trace) error {
	// For each (epoch, round), collect all observed values in SendTo->Observe.
	// Check that each StorePendingTransmission in that (epoch, round) has a Median from the observations.
	var sentObservations = make(map[protocol.EpochRound][]*big.Int)
	for _, trace := range traces {
		sendTo, isSendTo := trace.(*SendTo)
		if !isSendTo {
			continue
		}
		observe, isObserve := sendTo.Message.(protocol.MessageObserve)
		if !isObserve {
			continue
		}
		key := protocol.EpochRound{observe.Epoch, observe.Round}
		if _, found := sentObservations[key]; !found {
			sentObservations[key] = []*big.Int{}
		}
		var observationProto median.NumericalMedianObservationProto
		if err := proto.Unmarshal(observe.SignedObservation.Observation, &observationProto); err != nil {
			return err
		}
		value, err := median.ToBigInt(observationProto.Value)
		if err != nil {
			return err
		}
		sentObservations[key] = append(sentObservations[key], value)
	}
	for _, trace := range traces {
		storePendingTr, isStorePendingTr := trace.(*StorePendingTransmission)
		if !isStorePendingTr {
			continue
		}
		key := protocol.EpochRound{storePendingTr.Timestamp.Epoch, storePendingTr.Timestamp.Round}
		obss, hasObservations := sentObservations[key]
		if !hasObservations {
			return fmt.Errorf("epoch=%d round=%d error=an oracle is reporting a value for a round without observations, frame='%s'",
				storePendingTr.Timestamp.Epoch, storePendingTr.Timestamp.Round, trace)
		}
		storedMedian, err := evmreportcodec.ReportCodec{}.MedianFromReport(storePendingTr.Transmission.Report)
		if err != nil {
			return err
		}
		if !contains(obss, storedMedian) {
			return fmt.Errorf("error=value sent to the contract (%s) is not part of observed values from the data sources: %v, frame='%s'",
				storedMedian, obss, trace)
		}
	}
	return nil
}

func ReportedValueIsMedianOfObservations(traces []Trace) error {
	// For each (epoch, round) we record all the Receive->Observe messages until the first Broadcast->ReportReq.
	// We also record all the StorePendingTransmission by all oracle in that round
	// and compare the Median with the median of all the Observes received in time by the leader of the round.
	receivedObservations := map[protocol.EpochRound][]*big.Int{}
	reportStarted := map[protocol.EpochRound]struct{}{}
	for _, trace := range traces {
		switch typed := trace.(type) {
		case *Receive:
			observe, isObserve := typed.Message.(protocol.MessageObserve)
			if !isObserve {
				continue
			}
			key := protocol.EpochRound{observe.Epoch, observe.Round}
			if _, started := reportStarted[key]; started {
				continue // We ignore all observations after the report started. They will not be included in the report.
			}
			if _, found := receivedObservations[key]; !found {
				receivedObservations[key] = []*big.Int{}
			}
			var observationProto median.NumericalMedianObservationProto
			if err := proto.Unmarshal(observe.SignedObservation.Observation, &observationProto); err != nil {
				return err
			}
			value, err := median.ToBigInt(observationProto.Value)
			if err != nil {
				return err
			}
			receivedObservations[key] = append(receivedObservations[key], value)
		case *Broadcast:
			reportReq, isReportReq := typed.Message.(protocol.MessageReportReq)
			if !isReportReq {
				continue
			}
			key := protocol.EpochRound{reportReq.Epoch, reportReq.Round}
			reportStarted[key] = struct{}{}
		case *StorePendingTransmission:
			transmittedMedian, err := evmreportcodec.ReportCodec{}.MedianFromReport(typed.Transmission.Report)
			if err != nil {
				return err
			}
			key := protocol.EpochRound{typed.Timestamp.Epoch, typed.Timestamp.Round}
			if _, found := receivedObservations[key]; !found {
				return fmt.Errorf("epoch=%d, round=%d, error=transmission with median=%s before observations, frame='%s'",
					typed.Timestamp.Epoch, typed.Timestamp.Round, transmittedMedian, typed)
			}
			observedMedian, ok := getMedian(receivedObservations[key])
			if !ok {
				return fmt.Errorf("epoch=%d, round=%d, error=no observations recorded for the timestamp",
					typed.Timestamp.Epoch, typed.Timestamp.Round)
			}
			if observedMedian.Cmp(transmittedMedian) != 0 {
				return fmt.Errorf("epoch=%d, round=%d, error=transmitted value %s is not the median of recorded observation '%v', frame='%s'",
					typed.Timestamp.Epoch, typed.Timestamp.Round, transmittedMedian, receivedObservations[key], typed)
			}
		}
	}
	return nil
}

// Helpers

// getMedian assumes there at least one value in the input observations
func getMedian(obss []*big.Int) (*big.Int, bool) {
	if len(obss) == 0 {
		return &big.Int{}, false
	}
	sort.Slice(obss, func(i, j int) bool {
		return obss[i].Cmp(obss[j]) < 0
	})
	return obss[len(obss)/2], true
}

func contains(xs []*big.Int, y *big.Int) bool {
	for _, x := range xs {
		if x.Cmp(y) == 0 {
			return true
		}
	}
	return false
}
