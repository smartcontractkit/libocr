// Package shim contains implementations of internal types in terms of the external types
package shim

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type OCR3_1SerializingEndpoint[RI any] struct {
	chTelemetry  chan<- *serialization.TelemetryWrapper
	configDigest types.ConfigDigest
	endpoint     types.BinaryNetworkEndpoint2
	maxSigLen    int
	logger       commontypes.Logger
	metrics      *serializingEndpointMetrics
	pluginLimits ocr3_1types.ReportingPluginLimits
	n, f         int

	closeOnce    sync.Once
	subprocesses subprocesses.Subprocesses
	chCancel     chan struct{}
	chOut        chan protocol.MessageWithSender[RI]
	taper        loghelper.LogarithmicTaper
}

var _ protocol.NetworkEndpoint[struct{}] = (*OCR3_1SerializingEndpoint[struct{}])(nil)

func (n *OCR3_1SerializingEndpoint[RI]) sendTelemetry(t *serialization.TelemetryWrapper) {
	select {
	case n.chTelemetry <- t:
		n.metrics.sentMessagesTotal.Inc()
		n.taper.Reset(func(oldCount uint64) {
			n.logger.Info("OCR3_1SerializingEndpoint: stopped dropping telemetry", commontypes.LogFields{
				"droppedCount": oldCount,
			})
		})
	default:
		n.metrics.droppedMessagesTotal.Inc()
		n.taper.Trigger(func(newCount uint64) {
			n.logger.Warn("OCR3_1SerializingEndpoint: dropping telemetry", commontypes.LogFields{
				"droppedCount": newCount,
			})
		})
	}
}

func (n *OCR3_1SerializingEndpoint[RI]) serialize(msg protocol.Message[RI]) ([]byte, *serialization.MessageWrapper) {
	if !msg.CheckSize(n.n, n.f, n.pluginLimits, n.maxSigLen) {
		n.logger.Error("OCR3_1SerializingEndpoint: Dropping outgoing message because it fails size check", commontypes.LogFields{
			"limits": n.pluginLimits,
		})
		return nil, nil
	}
	sMsg, pbm, err := serialization.Serialize(msg)
	if err != nil {
		n.logger.Error("OCR3_1SerializingEndpoint: Failed to serialize", commontypes.LogFields{
			"message": msg,
		})
		return nil, nil
	}
	return sMsg, pbm
}

func (n *OCR3_1SerializingEndpoint[RI]) deserialize(raw []byte, msgPriority types.BinaryMessageOutboundPriority, msgType protocol.MessageType) (protocol.Message[RI], *serialization.MessageWrapper, error) {
	m, pbm, err := serialization.Deserialize[RI](n.n, raw)
	if err != nil {
		return nil, nil, err
	}

	if !m.CheckSize(n.n, n.f, n.pluginLimits, n.maxSigLen) {
		return nil, nil, fmt.Errorf("message failed size check")
	}

	if !m.CheckPriority(msgPriority) {
		return nil, nil, fmt.Errorf("message failed priority check")

	}
	if !m.CheckMessageType(msgType) {
		return nil, nil, fmt.Errorf("message failed message type")
	}

	return m, pbm, nil
}

func NewOCR3_1SerializingEndpoint[RI any](
	chTelemetry chan<- *serialization.TelemetryWrapper,
	configDigest types.ConfigDigest,
	endpoint types.BinaryNetworkEndpoint2,
	maxSigLen int,
	logger commontypes.Logger,
	metricsRegisterer prometheus.Registerer,
	pluginLimits ocr3_1types.ReportingPluginLimits,
	n, f int,
) (*OCR3_1SerializingEndpoint[RI], error) {
	serializing_endpoint := &OCR3_1SerializingEndpoint[RI]{
		chTelemetry,
		configDigest,
		endpoint,
		maxSigLen,
		logger,
		newSerializingEndpointMetrics(metricsRegisterer, logger),
		pluginLimits,
		n, f,

		sync.Once{},
		subprocesses.Subprocesses{},
		make(chan struct{}),
		make(chan protocol.MessageWithSender[RI]),
		loghelper.LogarithmicTaper{},
	}
	serializing_endpoint.subprocesses.Go(serializing_endpoint.run)
	return serializing_endpoint, nil
}

