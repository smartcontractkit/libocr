package networking

import (
	"go.uber.org/multierr"
	"time"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

var _ commontypes.BinaryNetworkEndpoint = &ocrEndpointV1V2{}

const (
	v2InactivityTimeout     = 1 * time.Minute
	v2Headstart             = v2InactivityTimeout / 2
	lastHeardReportInterval = 5 * time.Minute
)

type ocrEndpointV1V2 struct {
	logger     commontypes.Logger
	v2peerIDs  []ragetypes.PeerID // only useful for logging
	numOracles int
	endpoints  []commontypes.BinaryNetworkEndpoint
	chRecv     chan commontypes.BinaryMessageWithSender
	processes  subprocesses.Subprocesses
	chClose    chan struct{}
}

func (o *ocrEndpointV1V2) SendTo(payload []byte, to commontypes.OracleID) {
	for _, e := range o.endpoints {
		e.SendTo(payload, to)
	}
}

func (o *ocrEndpointV1V2) Broadcast(payload []byte) {
	for _, e := range o.endpoints {
		e.Broadcast(payload)
	}
}

func (o *ocrEndpointV1V2) Receive() <-chan commontypes.BinaryMessageWithSender {
	return o.chRecv
}

func (o *ocrEndpointV1V2) mergeRecvs() {
	chRecvs := make([]<-chan commontypes.BinaryMessageWithSender, 2)
	for i, e := range o.endpoints {
		chRecvs[i] = e.Receive()
	}
	const V1, V2 = 0, 1
	lastHeardV2 := make([]time.Time, o.numOracles)
	messagesSinceLastReportV1, messagesSinceLastReportV2 := make([]int, o.numOracles), make([]int, o.numOracles)
	messagesSinceStartupV1, messagesSinceStartupV2 := make([]int, o.numOracles), make([]int, o.numOracles)
	lastMessageWasV1OrV2 := make([]string, o.numOracles)
	switchesSinceLastReport, switchesSinceStartup := make([]int, o.numOracles), make([]int, o.numOracles)
	for i := 0; i < o.numOracles; i++ {
		lastHeardV2[i] = time.Now().Add(-v2InactivityTimeout + v2Headstart)
		lastMessageWasV1OrV2[i] = "none"
	}
	ticker := time.NewTicker(lastHeardReportInterval)
	defer ticker.Stop()
	for {
		select {
		case msg := <-chRecvs[V1]:
			if time.Since(lastHeardV2[msg.Sender]) > v2InactivityTimeout {
				select {
				case o.chRecv <- msg:
				case <-o.chClose:
					return
				}
				if lastMessageWasV1OrV2[msg.Sender] != "V1" {
					switchesSinceLastReport[msg.Sender]++
					switchesSinceStartup[msg.Sender]++
				}
				lastMessageWasV1OrV2[msg.Sender] = "V1"
			}

			messagesSinceLastReportV1[msg.Sender]++
			messagesSinceStartupV1[msg.Sender]++
		case msg := <-chRecvs[V2]:
			lastHeardV2[msg.Sender] = time.Now()
			select {
			case o.chRecv <- msg:
			case <-o.chClose:
				return
			}
			if lastMessageWasV1OrV2[msg.Sender] != "V2" {
				switchesSinceLastReport[msg.Sender]++
				switchesSinceStartup[msg.Sender]++
			}
			lastMessageWasV1OrV2[msg.Sender] = "V2"

			messagesSinceLastReportV2[msg.Sender]++
			messagesSinceStartupV2[msg.Sender]++
		case <-ticker.C:
			durationSinceLastHeardV2 := make([]time.Duration, len(lastHeardV2))
			now := time.Now()
			for i, lastTime := range lastHeardV2 {
				durationSinceLastHeardV2[i] = now.Sub(lastTime)
			}
			o.logger.Info("OCR endpoint v1v2 status report", commontypes.LogFields{
				"peerIDs":                   o.v2peerIDs,
				"durationSinceLastHeardV2":  durationSinceLastHeardV2,
				"messagesSinceLastReportV2": messagesSinceLastReportV2,
				"messagesSinceStartupV2":    messagesSinceStartupV2,
				"messagesSinceLastReportV1": messagesSinceLastReportV1,
				"messagesSinceStartupV1":    messagesSinceStartupV1,
				"switchesSinceLastReport":   switchesSinceLastReport,
				"switchesSinceStartup":      switchesSinceStartup,
				"lastMessageWasV1OrV2":      lastMessageWasV1OrV2,
			})
			for i := 0; i < o.numOracles; i++ {
				messagesSinceLastReportV1[i] = 0
				messagesSinceLastReportV2[i] = 0
				switchesSinceLastReport[i] = 0
			}
		case <-o.chClose:
			return
		}
	}
}

func (o *ocrEndpointV1V2) Start() error {
	succeeded := false
	defer func() {
		if !succeeded {
			o.Close()
		}
	}()

	for _, e := range o.endpoints {
		if err := e.Start(); err != nil {
			return err
		}
	}
	o.processes.Go(o.mergeRecvs)
	succeeded = true
	return nil
}

func (o *ocrEndpointV1V2) Close() error {
	close(o.chClose)
	o.processes.Wait()
	var allErrors error
	for _, e := range o.endpoints {
		allErrors = multierr.Append(allErrors, e.Close())
	}
	return allErrors
}

func newOCREndpointV1V2(
	logger loghelper.LoggerWithContext,
	configDigest ocr2types.ConfigDigest,
	peer *concretePeer,
	v1peerIDs []p2ppeer.ID,
	v2peerIDs []ragetypes.PeerID,
	v1bootstrappers []p2ppeer.AddrInfo,
	v2bootstrappers []ragetypes.PeerInfo,
	config EndpointConfig,
	failureThreshold int,
	limits BinaryNetworkEndpointLimits,
) (*ocrEndpointV1V2, error) {
	ocrV1, err := newOCREndpointV1(
		logger,
		configDigest,
		peer,
		v1peerIDs,
		v1bootstrappers,
		config,
		failureThreshold,
		limits,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create v1 ocr endpoint")
	}
	ocrV2, err := newOCREndpointV2(
		logger,
		configDigest,
		peer,
		v2peerIDs,
		v2bootstrappers,
		EndpointConfigV2{
			config.IncomingMessageBufferSize,
			config.OutgoingMessageBufferSize,
		},
		failureThreshold,
		limits,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create v2 ocr endpoint")
	}
	return &ocrEndpointV1V2{
		logger,
		v2peerIDs,
		len(v2peerIDs),
		[]commontypes.BinaryNetworkEndpoint{ocrV1, ocrV2},
		make(chan commontypes.BinaryMessageWithSender),
		subprocesses.Subprocesses{},
		make(chan struct{}),
	}, nil
}
