// Package shim contains implementations of internal types in terms of the external types
package shim

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/managed/limits"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type OCR3_1SerializingEndpoint[RI any] struct {
	chTelemetry            chan<- *serialization.TelemetryWrapper
	configDigest           types.ConfigDigest
	endpoint               types.BinaryNetworkEndpoint2
	maxSigLen              int
	logger                 commontypes.Logger
	metrics                *serializingEndpointMetrics
	pluginLimits           ocr3_1types.ReportingPluginLimits
	publicConfig           ocr3_1config.PublicConfig
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
	publicConfig ocr3_1config.PublicConfig,
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

func (n *OCR3_1SerializingEndpoint[RI]) toOutboundBinaryMessage(msg protocol.Message[RI]) (outboundBinaryMessage types.OutboundBinaryMessage, pbMessageForTelemetry *serialization.MessageWrapper) {
	if !msg.CheckSize(n.publicConfig.N(), n.publicConfig.F, n.pluginLimits, n.maxSigLen, n.publicConfig) {
		n.logger.Error("OCR3_1SerializingEndpoint: Dropping outgoing message because it fails size check", commontypes.LogFields{
			"limits": n.pluginLimits,
		})
		return nil, nil
	}
	payload, pbMessageForTelemetry, err := serialization.Serialize(msg)
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
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbMessageForTelemetry
	case protocol.MessageEpochStartRequest[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbMessageForTelemetry
	case protocol.MessageEpochStart[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbMessageForTelemetry
	case protocol.MessageRoundStart[RI]:
		return types.OutboundBinaryMessageRequest{
			types.SingleUseSizedLimitedResponsePolicy{
				n.serializedLengthLimits.MaxLenMsgObservation,
				time.Now().Add(n.publicConfig.DeltaProgress),
			},
			payload,
			types.BinaryMessagePriorityDefault,
		}, pbMessageForTelemetry
	case protocol.MessageObservation[RI]:
		return msg.RequestHandle.MakeResponse(payload), pbMessageForTelemetry
	case protocol.MessageProposal[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbMessageForTelemetry
	case protocol.MessagePrepare[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbMessageForTelemetry
	case protocol.MessageCommit[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbMessageForTelemetry
	case protocol.MessageReportSignatures[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityDefault}, pbMessageForTelemetry
	case protocol.MessageReportsPlusPrecursorRequest[RI]:
		return types.OutboundBinaryMessageRequest{
			types.SingleUseSizedLimitedResponsePolicy{
				n.serializedLengthLimits.MaxLenMsgReportsPlusPrecursor,
				time.Now().Add(3 * n.publicConfig.GetDeltaReportsPlusPrecursorRequest()),
			},
			payload,
			types.BinaryMessagePriorityDefault,
		}, pbMessageForTelemetry
	case protocol.MessageReportsPlusPrecursor[RI]:
		return msg.RequestHandle.MakeResponse(payload), pbMessageForTelemetry
	case protocol.MessageStateSyncSummary[RI]:
		return types.OutboundBinaryMessagePlain{payload, types.BinaryMessagePriorityLow}, pbMessageForTelemetry
	case protocol.MessageBlockSyncRequest[RI]:
		return types.OutboundBinaryMessageRequest{
			types.SingleUseSizedLimitedResponsePolicy{
				n.serializedLengthLimits.MaxLenMsgBlockSyncResponse,
				msg.RequestInfo.ExpiryTimestamp,
			},
			payload,
			types.BinaryMessagePriorityLow,
		}, pbMessageForTelemetry
	case protocol.MessageBlockSyncResponse[RI]:
		return msg.RequestHandle.MakeResponse(payload), pbMessageForTelemetry
	case protocol.MessageTreeSyncChunkRequest[RI]:
		return types.OutboundBinaryMessageRequest{
			types.SingleUseSizedLimitedResponsePolicy{
				n.serializedLengthLimits.MaxLenMsgTreeSyncChunkResponse,
				msg.RequestInfo.ExpiryTimestamp,
			},
			payload,
			types.BinaryMessagePriorityLow,
		}, pbMessageForTelemetry
	case protocol.MessageTreeSyncChunkResponse[RI]:
		return msg.RequestHandle.MakeResponse(payload), pbMessageForTelemetry
	case protocol.MessageBlobOffer[RI]:
		return types.OutboundBinaryMessageRequest{
			types.SingleUseSizedLimitedResponsePolicy{
				n.serializedLengthLimits.MaxLenMsgBlobOfferResponse,
				msg.RequestInfo.ExpiryTimestamp,
			},
			payload,
			types.BinaryMessagePriorityDefault,
		}, pbMessageForTelemetry
	case protocol.MessageBlobChunkRequest[RI]:
		return types.OutboundBinaryMessageRequest{
			types.SingleUseSizedLimitedResponsePolicy{
				n.serializedLengthLimits.MaxLenMsgBlobChunkResponse,
				msg.RequestInfo.ExpiryTimestamp,
			},
			payload,
			types.BinaryMessagePriorityDefault,
		}, pbMessageForTelemetry
	case protocol.MessageBlobChunkResponse[RI]:
		return msg.RequestHandle.MakeResponse(payload), pbMessageForTelemetry
	case protocol.MessageBlobOfferResponse[RI]:
		return msg.RequestHandle.MakeResponse(payload), pbMessageForTelemetry
	}

	panic("unreachable")
}

func (n *OCR3_1SerializingEndpoint[RI]) fromInboundBinaryMessage(inboundBinaryMessage types.InboundBinaryMessage) (message protocol.Message[RI], pbMessageForTelemetry *serialization.MessageWrapper, err error) {
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

	message, pbMessageForTelemetry, err = serialization.Deserialize[RI](n.publicConfig.N(), payload, requestHandle)
	if err != nil {
		return nil, nil, err
	}

	// Check InboundBinaryMessage type and priority. We can do this here because
	// for every protocol message type we know the corresponding
	// InboundBinaryMessage type and priority.
	switch message.(type) {
	case protocol.MessageNewEpochWish[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageNewEpochWish[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageNewEpochWish")
		}
	case protocol.MessageEpochStartRequest[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageEpochStartRequest[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageEpochStartRequest")
		}
	case protocol.MessageEpochStart[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageEpochStart[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageEpochStart")
		}
	case protocol.MessageRoundStart[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageRequest); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageRoundStart[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageRoundStart")
		}
	case protocol.MessageObservation[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageResponse); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageObservation[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageObservation")
		}
	case protocol.MessageProposal[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageProposal[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageProposal")
		}
	case protocol.MessagePrepare[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessagePrepare[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessagePrepare")
		}
	case protocol.MessageCommit[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageCommit[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageCommit")
		}
	case protocol.MessageReportSignatures[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageReportSignatures[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageReportSignatures")
		}
	case protocol.MessageReportsPlusPrecursorRequest[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageRequest); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageReportsPlusPrecursorRequest[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageReportsPlusPrecursorRequest")
		}
	case protocol.MessageReportsPlusPrecursor[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageResponse); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageReportsPlusPrecursor[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageReportsPlusPrecursor")
		}
	case protocol.MessageBlockSyncRequest[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageRequest); !ok || ibm.Priority != types.BinaryMessagePriorityLow {
			return protocol.MessageBlockSyncRequest[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageBlockSyncRequest")
		}
	case protocol.MessageBlockSyncResponse[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageResponse); !ok || ibm.Priority != types.BinaryMessagePriorityLow {
			return protocol.MessageBlockSyncResponse[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageBlockSync")
		}
	case protocol.MessageStateSyncSummary[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessagePlain); !ok || ibm.Priority != types.BinaryMessagePriorityLow {
			return protocol.MessageStateSyncSummary[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageStateSyncSummary")
		}
	case protocol.MessageTreeSyncChunkRequest[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageRequest); !ok || ibm.Priority != types.BinaryMessagePriorityLow {
			return protocol.MessageTreeSyncChunkRequest[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageTreeSyncRequest")
		}
	case protocol.MessageTreeSyncChunkResponse[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageResponse); !ok || ibm.Priority != types.BinaryMessagePriorityLow {
			return protocol.MessageTreeSyncChunkResponse[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageTreeSyncChunk")
		}
	case protocol.MessageBlobOffer[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageRequest); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageBlobOffer[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageBlobOffer")
		}
	case protocol.MessageBlobOfferResponse[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageResponse); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageBlobOfferResponse[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageBlobOfferResponse")
		}
	case protocol.MessageBlobChunkRequest[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageRequest); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageBlobChunkRequest[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageBlobChunkRequest")
		}
	case protocol.MessageBlobChunkResponse[RI]:
		if ibm, ok := inboundBinaryMessage.(types.InboundBinaryMessageResponse); !ok || ibm.Priority != types.BinaryMessagePriorityDefault {
			return protocol.MessageBlobChunkResponse[RI]{}, pbMessageForTelemetry, fmt.Errorf("wrong type or priority for MessageBlobChunkResponse")
		}
	}

	if !message.CheckSize(n.publicConfig.N(), n.publicConfig.F, n.pluginLimits, n.maxSigLen, n.publicConfig) {
		return nil, nil, fmt.Errorf("message failed size check")
	}

	return message, pbMessageForTelemetry, nil
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

				priority := getBinaryMessageOutboundPriorityFromInboundBinaryMessage(raw.InboundBinaryMessage)

				message, pbMessageForTelemetry, err := n.fromInboundBinaryMessage(raw.InboundBinaryMessage)
				if err != nil {
					n.logger.Error("OCR3_1SerializingEndpoint: Failed to deserialize", commontypes.LogFields{
						"error": err,
					})
					// TODO: This will falsely report a deserialization error (without relevant details) if priority or message type
					// don't match
					n.sendTelemetry(&serialization.TelemetryWrapper{
						Wrapped: &serialization.TelemetryWrapper_AssertionViolation{&serialization.TelemetryAssertionViolation{
							Violation: &serialization.TelemetryAssertionViolation_InvalidSerialization{&serialization.TelemetryAssertionViolationInvalidSerialization{
								ConfigDigest:        n.configDigest[:],
								SerializedMsgPrefix: truncateByteSlice(raw.InboundBinaryMessage.GetPayload(), 100),
								SerializedMsgLength: uint32(len(raw.InboundBinaryMessage.GetPayload())),
								Sender:              uint32(raw.Sender),
								Priority:            uint32(priority),
							}},
						}},
						UnixTimeNanoseconds: time.Now().UnixNano(),
					})
					break
				}

				redactPbMessageForTelemetryToSaveBandwidth(pbMessageForTelemetry)
				n.sendTelemetry(&serialization.TelemetryWrapper{
					Wrapped: &serialization.TelemetryWrapper_MessageReceived{&serialization.TelemetryMessageReceived{
						ConfigDigest: n.configDigest[:],
						Msg:          pbMessageForTelemetry,
						Sender:       uint32(raw.Sender),
						Priority:     uint32(priority),
					}},
					UnixTimeNanoseconds: time.Now().UnixNano(),
				})

				select {
				case n.chOut <- protocol.MessageWithSender[RI]{message, raw.Sender}:
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
	oMsg, pbMessageForTelemetry := n.toOutboundBinaryMessage(msg)
	if oMsg != nil {
		n.endpoint.SendTo(oMsg, to)
		redactPbMessageForTelemetryToSaveBandwidth(pbMessageForTelemetry)
		n.sendTelemetry(&serialization.TelemetryWrapper{
			Wrapped: &serialization.TelemetryWrapper_MessageSent{&serialization.TelemetryMessageSent{
				ConfigDigest: n.configDigest[:],
				Msg:          pbMessageForTelemetry,
				Receiver:     uint32(to),
				Priority:     uint32(getBinaryMessageOutboundPriorityFromOutboundBinaryMessage(oMsg)),
			}},
			UnixTimeNanoseconds: time.Now().UnixNano(),
		})
	}
}

func (n *OCR3_1SerializingEndpoint[RI]) Broadcast(msg protocol.Message[RI]) {
	oMsg, pbMessageForTelemetry := n.toOutboundBinaryMessage(msg)
	if oMsg != nil {
		n.endpoint.Broadcast(oMsg)
		redactPbMessageForTelemetryToSaveBandwidth(pbMessageForTelemetry)
		n.sendTelemetry(&serialization.TelemetryWrapper{
			Wrapped: &serialization.TelemetryWrapper_MessageBroadcast{&serialization.TelemetryMessageBroadcast{
				ConfigDigest: n.configDigest[:],
				Msg:          pbMessageForTelemetry,
				Priority:     uint32(getBinaryMessageOutboundPriorityFromOutboundBinaryMessage(oMsg)),
			}},
			UnixTimeNanoseconds: time.Now().UnixNano(),
		})
	}
}

func (n *OCR3_1SerializingEndpoint[RI]) Receive() <-chan protocol.MessageWithSender[RI] {
	return n.chOut
}

func redactPbMessageForTelemetryToSaveBandwidth(pbMessageForTelemetry *serialization.MessageWrapper) {
	switch pbMessageForTelemetry.Msg.(type) {
	case *serialization.MessageWrapper_MessageReportsPlusPrecursor:
		mrpc := pbMessageForTelemetry.GetMessageReportsPlusPrecursor()
		if mrpc != nil {
			mrpc.ReportsPlusPrecursor = nil
		}
	case *serialization.MessageWrapper_MessageBlockSyncResponse:
		mbsr := pbMessageForTelemetry.GetMessageBlockSyncResponse()
		if mbsr != nil {
			mbsr.AttestedStateTransitionBlocks = nil
		}
	case *serialization.MessageWrapper_MessageTreeSyncChunkResponse:
		mtscr := pbMessageForTelemetry.GetMessageTreeSyncChunkResponse()
		if mtscr != nil {
			mtscr.KeyValues = nil
			mtscr.BoundingLeaves = nil
		}
	case *serialization.MessageWrapper_MessageBlobChunkResponse:
		mbcr := pbMessageForTelemetry.GetMessageBlobChunkResponse()
		if mbcr != nil {
			mbcr.Chunk = nil
			mbcr.Proof = nil
		}
	}
}

func getBinaryMessageOutboundPriorityFromInboundBinaryMessage(inboundBinaryMessage types.InboundBinaryMessage) types.BinaryMessageOutboundPriority {
	switch inboundBinaryMessage := inboundBinaryMessage.(type) {
	case types.InboundBinaryMessagePlain:
		return inboundBinaryMessage.Priority
	case types.InboundBinaryMessageRequest:
		return inboundBinaryMessage.Priority
	case types.InboundBinaryMessageResponse:
		return inboundBinaryMessage.Priority
	}
	panic("getBinaryMessageOutboundPriorityFromInboundBinaryMessage: unreachable")
}

func getBinaryMessageOutboundPriorityFromOutboundBinaryMessage(outboundBinaryMessage types.OutboundBinaryMessage) types.BinaryMessageOutboundPriority {
	switch outboundBinaryMessage := outboundBinaryMessage.(type) {
	case types.OutboundBinaryMessagePlain:
		return outboundBinaryMessage.Priority
	case types.OutboundBinaryMessageRequest:
		return outboundBinaryMessage.Priority
	case types.OutboundBinaryMessageResponse:
		return outboundBinaryMessage.Priority
	}
	panic("getBinaryMessageOutboundPriorityFromOutboundBinaryMessage: unreachable")
}

func truncateByteSlice(b []byte, maxLength int) []byte {
	if len(b) > maxLength {
		return b[:maxLength]
	}
	return b
}