func (n *OCR3_1SerializingEndpoint[RI]) run() {
	// we close chOut here rather than in Close() so as to forward the closure
	// of the underlyig BinaryNetworkEndpoint2 immediately
	defer close(n.chOut)

	chInboundMsgWithSender := n.endpoint.Receive()
	for {
		select {
		case inboundMsgWithSender, ok := <-chInboundMsgWithSender:
			if !ok {
				return
			}
			var msgToChOut protocol.Message[RI]
			switch msg := inboundMsgWithSender.InboundBinaryMessage.(type) {
			case types.InboundBinaryMessagePlain:
				m, pbm, err := n.deserialize(msg.Payload, msg.Priority, protocol.MessageTypePlain)
				if err != nil {
					n.logger.Error("OCR3_1SerializingEndpoint: Failed to deserialize", commontypes.LogFields{
						"error": err,
						"type":  protocol.MessageTypePlain,
					})
					n.sendTelemetryMessageAssertionViolation(msg.Payload, inboundMsgWithSender.Sender)
					break
				}
				msgToChOut = m
				n.sendTelemetryMessageReceived(pbm, inboundMsgWithSender.Sender)
			case types.InboundBinaryMessageRequest:
				m, pbm, err := n.deserialize(msg.Payload, msg.Priority, protocol.MessageTypeRequest)
				if err != nil {
					n.logger.Error("OCR3_1SerializingEndpoint: Failed to deserialize", commontypes.LogFields{
						"error": err,
						"type":  protocol.MessageTypeRequest,
					})
					n.sendTelemetryMessageAssertionViolation(msg.Payload, inboundMsgWithSender.Sender)
					break
				}
				if _, ok := m.(protocol.SerializableRequestMessage[RI]); !ok {
					n.logger.Error("OCR3_1SerializingEndpoint: message type assertion failed", commontypes.LogFields{})
					n.sendTelemetryMessageAssertionViolation(msg.Payload, inboundMsgWithSender.Sender)
					break
				}
				msgToChOut = m.(protocol.SerializableRequestMessage[RI]).NewInboundRequestMessage(msg.RequestHandle)
				n.sendTelemetryMessageReceived(pbm, inboundMsgWithSender.Sender)
			case types.InboundBinaryMessageResponse:
				m, pbm, err := n.deserialize(msg.Payload, msg.Priority, protocol.MessageTypeResponse)
				if err != nil {
					n.logger.Error("OCR3_1SerializingEndpoint: Failed to deserialize", commontypes.LogFields{
						"error": err,
						"type":  protocol.MessageTypeResponse,
					})
					n.sendTelemetryMessageAssertionViolation(msg.Payload, inboundMsgWithSender.Sender)
					break
				}
				if _, ok := m.(protocol.SerializableResponseMessage[RI]); !ok {
					n.logger.Error("OCR3_1SerializingEndpoint: message type assertion failed", commontypes.LogFields{})
					n.sendTelemetryMessageAssertionViolation(msg.Payload, inboundMsgWithSender.Sender)
					break
				}
				msgToChOut = m.(protocol.SerializableResponseMessage[RI]).NewInboundResponseMessage()
				n.sendTelemetryMessageReceived(pbm, inboundMsgWithSender.Sender)
			default:
				panic(fmt.Sprintf("OCR3_1SerializingEndpoint: Unexpected inbound binary message type type %T", msg))
			}
			if msgToChOut != nil {
				select {
				case n.chOut <- protocol.MessageWithSender[RI]{msgToChOut, inboundMsgWithSender.Sender}:
				case <-n.chCancel:
					return
				}
			}
		case <-n.chCancel:
			return
		}
	}
}

// Close closes the SerializingEndpoint. It does *not* close the underlying endpoint.
func (n *OCR3_1SerializingEndpoint[RI]) Close() error {
	n.closeOnce.Do(func() {
		close(n.chCancel)
		n.metrics.Close()
		n.subprocesses.Wait()
	})
	return nil
}

