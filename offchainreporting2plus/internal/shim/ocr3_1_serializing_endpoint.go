// Package shim contains implementations of internal types in terms of the external types
package shim

import (
	"fmt"
	"sync"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/managed/limits"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/serialization"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/subprocesses"
	"github.com/prometheus/client_golang/prometheus"
)

type OCR3_1SerializingEndpoint[RI any] struct {
	chTelemetry            chan<- *serialization.TelemetryWrapper
	configDigest           types.ConfigDigest
	endpoint               types.BinaryNetworkEndpoint2
	maxSigLen              int
	logger                 commontypes.Logger
	metrics                *serializingEndpointMetrics
	pluginLimits           ocr3_1types.ReportingPluginLimits
	publicConfig           ocr3config.PublicConfig
	serializedLengthLimits limits.OCR3_1SerializedLengthLimits

	mutex        sync.Mutex
	subprocesses subprocesses.Subprocesses
	started      bool
	closed       bool
	closedChOut  bool
	chCancel     chan struct{}
	chOut        chan protocol.MessageWithSender[RI]
	taper        loghelper.LogarithmicTaper
}

var _ protocol.NetworkEndpoint[struct{}] = (*OCR3_1SerializingEndpoint[struct{}])(nil)

func NewOCR3_1SerializingEndpoint[RI any](
	chTelemetry chan<- *serialization.TelemetryWrapper,
	configDigest types.ConfigDigest,
	endpoint types.BinaryNetworkEndpoint2,
	maxSigLen int,
	logger commontypes.Logger,
	metricsRegisterer prometheus.Registerer,
	pluginLimits ocr3_1types.ReportingPluginLimits,
	publicConfig ocr3config.PublicConfig,
	serializedLengthLimits limits.OCR3_1SerializedLengthLimits,
) *OCR3_1SerializingEndpoint[RI] {
	return &OCR3_1SerializingEndpoint[RI]{
		chTelemetry,
		configDigest,
		endpoint,
		maxSigLen,
		logger,
		newSerializingEndpointMetrics(metricsRegisterer, logger),
		pluginLimits,
		publicConfig,
		serializedLengthLimits,

		sync.Mutex{},
		subprocesses.Subprocesses{},
		false,
		false,
		false,
		make(chan struct{}),
		make(chan protocol.MessageWithSender[RI]),
		loghelper.LogarithmicTaper{},
	}
}

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

