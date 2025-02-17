package shim

// type OCR31SerializingEndpoint[RI any] struct {
// 	chTelemetry  chan<- *serialization.TelemetryWrapper
// 	configDigest types.ConfigDigest
// 	endpoint     ocr2types.BinaryNetworkEndpoint2
// 	maxSigLen    int
// 	logger       commontypes.Logger
// 	metrics      *serializingEndpointMetrics
// 	pluginLimits ocr3types.ReportingPluginLimits
// 	n, f         int

// 	mutex        sync.Mutex
// 	subprocesses subprocesses.Subprocesses
// 	started      bool
// 	closed       bool
// 	closedChOut  bool
// 	chCancel     chan struct{}
// 	chOut        chan protocol.MessageWithSender[RI]
// 	taper        loghelper.LogarithmicTaper
// }

// func (n *OCR31SerializingEndpoint[RI]) SendTo(msg protocol.Message[RI], to commontypes.OracleID) {
// 	sMsg, pbm := n.serialize(msg)
// 	if sMsg != nil {
// 		n.endpoint.SendTo(sMsg, to)
// 		n.sendTelemetry(&serialization.TelemetryWrapper{
// 			Wrapped: &serialization.TelemetryWrapper_MessageSent{&serialization.TelemetryMessageSent{
// 				ConfigDigest:  n.configDigest[:],
// 				Msg:           pbm,
// 				SerializedMsg: sMsg,
// 				Receiver:      uint32(to),
// 			}},
// 			UnixTimeNanoseconds: time.Now().UnixNano(),
// 		})
// 	}
// }

// func (n *OCR31SerializingEndpoint[RI]) Receive() <-chan protocol.MessageWithSender[RI] {
// 	return n.chOut
// }

// func (n *OCR31SerializingEndpoint[RI]) TranslateLoop() {
// 	msg := <-n.endpoint.Receive()
// 	switch msg.InboundBinaryMessage.(type) {
// 	case ocr2types.InboundBinaryMessagePlain:
// 	case ocr2types.InboundBinaryMessageResponse:
// 	case ocr2types.InboundBinaryMessageRequest:
// 		// in this case, we need to attach a handle to what we received.
// 	}
// }