func (n *OCR3_1SerializingEndpoint[RI]) SendTo(msg protocol.Message[RI], to commontypes.OracleID) {
	var sMsg []byte
	var pbm *serialization.MessageWrapper
	if msg.CheckMessageType(protocol.MessageTypeRequest) {
		if _, ok := msg.(protocol.RequestMessage[RI]); !ok {
			n.logger.Error("OCR3_1SerializingEndpoint: message type assertion failed", commontypes.LogFields{})
		}
		sMsg, pbm = n.serialize(msg.(protocol.RequestMessage[RI]).GetSerializableRequestMessage())
	} else if msg.CheckMessageType(protocol.MessageTypeResponse) {
		if _, ok := msg.(protocol.ResponseMessage[RI]); !ok {
			n.logger.Error("OCR3_1SerializingEndpoint: message type assertion failed", commontypes.LogFields{})
		}
		sMsg, pbm = n.serialize(msg.(protocol.ResponseMessage[RI]).GetSerializableResponseMessage())
	} else if msg.CheckMessageType(protocol.MessageTypePlain) {
		sMsg, pbm = n.serialize(msg)
	} else {
		panic(fmt.Sprintf("OCR3_1SerializingEndpoint: Unexpected message type %T", msg))
	}
	if sMsg != nil {
		n.endpoint.SendTo(msg.GetOutboundBinaryMessage(sMsg), to)
		n.sendTelemetryMessageSent(pbm, sMsg, to)
	}
}

func (n *OCR3_1SerializingEndpoint[RI]) Broadcast(msg protocol.Message[RI]) {
	var sMsg []byte
	var pbm *serialization.MessageWrapper
	if msg.CheckMessageType(protocol.MessageTypeRequest) {
		if _, ok := msg.(protocol.RequestMessage[RI]); !ok {
			n.logger.Error("OCR3_1SerializingEndpoint: message type assertion failed", commontypes.LogFields{})
		}
		sMsg, pbm = n.serialize(msg.(protocol.RequestMessage[RI]).GetSerializableRequestMessage())
	} else if msg.CheckMessageType(protocol.MessageTypeResponse) {
		if _, ok := msg.(protocol.ResponseMessage[RI]); !ok {
			n.logger.Error("OCR3_1SerializingEndpoint: message type assertion failed", commontypes.LogFields{})
		}
		sMsg, pbm = n.serialize(msg.(protocol.ResponseMessage[RI]).GetSerializableResponseMessage())
	} else if msg.CheckMessageType(protocol.MessageTypePlain) {
		sMsg, pbm = n.serialize(msg)
	} else {
		panic(fmt.Sprintf("OCR3_1SerializingEndpoint: Unexpected message type %T", msg))
	}
	if sMsg != nil {
		n.endpoint.Broadcast(msg.GetOutboundBinaryMessage(sMsg))
		n.sendTelemetryMessageBroadcast(pbm, sMsg)
	}
}

func (n *OCR3_1SerializingEndpoint[RI]) Receive() <-chan protocol.MessageWithSender[RI] {
	return n.chOut
}

func (n *OCR3_1SerializingEndpoint[RI]) sendTelemetryMessageReceived(pbm *serialization.MessageWrapper, from commontypes.OracleID) {
	n.sendTelemetry(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_MessageReceived{&serialization.TelemetryMessageReceived{
			ConfigDigest: n.configDigest[:],
			Msg:          pbm,
			Sender:       uint32(from),
		}},
		UnixTimeNanoseconds: time.Now().UnixNano(),
	})
}

func (n *OCR3_1SerializingEndpoint[RI]) sendTelemetryMessageBroadcast(pbm *serialization.MessageWrapper, sMsg []byte) {
	n.sendTelemetry(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_MessageBroadcast{&serialization.TelemetryMessageBroadcast{
			ConfigDigest:  n.configDigest[:],
			Msg:           pbm,
			SerializedMsg: sMsg,
		}},
		UnixTimeNanoseconds: time.Now().UnixNano(),
	})
}

func (n *OCR3_1SerializingEndpoint[RI]) sendTelemetryMessageSent(pbm *serialization.MessageWrapper, sMsg []byte, to commontypes.OracleID) {
	n.sendTelemetry(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_MessageSent{&serialization.TelemetryMessageSent{
			ConfigDigest:  n.configDigest[:],
			Msg:           pbm,
			SerializedMsg: sMsg,
			Receiver:      uint32(to),
		}},
		UnixTimeNanoseconds: time.Now().UnixNano(),
	})
}

func (n *OCR3_1SerializingEndpoint[RI]) sendTelemetryMessageAssertionViolation(sMsg []byte, from commontypes.OracleID) {
	n.sendTelemetry(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_AssertionViolation{&serialization.TelemetryAssertionViolation{
			Violation: &serialization.TelemetryAssertionViolation_InvalidSerialization{&serialization.TelemetryAssertionViolationInvalidSerialization{
				ConfigDigest:  n.configDigest[:],
				SerializedMsg: sMsg,
				Sender:        uint32(from),
			}},
		}},
		UnixTimeNanoseconds: time.Now().UnixNano(),
	})
}