func (n *OCR3_1SerializingEndpoint[RI]) toOutboundBinaryMessage(msg protocol.Message[RI]) (types.OutboundBinaryMessage, *serialization.MessageWrapper) {
	if !msg.CheckSize(n.publicConfig.N(), n.publicConfig.F, n.pluginLimits, n.maxSigLen) {
		n.logger.Error("OCR3_1SerializingEndpoint: Dropping outgoing message because it fails size check", commontypes.LogFields{
			"limits": n.pluginLimits,
		})
		return nil, nil
	}
	payload, pbm, err := serialization.Serialize(msg)
	if err != nil {
		n.logger.Error("OCR3_1SerializingEndpoint: Failed to serialize", commontypes.LogFields{
			"message": msg,
		})
		return nil, nil
	}

	// Convert message into OutboundBinaryMessage. We can do this here because
	// for every protocol message type we know the corresponding
	// OutboundBinaryMessage type and priority.
	switch msg := msg.(type) {
	case protocol.MessageNewEpochWish[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbm
	case protocol.MessageEpochStartRequest[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbm
	case protocol.MessageEpochStart[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbm
	case protocol.MessageRoundStart[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbm
	case protocol.MessageObservation[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbm
	case protocol.MessageProposal[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbm
	case protocol.MessagePrepare[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbm
	case protocol.MessageCommit[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbm
	case protocol.MessageReportSignatures[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbm
	case protocol.MessageCertifiedCommitRequest[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbm
	case protocol.MessageCertifiedCommit[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbm
	case protocol.MessageStateSyncSummary[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityLow}, pbm
	case protocol.MessageBlockSyncRequest[RI]:
		return types.OutboundBinaryMessageRequest{
			types.SingleUseSizedLimitedResponsePolicy{
				n.serializedLengthLimits.MaxLenMsgBlockSyncResponse,
				time.Now().Add(protocol.DeltaMaxBlockSyncRequest),
			},
			payload,
			types.BinaryMessagePriorityLow,
		}, pbm
	case protocol.MessageBlockSyncResponse[RI]:
		return msg.RequestHandle.MakeResponse(payload), pbm
	case protocol.MessageTreeSyncChunkRequest[RI]:
		return types.OutboundBinaryMessageRequest{
			types.SingleUseSizedLimitedResponsePolicy{
				n.serializedLengthLimits.MaxLenMsgTreeSyncChunkResponse,
				time.Now().Add(protocol.DeltaMaxTreeSyncRequest),
			},
			payload,
			types.BinaryMessagePriorityLow,
		}, pbm
	case protocol.MessageTreeSyncChunkResponse[RI]:
		return msg.RequestHandle.MakeResponse(payload), pbm
	case protocol.MessageBlobOffer[RI]:
		return types.OutboundBinaryMessageRequest{
			types.SingleUseSizedLimitedResponsePolicy{
				n.serializedLengthLimits.MaxLenMsgBlobOfferResponse,
				msg.RequestInfo.ExpiryTimestamp,
			},
			payload,
			types.BinaryMessagePriorityDefault,
		}, pbm
	case protocol.MessageBlobChunkRequest[RI]:
		return types.OutboundBinaryMessageRequest{
			types.SingleUseSizedLimitedResponsePolicy{
				n.serializedLengthLimits.MaxLenMsgBlobChunkResponse,
				msg.RequestInfo.ExpiryTimestamp,
			},
			payload,
			types.BinaryMessagePriorityDefault,
		}, pbm
	case protocol.MessageBlobChunkResponse[RI]:
		return msg.RequestHandle.MakeResponse(payload), pbm
	case protocol.MessageBlobOfferResponse[RI]:
		return msg.RequestHandle.MakeResponse(payload), pbm
	}

	panic("unreachable")
}

func (n *OCR3_1SerializingEndpoint[RI]) fromInboundBinaryMessage(inboundBinaryMessage types.InboundBinaryMessage) (protocol.Message[RI], *serialization.MessageWrapper, error) {
	var payload []byte
	var requestHandle types.RequestHandle
	switch m := inboundBinaryMessage.(type) {
	case types.InboundBinaryMessagePlain:
		payload = m.Payload
	case types.InboundBinaryMessageRequest:
		payload = m.Payload
		requestHandle = m.RequestHandle
	case types.InboundBinaryMessageResponse:
		payload = m.Payload
	}

	m, pbm, err := serialization.Deserialize[RI](n.publicConfig.N(), payload, requestHandle)
	if err != nil {
		return nil, nil, err
	}

	// Check InboundBinaryMessage type and priority. We can do this here because
	// for every protocol message type we know the corresponding
	// InboundBinaryMessage type and priority.
	switch m.(type) {
	case protocol.MessageNewEpochWish[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageNewEpochWish[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageNewEpochWish")
		}
	case protocol.MessageEpochStartRequest[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageEpochStartRequest[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageEpochStartRequest")
		}
	case protocol.MessageEpochStart[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageEpochStart[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageEpochStart")
		}
	case protocol.MessageRoundStart[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageRoundStart[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageRoundStart")
		}
	case protocol.MessageObservation[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageObservation[RI]{}, pbm, fmt.Errorf("wrong type or request ID for MessageObservation")
		}
	case protocol.MessageProposal[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageProposal[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageProposal")
		}
	case protocol.MessagePrepare[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessagePrepare[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessagePrepare")
		}
	case protocol.MessageCommit[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageCommit[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageCommit")
		}
	case protocol.MessageReportSignatures[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageReportSignatures[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageReportSignatures")
		}
	case protocol.MessageCertifiedCommitRequest[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageCertifiedCommitRequest[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageCertifiedCommitRequest")
		}
	case protocol.MessageCertifiedCommit[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageCertifiedCommit[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageCertifiedCommit")
		}
	case protocol.MessageBlockSyncRequest[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageRequest); !ok || ibm.Priority != types.BinaryMessagePriorityLow {
			return protocol.MessageBlockSyncRequest[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageBlockSyncRequest")
		}
	case protocol.MessageBlockSyncResponse[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageResponse); !ok || ibm.Priority != types.BinaryMessagePriorityLow {
			return protocol.MessageBlockSyncResponse[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageBlockSync")
		}
	case protocol.MessageStateSyncSummary[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityLow {
			return protocol.MessageStateSyncSummary[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageStateSyncSummary")
		}
	case protocol.MessageTreeSyncChunkRequest[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageRequest); !ok || ibm.Priority != types.BinaryMessagePriorityLow {
			return protocol.MessageTreeSyncChunkRequest[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageTreeSyncRequest")
		}
	case protocol.MessageTreeSyncChunkResponse[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageResponse); !ok || ibm.Priority != types.BinaryMessagePriorityLow {
			return protocol.MessageTreeSyncChunkResponse[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageTreeSyncChunk")
		}
	case protocol.MessageBlobOffer[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageRequest); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageBlobOffer[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageBlobOffer")
		}
	case protocol.MessageBlobOfferResponse[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageResponse); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageBlobOfferResponse[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageBlobOfferResponse")
		}
	case protocol.MessageBlobChunkRequest[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageRequest); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageBlobChunkRequest[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageBlobChunkRequest")
		}
	case protocol.MessageBlobChunkResponse[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageResponse); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageBlobChunkResponse[RI]{}, pbm, fmt.Errorf("wrong type or priority for MessageBlobChunkResponse")
		}
	}

	if !m.CheckSize(n.publicConfig.N(), n.publicConfig.F, n.pluginLimits, n.maxSigLen) {
		return nil, nil, fmt.Errorf("message failed size check")
	}

	return m, pbm, nil
}

// Start starts the SerializingEndpoint. It will also start the underlying endpoint.
func (n *OCR3_1SerializingEndpoint[RI]) Start() error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.started {
		return fmt.Errorf("cannot start already started SerializingEndpoint")
	}
	n.started = true

	// irrelevant detail: Start() is not needed for BinaryNetworkEndpoint2
	// if err := n.endpoint.Start(); err != nil {
	// 	return fmt.Errorf("error while starting OCR3_1SerializingEndpoint: %w", err)
	// }

	n.subprocesses.Go(func() {
		chRaw := n.endpoint.Receive()
		for {
			select {
			case raw, ok := <-chRaw:
				if !ok {
					n.mutex.Lock()
					defer n.mutex.Unlock()
					n.closedChOut = true
					close(n.chOut)
					return
				}

				m, pbm, err := n.fromInboundBinaryMessage(raw.InboundBinaryMessage)
				if err != nil {
					n.logger.Error("OCR3_1SerializingEndpoint: Failed to deserialize", commontypes.LogFields{
						"error": err,
					})
					// TODO: This will falsely report a deserialization error (without relevant details) if priority or message type
					// don't match
					n.sendTelemetry(&serialization.TelemetryWrapper{
						Wrapped: &serialization.TelemetryWrapper_AssertionViolation{&serialization.TelemetryAssertionViolation{
							Violation: &serialization.TelemetryAssertionViolation_InvalidSerialization{&serialization.TelemetryAssertionViolationInvalidSerialization{
								ConfigDigest:  n.configDigest[:],
								SerializedMsg: raw.InboundBinaryMessage.GetPayload(),
								Sender:        uint32(raw.Sender),
							}},
						}},
						UnixTimeNanoseconds: time.Now().UnixNano(),
					})
					break
				}

				n.sendTelemetry(&serialization.TelemetryWrapper{
					Wrapped: &serialization.TelemetryWrapper_MessageReceived{&serialization.TelemetryMessageReceived{
						ConfigDigest: n.configDigest[:],
						Msg:          pbm,
						Sender:       uint32(raw.Sender),
					}},
					UnixTimeNanoseconds: time.Now().UnixNano(),
				})

				select {
				case n.chOut <- protocol.MessageWithSender[RI]{m, raw.Sender}:
				case <-n.chCancel:
					return
				}
			case <-n.chCancel:
				return
			}
		}
	})

	return nil
}

// Close closes the SerializingEndpoint. It will also close the underlying endpoint.
func (n *OCR3_1SerializingEndpoint[RI]) Close() error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.started && !n.closed {
		n.closed = true
		close(n.chCancel)
		n.subprocesses.Wait()

		if !n.closedChOut {
			n.closedChOut = true
			close(n.chOut)
		}

		err := n.endpoint.Close()
		n.metrics.Close()

		return err
	}

	return nil
}

func (n *OCR3_1SerializingEndpoint[RI]) SendTo(msg protocol.Message[RI], to commontypes.OracleID) {
	oMsg, pbm := n.toOutboundBinaryMessage(msg)
	if oMsg != nil {
		n.endpoint.SendTo(oMsg, to)
		n.sendTelemetry(&serialization.TelemetryWrapper{
			Wrapped: &serialization.TelemetryWrapper_MessageSent{&serialization.TelemetryMessageSent{
				ConfigDigest:  n.configDigest[:],
				Msg:           pbm,
				SerializedMsg: oMsg.GetPayload(),
				// TODO: What about priority or message type?
				Receiver: uint32(to),
			}},
			UnixTimeNanoseconds: time.Now().UnixNano(),
		})
	}
}

func (n *OCR3_1SerializingEndpoint[RI]) Broadcast(msg protocol.Message[RI]) {
	oMsg, pbm := n.toOutboundBinaryMessage(msg)
	if oMsg != nil {
		n.endpoint.Broadcast(oMsg)
		n.sendTelemetry(&serialization.TelemetryWrapper{
			Wrapped: &serialization.TelemetryWrapper_MessageBroadcast{&serialization.TelemetryMessageBroadcast{
				ConfigDigest: n.configDigest[:],
				Msg:          pbm,
				// TODO: What about priority or message type?
				SerializedMsg: oMsg.GetPayload(),
			}},
			UnixTimeNanoseconds: time.Now().UnixNano(),
		})
	}
}

func (n *OCR3_1SerializingEndpoint[RI]) Receive() <-chan protocol.MessageWithSender[RI] {
	return n.chOut
}
