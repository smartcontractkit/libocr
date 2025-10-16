package ragep2pnew

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/networking/ragep2pwrapper"
	"github.com/RoSpaceDev/libocr/ragep2p/internal/knock"
	"github.com/RoSpaceDev/libocr/ragep2p/internal/mtls"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/demuxer"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/frame"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/muxer"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/overheadawareconn"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/ratelimit"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/ratelimitaggregator"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/stream2types"

	"github.com/RoSpaceDev/libocr/ragep2p/types"
	"github.com/RoSpaceDev/libocr/subprocesses"

	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/internaltypes"
)

// Maximum number of streams with another peer that can be opened on a host
// Deprecated: use types.MaxStreamsPerPeer instead
const MaxStreamsPerPeer = types.MaxStreamsPerPeer

// Maximum stream name length
// Deprecated: use types.MaxStreamNameLength instead
const MaxStreamNameLength = types.MaxStreamNameLength

// Maximum length of messages sent with ragep2p
// Deprecated: use types.MaxMessageLength instead
const MaxMessageLength = types.MaxMessageLength

// The 5 second value is cribbed from go standard library's tls package as of version 1.16 and later
// https://cs.opensource.google/go/go/+/master:src/crypto/tls/conn.go;drc=059a9eedf45f4909db6a24242c106be15fb27193;l=1454
const netTimeout = 5 * time.Second

type hostState uint8

const (
	_ hostState = iota
	hostStatePending
	hostStateOpen
	hostStateClosed
)

type peerStreamOpenRequest struct {
	streamID       internaltypes.StreamID
	streamName     string
	streamPriority stream2types.StreamPriority
	limits         stream2types.ValidatedStream2Limits
}

type peerStreamOpenResponse struct {
	err error
}

type peerStreamUpdateLimitsRequest struct {
	streamID internaltypes.StreamID
	limits   stream2types.ValidatedStream2Limits
}

type peerStreamUpdateLimitsResponse struct {
	err error
}

type peerStreamCloseRequest struct {
	streamID internaltypes.StreamID
}

type peerStreamCloseResponse struct {
	peerHasNoStreams bool
	err              error
}

type newConnNotification struct {
	chConnTerminated <-chan struct{}
}

type streamStateNotification struct {
	streamID   internaltypes.StreamID
	streamName string // Used for sanity check, populated only on stream open and empty on stream close
	open       bool
}

type peerConnLifeCycle struct {
	connCancel       context.CancelFunc
	connSubs         subprocesses.Subprocesses
	chConnTerminated <-chan struct{}
}

type peer struct {
	chDone <-chan struct{}

	other  types.PeerID
	logger loghelper.LoggerWithContext

	metrics *peerMetrics

	incomingConnsLimiterMu sync.Mutex
	incomingConnsLimiter   *ratelimit.TokenBucket

	rateLimitAggregator *ratelimitaggregator.Aggregator

	connLifeCycleMu sync.Mutex
	connLifeCycle   peerConnLifeCycle

	mux   *muxer.Muxer
	demux *demuxer.Demuxer

	chNewConnNotification chan<- newConnNotification

	chOtherStreamStateNotification chan<- streamStateNotification
	chSelfStreamStateNotification  <-chan streamStateNotification

	chStreamOpenRequest  chan<- peerStreamOpenRequest
	chStreamOpenResponse <-chan peerStreamOpenResponse

	chStreamUpdateLimitsRequest  chan<- peerStreamUpdateLimitsRequest
	chStreamUpdateLimitsResponse <-chan peerStreamUpdateLimitsResponse

	chStreamCloseRequest  chan<- peerStreamCloseRequest
	chStreamCloseResponse <-chan peerStreamCloseResponse
}

type HostConfig struct {
	// DurationBetweenDials is the minimum duration between two dials. It is
	// not the exact duration because of jitter.
	DurationBetweenDials time.Duration
}

// A Host allows users to establish Streams with other peers identified by their
// PeerID. The host will transparently handle peer discovery, secure connection
// (re)establishment, multiplexing streams over the connection and rate
// limiting.
type Host struct {
	// Constructor args
	config            HostConfig
	keyring           types.PeerKeyring
	listenAddresses   []string
	discoverer        Discoverer
	logger            loghelper.LoggerWithContext
	metricsRegisterer prometheus.Registerer

	hostMetrics *hostMetrics

	// Derived from keyring
	id      types.PeerID
	tlsCert tls.Certificate

	// Host state
	stateMu sync.Mutex
	state   hostState

	// Manage various subprocesses of host
	subprocesses subprocesses.Subprocesses
	ctx          context.Context
	cancel       context.CancelFunc

	// Peers
	peersMu sync.Mutex
	peers   map[types.PeerID]*peer
}

// NewHost creates a new Host with the provided config, Ed25519 secret key,
// network listen address. A Discoverer is also provided to NewHost for
// discovering addresses of peers.
func NewHost(
	config HostConfig,
	keyring types.PeerKeyring,
	listenAddresses []string,
	discoverer Discoverer,
	logger commontypes.Logger,
	metricsRegisterer prometheus.Registerer,
) (*Host, error) {
	if len(listenAddresses) == 0 {
		return nil, fmt.Errorf("no listen addresses provided")
	}

	id := types.PeerIDFromKeyring(keyring)

	tlsCert, err := mtls.NewMinimalX509CertFromKeyring(keyring)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate from keyring for host: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Host{
		config,
		keyring,
		listenAddresses,
		discoverer,
		// peerID might already be set to the same value if we are managed, but we don't take any chances
		loghelper.MakeRootLoggerWithContext(logger).MakeChild(commontypes.LogFields{"id": "ragep2p", "peerID": id}),
		metricsRegisterer,

		newHostMetrics(metricsRegisterer, logger, id),

		id,
		tlsCert,

		sync.Mutex{},
		hostStatePending,

		subprocesses.Subprocesses{},
		ctx,
		cancel,

		sync.Mutex{},
		map[types.PeerID]*peer{},
	}, nil
}

// Start listening on the network interfaces and dialling peers.
func (ho *Host) Start() error {
	succeeded := false
	defer func() {
		if !succeeded {
			ho.Close()
		}
	}()
	ho.logger.Trace("ragep2p Start()", commontypes.LogFields{"listenAddresses": ho.listenAddresses})
	ho.stateMu.Lock()
	defer ho.stateMu.Unlock()

	if ho.state != hostStatePending {
		return fmt.Errorf("cannot Start() host that has already been started")
	}
	ho.state = hostStateOpen

	ho.subprocesses.Go(func() {
		ho.dialLoop()
	})
	for _, addr := range ho.listenAddresses {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("net.Listen(%q) failed: %w", addr, err)
		}
		ho.subprocesses.Go(func() {
			ho.listenLoop(ln)
		})
	}

	err := ho.discoverer.Start(Wrapped(ho), ho.keyring, ho.logger)
	if err != nil {
		return fmt.Errorf("failed to start discoverer: %w", err)
	}

	succeeded = true
	return nil
}

func remotePeerIDField(other types.PeerID) commontypes.LogFields {
	return commontypes.LogFields{
		"remotePeerID": other,
	}
}

// Caller should hold peersMu.
func (ho *Host) findOrCreatePeer(other types.PeerID) *peer {
	if _, ok := ho.peers[other]; !ok {
		logger := ho.logger.MakeChild(remotePeerIDField(other))

		metrics := newPeerMetrics(ho.metricsRegisterer, logger, ho.id, other)

		chDone := make(chan struct{})

		chConnTerminated := make(chan struct{})
		// close so that we re-dial and establish a connection
		close(chConnTerminated)

		mux := muxer.NewMuxer(logger)
		demux := demuxer.NewDemuxer()

		chNewConnNotification := make(chan newConnNotification)

		chOtherStreamStateNotification := make(chan streamStateNotification)
		chSelfStreamStateNotification := make(chan streamStateNotification)

		chStreamOpenRequest := make(chan peerStreamOpenRequest)
		chStreamOpenResponse := make(chan peerStreamOpenResponse)

		chStreamUpdateLimitsRequest := make(chan peerStreamUpdateLimitsRequest)
		chStreamUpdateLimitsResponse := make(chan peerStreamUpdateLimitsResponse)

		chStreamCloseRequest := make(chan peerStreamCloseRequest)
		chStreamCloseResponse := make(chan peerStreamCloseResponse)

		incomingConnsLimiter := ratelimit.NewTokenBucket(incomingConnsRateLimit(ho.config.DurationBetweenDials), 4, true)

		rateLimitAggregator := ratelimitaggregator.NewAggregator(logger)
		metrics.SetRateLimits(rateLimitAggregator.Aggregates())

		p := peer{
			chDone,

			other,
			logger,

			metrics,

			sync.Mutex{},
			incomingConnsLimiter,

			rateLimitAggregator,

			sync.Mutex{},
			peerConnLifeCycle{
				func() {},
				subprocesses.Subprocesses{},
				chConnTerminated,
			},

			mux,
			demux,

			chNewConnNotification,

			chOtherStreamStateNotification,
			chSelfStreamStateNotification,

			chStreamOpenRequest,
			chStreamOpenResponse,

			chStreamUpdateLimitsRequest,
			chStreamUpdateLimitsResponse,

			chStreamCloseRequest,
			chStreamCloseResponse,
		}
		ho.peers[other] = &p

		ho.subprocesses.Go(func() {
			peerLoop(
				ho.ctx,
				chDone,
				rateLimitAggregator,
				chNewConnNotification,
				chOtherStreamStateNotification,
				chSelfStreamStateNotification,
				mux,
				demux,
				chStreamOpenRequest,
				chStreamOpenResponse,
				chStreamUpdateLimitsRequest,
				chStreamUpdateLimitsResponse,
				chStreamCloseRequest,
				chStreamCloseResponse,
				logger,
				metrics,
			)
		})
	}
	return ho.peers[other]
}

func peerLoop(
	ctx context.Context,
	chDone chan<- struct{},
	rateLimitAggregator *ratelimitaggregator.Aggregator,
	chNewConnNotification <-chan newConnNotification,
	chOtherStreamStateNotification <-chan streamStateNotification,
	chSelfStreamStateNotification chan<- streamStateNotification,
	mux *muxer.Muxer,
	demux *demuxer.Demuxer,
	chStreamOpenRequest <-chan peerStreamOpenRequest,
	chStreamOpenResponse chan<- peerStreamOpenResponse,
	chStreamUpdateLimitsRequest <-chan peerStreamUpdateLimitsRequest,
	chStreamUpdateLimitsResponse chan<- peerStreamUpdateLimitsResponse,
	chStreamCloseRequest <-chan peerStreamCloseRequest,
	chStreamCloseResponse chan<- peerStreamCloseResponse,
	logger loghelper.LoggerWithContext,
	metrics *peerMetrics,
) {
	defer close(chDone)
	defer logger.Info("peerLoop exiting", nil)

	defer metrics.Close()

	type stream struct {
		name                      string
		priority                  stream2types.StreamPriority
		messagesLimit, bytesLimit types.TokenBucketParams
	}
	streams := map[internaltypes.StreamID]stream{}
	otherStreams := map[internaltypes.StreamID]struct{}{}

	var chConnTerminated <-chan struct{}

	pendingSelfStreamStateNotifications := map[internaltypes.StreamID]bool{}
	var selfStreamStateNotification streamStateNotification
	var chSelfStreamStateNotificationOrNil chan<- streamStateNotification

	for {
		chSelfStreamStateNotificationOrNil = nil
		// fake loop, we only perform zero or one iteration of this
		for streamID, state := range pendingSelfStreamStateNotifications {
			chSelfStreamStateNotificationOrNil = chSelfStreamStateNotification
			selfStreamStateNotification = streamStateNotification{
				streamID,
				streams[streamID].name,
				state,
			}
			break
		}

		select {
		case chSelfStreamStateNotificationOrNil <- selfStreamStateNotification:

			delete(pendingSelfStreamStateNotifications, selfStreamStateNotification.streamID)

			// if the stream has been opened by the other end already, switch it on right away
			if _, other := otherStreams[selfStreamStateNotification.streamID]; other && selfStreamStateNotification.open {
				if !mux.EnableStream(selfStreamStateNotification.streamID) {
					logger.Error("Assumption violation. Failed to enable stream on muxer. This stream may not work correctly.", commontypes.LogFields{
						"streamStateNotification": selfStreamStateNotification,
					})
				}
			}

		case notification := <-chNewConnNotification:
			logger.Trace("New connection, creating pending notifications of all streams", nil)

			chConnTerminated = notification.chConnTerminated
			for streamID := range streams {
				pendingSelfStreamStateNotifications[streamID] = true
			}

		case <-chConnTerminated:
			chConnTerminated = nil
			logger.Trace("Connection terminated, pausing all streams", nil)

			// Clear pending notifications
			pendingSelfStreamStateNotifications = map[internaltypes.StreamID]bool{}

			// Reset streams on other side
			otherStreams = map[internaltypes.StreamID]struct{}{}

			// Pause all streams on our side
			for streamID := range streams {
				if !mux.DisableStream(streamID) {
					logger.Error("Assumption violation. Failed to disable stream on muxer. This stream may not work correctly.", commontypes.LogFields{
						"streamID": streamID,
					})
				}
			}

			logger.Trace("Connection terminated, paused all streams", nil)

		case notification := <-chOtherStreamStateNotification:
			logger.Trace("Received stream state notification", commontypes.LogFields{
				"notification": notification,
			})

			_, other := otherStreams[notification.streamID]
			if other == notification.open {
				break
			}
			if notification.open {
				otherStreams[notification.streamID] = struct{}{}
			} else {
				delete(otherStreams, notification.streamID)
			}
			if _, ok := streams[notification.streamID]; ok {
				selfStreamName := streams[notification.streamID].name
				if notification.open && selfStreamName != notification.streamName {
					logger.Warn("Name mismatch between self and other stream", commontypes.LogFields{
						"localStreamName":  selfStreamName,
						"remoteStreamName": notification.streamName,
					})
				}
				if notification.open {
					if !mux.EnableStream(notification.streamID) {
						logger.Error("Assumption violation. Failed to enable stream on muxer. This stream may not work correctly.", commontypes.LogFields{
							"streamStateNotification": notification,
						})
					}
				} else {
					if !mux.DisableStream(notification.streamID) {
						logger.Error("Assumption violation. Failed to disable stream on muxer. This stream may not work correctly.", commontypes.LogFields{
							"streamStateNotification": notification,
						})
					}
				}
			}

		case req := <-chStreamOpenRequest:
			if _, ok := streams[req.streamID]; ok {
				chStreamOpenResponse <- peerStreamOpenResponse{
					fmt.Errorf("stream already exists"),
				}
			} else if len(streams) >= MaxStreamsPerPeer {
				chStreamOpenResponse <- peerStreamOpenResponse{
					fmt.Errorf("too many streams, expected at most %d", MaxStreamsPerPeer),
				}
			} else {
				if !mux.AddStream(req.streamID, req.streamName, req.streamPriority, req.limits.MaxOutgoingBufferedMessages) {
					logger.Error("Assumption violation. Failed to add already existing stream to muxer. This stream may not work correctly.", commontypes.LogFields{
						"streamOpenRequest": req,
					})
					// let's try to fix the problem by removing and adding the stream again
					_ = mux.RemoveStream(req.streamID)
					_ = mux.AddStream(req.streamID, req.streamName, req.streamPriority, req.limits.MaxOutgoingBufferedMessages)
				}

				if !demux.AddStream(req.streamID, req.limits.MaxIncomingBufferedMessages, req.limits.MaxMessageLength, req.limits.MessagesLimit, req.limits.BytesLimit) {
					logger.Error("Assumption violation. Failed to add already existing stream to demuxer. This stream may not work correctly.", commontypes.LogFields{
						"streamOpenRequest": req,
					})
					// let's try to fix the problem by removing and adding the stream again
					demux.RemoveStream(req.streamID)
					demux.AddStream(req.streamID, req.limits.MaxIncomingBufferedMessages, req.limits.MaxMessageLength, req.limits.MessagesLimit, req.limits.BytesLimit)
				}

				rateLimitAggregator.AddStream(req.limits.MessagesLimit, req.limits.BytesLimit)
				metrics.SetRateLimits(rateLimitAggregator.Aggregates())

				streams[req.streamID] = stream{
					req.streamName,
					req.streamPriority,
					req.limits.MessagesLimit,
					req.limits.BytesLimit,
				}
				if chConnTerminated != nil {
					pendingSelfStreamStateNotifications[req.streamID] = true
				}
				chStreamOpenResponse <- peerStreamOpenResponse{
					nil,
				}
			}

		case req := <-chStreamUpdateLimitsRequest:
			if oldS, ok := streams[req.streamID]; ok {
				s := stream{
					oldS.name,
					oldS.priority,
					req.limits.MessagesLimit,
					req.limits.BytesLimit,
				}
				streams[req.streamID] = s

				if !mux.UpdateStream(req.streamID, req.limits.MaxOutgoingBufferedMessages) {
					logger.Error("Assumption violation. Failed to update stream on muxer. This stream may not work correctly.", commontypes.LogFields{
						"streamUpdateLimitsRequest": req,
					})
					// let's try to fix the problem by removing and adding the stream again
					_ = mux.RemoveStream(req.streamID)
					_ = mux.AddStream(req.streamID, s.name, s.priority, req.limits.MaxOutgoingBufferedMessages)
				}

				if !demux.UpdateStream(
					req.streamID,
					req.limits.MaxIncomingBufferedMessages,
					req.limits.MaxMessageLength,
					req.limits.MessagesLimit,
					req.limits.BytesLimit,
				) {
					logger.Error("Assumption violation. Failed to update stream on demuxer. This stream may not work correctly.", commontypes.LogFields{
						"streamUpdateLimitsRequest": req,
					})
					// let's try to fix the problem by removing and adding the stream again
					demux.RemoveStream(req.streamID)
					_ = demux.AddStream(req.streamID, req.limits.MaxIncomingBufferedMessages, req.limits.MaxMessageLength, req.limits.MessagesLimit, req.limits.BytesLimit)
				}

				rateLimitAggregator.RemoveStream(oldS.messagesLimit, oldS.bytesLimit)
				rateLimitAggregator.AddStream(s.messagesLimit, s.bytesLimit)
				metrics.SetRateLimits(rateLimitAggregator.Aggregates())

				chStreamUpdateLimitsResponse <- peerStreamUpdateLimitsResponse{nil}
			} else {
				chStreamUpdateLimitsResponse <- peerStreamUpdateLimitsResponse{
					fmt.Errorf("stream not found"),
				}
			}

		case req := <-chStreamCloseRequest:
			if s, ok := streams[req.streamID]; ok {
				if !mux.RemoveStream(req.streamID) {
					logger.Error("Assumption violation. Failed to remove stream from muxer. Proceeding anyways because what else can I do...", commontypes.LogFields{
						"streamCloseRequest": req,
					})
				}

				demux.RemoveStream(req.streamID)

				rateLimitAggregator.RemoveStream(s.messagesLimit, s.bytesLimit)
				metrics.SetRateLimits(rateLimitAggregator.Aggregates())

				delete(streams, req.streamID)
				if chConnTerminated != nil {
					pendingSelfStreamStateNotifications[req.streamID] = false
				}
				chStreamCloseResponse <- peerStreamCloseResponse{
					len(streams) == 0,
					nil,
				}

				if len(streams) == 0 {
					return
				}
			} else {
				chStreamCloseResponse <- peerStreamCloseResponse{
					false,
					fmt.Errorf("stream not found"),
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

// Close stops listening on the network interface(s) and closes all active
// streams.
func (ho *Host) Close() error {
	ho.stateMu.Lock()
	defer ho.stateMu.Unlock()

	if ho.state != hostStateOpen {
		return fmt.Errorf("cannot Close() host that isn't open")
	}
	ho.logger.Info("Host closing discoverer", nil)
	err := ho.discoverer.Close()
	ho.logger.Info("Host winding down", nil)
	ho.state = hostStateClosed
	ho.cancel()
	ho.subprocesses.Wait()
	ho.hostMetrics.Close()
	ho.logger.Info("Host exiting", nil)
	if err != nil {
		return fmt.Errorf("failed to close discoverer: %w", err)
	}
	return nil
}

func (ho *Host) ID() types.PeerID {
	return ho.id
}

func (ho *Host) dialLoop() {
	type dialState struct {
		next uint
	}
	dialStates := make(map[types.PeerID]*dialState)
	for {
		var dialProcesses subprocesses.Subprocesses
		ho.peersMu.Lock()
		peers := make([]*peer, 0, len(ho.peers))
		for pid, p := range ho.peers {
			peers = append(peers, p)
			if dialStates[pid] == nil {
				dialStates[pid] = &dialState{0}
			}
		}
		// Some peers may have been discarded, garbage collect dial states
		for pid := range dialStates {
			if ho.peers[pid] == nil {
				delete(dialStates, pid)
			}
		}
		ho.peersMu.Unlock()
		for _, p := range peers {
			ds := dialStates[p.other]
			dialProcesses.Go(func() {
				p.connLifeCycleMu.Lock()
				chConnTerminated := p.connLifeCycle.chConnTerminated
				p.connLifeCycleMu.Unlock()
				select {
				case <-chConnTerminated:
					p.logger.Debug("Dialing", nil)
				default:
					p.logger.Trace("Dial skip", nil)
					return
				}

				addresses, err := ho.discoverer.FindPeer(p.other)
				if err != nil {
					p.logger.Warn("Discoverer error", commontypes.LogFields{"error": err})
					return
				}
				if len(addresses) == 0 {
					p.logger.Warn("Discoverer found no addresses", nil)
					return
				}

				address := string(addresses[ds.next%uint(len(addresses))])

				// We used to increment this only on dial error but a connection might fail after the Dial itself has
				// succeeded (eg. this happens with self-dials where the connection is reset after the incorrect knock
				// is received). Tracking an error so far down the stack is much harder so increment every time to give
				// a fair chance to every address.
				ds.next++

				logger := p.logger.MakeChild(commontypes.LogFields{"direction": "out", "remoteAddr": address})

				dialer := net.Dialer{
					Timeout: ho.config.DurationBetweenDials,
				}
				conn, err := dialer.DialContext(ho.ctx, "tcp", address)
				if err != nil {
					logger.Warn("Dial error", commontypes.LogFields{"error": err})
					return
				}

				logger.Trace("Dial succeeded", nil)
				ho.subprocesses.Go(func() {
					ho.handleOutgoingConnection(conn, p.other, logger)
				})
			})

		}
		dialProcesses.Wait()

		select {
		//case <-time.After(5 * time.Second): // good for testing simultaneous dials, real version is on next line
		case <-time.After(ho.config.DurationBetweenDials + time.Duration(rand.Float32()*float32(ho.config.DurationBetweenDials))):
		case <-ho.ctx.Done():
			ho.logger.Trace("Host.dialLoop exiting", nil)
			return
		}
	}
}

func (ho *Host) listenLoop(ln net.Listener) {
	ho.subprocesses.Go(func() {
		<-ho.ctx.Done()
		if err := ln.Close(); err != nil {
			ho.logger.Warn("Failed to close listener", commontypes.LogFields{"error": err})
		}
	})

	for {
		conn, err := ln.Accept()
		ho.hostMetrics.inboundDialsTotal.Inc()
		if err != nil {
			ho.logger.Info("Exiting Host.listenLoop due to error while Accepting", commontypes.LogFields{"error": err})
			return
		}
		ho.subprocesses.Go(func() {
			ho.handleIncomingConnection(conn)
		})
	}
}

func (ho *Host) handleOutgoingConnection(conn net.Conn, other types.PeerID, logger loghelper.LoggerWithContext) {
	shouldClose := true
	defer func() {
		if shouldClose {
			if err := safeClose(conn); err != nil {
				logger.Warn("Failed to close outgoing connection", commontypes.LogFields{"error": err})
			}
		}
	}()

	knck, err := knock.BuildKnock(other, ho.id, ho.keyring)
	if err != nil {
		logger.Warn("Error while building knock", commontypes.LogFields{"error": err})
		return
	}
	if err := conn.SetWriteDeadline(time.Now().Add(netTimeout)); err != nil {
		logger.Warn("Closing connection, error during SetWriteDeadline", commontypes.LogFields{"error": err})
		return
	}
	if _, err := conn.Write(knck); err != nil {
		logger.Warn("Error while sending knock", commontypes.LogFields{"error": err})
		return
	}

	ho.peersMu.Lock()
	peer, ok := ho.peers[other]
	ho.peersMu.Unlock()
	if !ok {
		// peer must have been deleted in the time between the dial being
		// started and now
		return
	}

	shouldClose = false

	overheadAwareConn := overheadawareconn.NewOverheadAwareConn(
		conn,
		peer.metrics.rawconnReadBytesTotal,
		peer.metrics.rawconnWrittenBytesTotal,
	)

	tlsConfig := newTLSConfig(
		ho.tlsCert,
		mtls.VerifyCertMatchesPubKey(other),
	)
	tlsConn := tls.Client(overheadAwareConn, tlsConfig)
	ho.handleConnection(false, overheadAwareConn, tlsConn, peer, logger)
}

func (ho *Host) handleIncomingConnection(conn net.Conn) {
	remoteAddrLogFields := commontypes.LogFields{"direction": "in", "remoteAddr": conn.RemoteAddr()}
	logger := ho.logger.MakeChild(remoteAddrLogFields)
	shouldClose := true
	defer func() {
		if shouldClose {
			if err := safeClose(conn); err != nil {
				logger.Warn("Failed to close incoming connection", commontypes.LogFields{"error": err})
			}
		}
	}()

	knck := make([]byte, knock.KnockSize)
	if err := conn.SetReadDeadline(time.Now().Add(netTimeout)); err != nil {
		logger.Warn("Closing connection, error during SetReadDeadline", commontypes.LogFields{"error": err})
		return
	}
	n, err := conn.Read(knck)
	if err != nil {
		logger.Warn("Error while reading knock", commontypes.LogFields{"error": err})
		return
	}
	if n != knock.KnockSize {
		logger.Warn("Knock too short", nil)
		return
	}

	other, err := knock.VerifyKnock(ho.id, knck)
	if err != nil {
		if errors.Is(err, knock.ErrFromSelfDial) {
			logger.Info("Self-dial knock, dropping connection. Someone has likely misconfigured their announce addresses.", nil)
		} else {
			logger.Warn("Invalid knock", commontypes.LogFields{"error": err})
		}
		return
	}

	ho.peersMu.Lock()
	peer, ok := ho.peers[*other]
	ho.peersMu.Unlock()
	if !ok {
		logger.Warn("Received incoming connection from an unknown peer, closing", remotePeerIDField(*other))
		return
	}
	logger = peer.logger.MakeChild(remoteAddrLogFields) // introduce remotePeerID in our logs since we now know it
	overheadAwareConn := overheadawareconn.NewOverheadAwareConn(
		conn,
		peer.metrics.rawconnReadBytesTotal,
		peer.metrics.rawconnWrittenBytesTotal,
	)

	shouldClose = false

	tlsConfig := newTLSConfig(
		ho.tlsCert,
		mtls.VerifyCertMatchesPubKey(*other),
	)
	tlsConn := tls.Server(overheadAwareConn, tlsConfig)
	ho.handleConnection(true, overheadAwareConn, tlsConn, peer, logger)
}

func (ho *Host) handleConnection(
	incoming bool,
	overheadAwareConn *overheadawareconn.OverheadAwareConn,
	tlsConn *tls.Conn,
	peer *peer,
	logger loghelper.LoggerWithContext,
) {
	shouldClose := true
	defer func() {
		if shouldClose {
			if err := safeClose(tlsConn); err != nil {
				logger.Warn("Failed to close connection", commontypes.LogFields{"error": err})
			}
		}
	}()

	// Handshake reads and write to the connection. Set a deadline to prevent tarpitting
	if err := tlsConn.SetDeadline(time.Now().Add(netTimeout)); err != nil {
		logger.Warn("Closing connection, error during SetDeadline", commontypes.LogFields{"error": err})
		return
	}
	// Perform handshake so that we know the public key
	if err := tlsConn.Handshake(); err != nil {
		logger.Warn("Closing connection, error during Handshake", commontypes.LogFields{"error": err})
		return
	}
	// Disable deadline. Whoever uses the connection next will have to set their own timeouts.
	if err := tlsConn.SetDeadline(time.Time{}); err != nil {
		logger.Warn("Closing connection, error during SetDeadline", commontypes.LogFields{"error": err})
		return
	}

	// get public key
	pubKey, err := mtls.PubKeyFromCert(tlsConn.ConnectionState().PeerCertificates[0])
	if err != nil {
		logger.Warn("Closing connection, error getting public key", commontypes.LogFields{"error": err})
		return
	}
	if peer.other != types.PeerIDFromPeerPublicKey(pubKey) {
		logger.Warn("TLS handshake PeerID mismatch", commontypes.LogFields{
			"expected": peer.other,
			"actual":   types.PeerIDFromPeerPublicKey(pubKey),
		})
		return
	}

	if incoming {
		peer.incomingConnsLimiterMu.Lock()
		allowed := peer.incomingConnsLimiter.RemoveTokens(1)
		peer.incomingConnsLimiterMu.Unlock()
		if !allowed {
			logger.Warn("Incoming connection rate limited", nil)
			return
		}
	}

	overheadAwareConn.SetupComplete()

	logger.Info("Connection established", nil)
	peer.metrics.connEstablishedTotal.Inc()
	if incoming {
		peer.metrics.connEstablishedInboundTotal.Inc()
	}

	// the lock here ensures there is at most one active connection at any time.
	// it also prevents races on connLifeCycle.connSubs.
	peer.connLifeCycleMu.Lock()
	peer.connLifeCycle.connCancel()
	peer.connLifeCycle.connSubs.Wait()
	connCtx, connCancel := context.WithCancel(ho.ctx)
	chConnTerminated := make(chan struct{})
	peer.connLifeCycle.connCancel = connCancel
	peer.connLifeCycle.chConnTerminated = chConnTerminated
	peer.connLifeCycle.connSubs.Go(func() {
		defer connCancel()
		authenticatedConnectionLoop(
			connCtx,
			overheadAwareConn,
			tlsConn,
			peer.chOtherStreamStateNotification,
			peer.chSelfStreamStateNotification,
			peer.mux,
			peer.demux,
			chConnTerminated,
			logger,
			peer.metrics,
		)
	})
	peer.connLifeCycleMu.Unlock()

	select {
	case peer.chNewConnNotification <- newConnNotification{chConnTerminated}:
		// keep the connection
		shouldClose = false
	case <-peer.chDone:
	case <-ho.ctx.Done():
	}
}

// NewStream creates a new bidirectional stream with peer other for streamName.
// It is parameterized with a maxMessageLength, the maximum size of a message in
// bytes and two parameters for rate limiting.
//
// Deprecated: Please switch to NewStream2.
func (ho *Host) NewStream(
	other types.PeerID,
	streamName string,
	outgoingBufferSize int, // number of messages that fit in the outgoing buffer
	incomingBufferSize int, // number of messages that fit in the incoming buffer
	maxMessageLength int,
	messagesLimit types.TokenBucketParams, // rate limit for (the number of) incoming messages
	bytesLimit types.TokenBucketParams, // rate limit for (the accumulated size in bytes of) incoming messages
) (*Stream, error) {
	stream2, err := ho.NewStream2(
		other,
		streamName,
		stream2types.StreamPriorityDefault,
		Stream2Limits{
			outgoingBufferSize,
			incomingBufferSize,
			maxMessageLength,
			messagesLimit,
			bytesLimit,
		},
	)
	if err != nil {
		return nil, err
	}

	return newStreamFromStream2(stream2)
}

// NewStream2 creates a new bidirectional stream with peer other for streamName.
// It is parameterized with a maxMessageLength, the maximum size of a message in
// bytes and two parameters for rate limiting. Compared to Stream, Stream2
// introduces an additional parameter: the message priority level.
func (ho *Host) NewStream2(
	other types.PeerID,
	streamName string,
	priority stream2types.StreamPriority,
	limits Stream2Limits,
) (Stream2, error) {
	if other == ho.id {
		return nil, fmt.Errorf("stream with self is forbidden")
	}

	if len(streamName) == 0 {
		return nil, fmt.Errorf("streamName cannot be empty")
	}
	if types.MaxStreamNameLength < len(streamName) {
		return nil, fmt.Errorf("streamName '%v' is longer than maximum length %v", streamName, types.MaxStreamNameLength)
	}

	validatedLimits, err := limits.Validate()
	if err != nil {
		return nil, err
	}

	ho.peersMu.Lock()
	defer ho.peersMu.Unlock()
	p := ho.findOrCreatePeer(other)

	sid := internaltypes.MakeStreamID(ho.id, other, streamName)

	var response peerStreamOpenResponse
	select {
	// it's important that we hold peersMu here. otherwise the peer could have
	// shut down and we'd block on the send until the host is shut down
	case p.chStreamOpenRequest <- peerStreamOpenRequest{
		sid,
		streamName,
		priority,
		validatedLimits,
	}:
		response = <-p.chStreamOpenResponse
		if response.err != nil {
			return nil, response.err
		}
	case <-ho.ctx.Done():
		return nil, fmt.Errorf("host shut down")
	}

	streamLogger := loghelper.MakeRootLoggerWithContext(p.logger).MakeChild(commontypes.LogFields{
		"streamID":   sid,
		"streamName": streamName,
	})

	ctx, cancel := context.WithCancel(ho.ctx)
	s := stream2{
		sync.Mutex{},
		false,

		streamName,
		other,
		sid,

		limits.MaxOutgoingBufferedMessages,
		limits.MaxMessageLength,
		ho,

		subprocesses.Subprocesses{},
		ctx,
		cancel,
		streamLogger,
		make(chan InboundBinaryMessage, 1),

		p.mux,
		p.demux,

		p.chStreamUpdateLimitsRequest,
		p.chStreamUpdateLimitsResponse,

		p.chStreamCloseRequest,
		p.chStreamCloseResponse,
	}

	s.subprocesses.Go(func() {
		s.receiveLoop()
	})

	streamLogger.Info("NewStream2 succeeded", commontypes.LogFields{
		"maxOutgoingBufferedMessages": limits.MaxOutgoingBufferedMessages,
		"maxIncomingBufferedMessages": limits.MaxIncomingBufferedMessages,
		"maxMessageLength":            limits.MaxMessageLength,
		"messagesLimit":               limits.MessagesLimit,
		"bytesLimit":                  limits.BytesLimit,
	})

	return &s, nil
}

/////////////////////////////////////////////
// authenticated connection handling
//////////////////////////////////////////////

func authenticatedConnectionLoop(
	ctx context.Context,
	overheadAwareConn *overheadawareconn.OverheadAwareConn,
	conn net.Conn,
	chOtherStreamStateNotification chan<- streamStateNotification,
	chSelfStreamStateNotification <-chan streamStateNotification,
	mux *muxer.Muxer,
	demux *demuxer.Demuxer,
	chTerminated chan<- struct{},
	logger loghelper.LoggerWithContext,
	metrics *peerMetrics,
) {
	defer func() {
		close(chTerminated)
		logger.Info("authenticatedConnectionLoop: exited", nil)
	}()

	var subs subprocesses.Subprocesses
	defer subs.Wait()

	defer func() {
		if err := safeClose(conn); err != nil {
			logger.Warn("Failed to close connection", commontypes.LogFields{"error": err})
		}
	}()

	childCtx, childCancel := context.WithCancel(ctx)
	defer childCancel()

	chReadTerminated := make(chan struct{})
	subs.Go(func() {
		authenticatedConnectionReadLoop(
			childCtx,
			overheadAwareConn,
			conn,
			chOtherStreamStateNotification,
			demux,
			chReadTerminated,
			logger,
			metrics,
		)
	})

	chWriteTerminated := make(chan struct{})
	subs.Go(func() {
		authenticatedConnectionWriteLoop(
			childCtx,
			conn,
			chSelfStreamStateNotification,
			mux,
			demux, // added for request/response tracking
			chWriteTerminated,
			logger,
			metrics,
		)
	})

	select {
	case <-ctx.Done():
	case <-chReadTerminated:
	case <-chWriteTerminated:
	}

	logger.Info("authenticatedConnectionLoop: winding down", nil)
}

func authenticatedConnectionReadLoop(
	ctx context.Context,
	overheadAwareConn *overheadawareconn.OverheadAwareConn,
	conn net.Conn,
	chOtherStreamStateNotification chan<- streamStateNotification,
	demux *demuxer.Demuxer,
	chReadTerminated chan<- struct{},
	logger loghelper.LoggerWithContext,
	metrics *peerMetrics,
) {
	defer close(chReadTerminated)

	readInternal := func(buf []byte) bool {
		_, err := io.ReadFull(conn, buf)
		if err != nil {
			logger.Warn("Error reading from connection", commontypes.LogFields{"error": err})
			return false
		}

		metrics.connReadProcessedBytesTotal.Add(float64(len(buf)))

		if err := overheadAwareConn.AddDeliveredApplicationDataBytes(len(buf)); err != nil {
			logger.Warn("OverheadAwareConn is asking us to drop the connection", commontypes.LogFields{"error": err})
			return false
		}

		return true
	}

	skipInternal := func(n int) bool {
		r, err := io.Copy(io.Discard, io.LimitReader(conn, int64(n)))
		if err != nil || r != int64(n) {
			logger.Warn("Error reading from connection", commontypes.LogFields{"error": err})
			return false
		}

		metrics.connReadSkippedBytesTotal.Add(float64(n))

		if err := overheadAwareConn.AddDeliveredApplicationDataBytes(n); err != nil {
			logger.Warn("OverheadAwareConn is asking us to drop the connection", commontypes.LogFields{"error": err})
			return false
		}

		return true
	}

	// We taper some logs to prevent an adversary from spamming our logs
	limitsExceededTaper := loghelper.LogarithmicTaper{}
	// Note that we never reset this taper. There shouldn't be many messages
	// with unknown stream id.
	unknownStreamIDTaper := loghelper.LogarithmicTaper{}

	// We keep track of stream names for logging.
	// Note that entries in this map are not checked for truthfulness, the remote
	// could lie about the stream name.
	remoteStreamNameByID := make(map[internaltypes.StreamID]string)

	logWithHeaderInternal := func(header frame.FrameHeader) commontypes.Logger {
		return logger.MakeChild(commontypes.LogFields{
			"payloadLength":    header.GetPayloadSize(),
			"streamID":         header.GetStreamID(),
			"remoteStreamName": remoteStreamNameByID[header.GetStreamID()],
		})
	}

	logNegativeDemuxResultInternal := func(header frame.FrameHeader, demuxShouldPushResult demuxer.ShouldPushResult) {
		switch demuxShouldPushResult {
		case demuxer.ShouldPushResultMessageTooBig:
			limitsExceededTaper.Trigger(func(count uint64) {
				logWithHeaderInternal(header).Warn("authenticatedConnectionReadLoop: message too big, dropping message", commontypes.LogFields{
					"limitsExceededDroppedCount": count,
				})
			})
		case demuxer.ShouldPushResultMessagesLimitExceeded:
			limitsExceededTaper.Trigger(func(count uint64) {
				logWithHeaderInternal(header).Warn("authenticatedConnectionReadLoop: message limit exceeded, dropping message", commontypes.LogFields{
					"limitsExceededDroppedCount": count,
				})
			})
		case demuxer.ShouldPushResultBytesLimitExceeded:
			limitsExceededTaper.Trigger(func(count uint64) {
				logWithHeaderInternal(header).Warn("authenticatedConnectionReadLoop: bytes limit exceeded, dropping message", commontypes.LogFields{
					"limitsExceededDroppedCount": count,
				})
			})
		case demuxer.ShouldPushResultUnknownStream:
			unknownStreamIDTaper.Trigger(func(count uint64) {
				logWithHeaderInternal(header).Warn("authenticatedConnectionReadLoop: unknown stream id, dropping message", commontypes.LogFields{
					"unknownStreamIDDroppedCount": count,
				})
			})
		case demuxer.ShouldPushResultResponseRejected:
			limitsExceededTaper.Trigger(func(count uint64) {
				logWithHeaderInternal(header).Warn(
					"authenticatedConnectionReadLoop: response rejected, dropping message", commontypes.LogFields{
						"limitsExceededDroppedCount": count,
					},
				)
			})
		case demuxer.ShouldPushResultYes:
			logger.Critical("authenticatedConnectionReadLoop.logNegativeDemuxResultInternal should never hit shouldPushResultYes", nil)
		}
	}

	demuxPushInternal := func(header frame.FrameHeader, msg InboundBinaryMessage) {
		switch demux.PushMessage(header.GetStreamID(), msg) {
		case demuxer.PushResultDropped:
			logWithHeaderInternal(header).Trace("authenticatedConnectionReadLoop: demuxer is overflowing for stream, dropping oldest message", nil)
		case demuxer.PushResultUnknownStream:
			unknownStreamIDTaper.Trigger(func(count uint64) {
				logWithHeaderInternal(header).Warn("authenticatedConnectionReadLoop: unknown stream id, dropping message", commontypes.LogFields{
					"unknownStreamIDDroppedCount": count,
				})
			})
		case demuxer.PushResultSuccess:
		}
	}

	// We keep track of the number of open & close frames that we have received.
	openCloseFramesReceived := 0
	const maxOpenCloseFramesReceived = 2 * MaxStreamsPerPeer

	maxOpenCloseFramesExceededInternal := func(streamID internaltypes.StreamID, payloadSize int) bool {
		if openCloseFramesReceived <= maxOpenCloseFramesReceived {
			return false
		}

		childLogger := logger.MakeChild(commontypes.LogFields{
			"payloadLength":    payloadSize,
			"streamID":         streamID,
			"remoteStreamName": remoteStreamNameByID[streamID],
		})
		childLogger.Warn("authenticatedConnectionReadLoop: peer received too many open/close frames, closing connection",
			commontypes.LogFields{
				"maxOpenCloseFramesReceived": maxOpenCloseFramesReceived,
			})
		return true
	}

	frameHeaderReader := frame.MakeFrameHeaderReader(readInternal)

	for {
		header, err := frameHeaderReader.ReadFrameHeader()
		if err != nil {
			if errors.Is(err, frame.ErrReadFrameHeaderReadFailed) {
				// no need to log, readInternal will have already done so
			} else {
				logger.Warn("authenticatedConnectionReadLoop: could not read frame header, closing connection", commontypes.LogFields{
					"error": err,
				})
			}
			return
		}

		switch header := header.(type) {
		case frame.OpenStreamFrameHeader:
			openCloseFramesReceived++

			streamName := make([]byte, header.PayloadSize)
			if !readInternal(streamName) {
				return
			}
			remoteStreamNameByID[header.StreamID] = string(streamName)

			select {
			case chOtherStreamStateNotification <- streamStateNotification{header.StreamID, string(streamName), true}:
			case <-ctx.Done():
				return
			}

			if maxOpenCloseFramesExceededInternal(header.StreamID, header.PayloadSize) {
				return
			}
		case frame.CloseStreamFrameHeader:
			openCloseFramesReceived++

			delete(remoteStreamNameByID, header.StreamID)
			select {
			case chOtherStreamStateNotification <- streamStateNotification{header.StreamID, "", false}:
			case <-ctx.Done():
				return
			}

			if maxOpenCloseFramesExceededInternal(header.StreamID, 0) {
				return
			}
		case frame.MessagePlainFrameHeader:
			demuxShouldPushResult := demux.ShouldPush(header.StreamID, header.PayloadSize)
			switch demuxShouldPushResult {
			case demuxer.ShouldPushResultMessagesLimitExceeded, demuxer.ShouldPushResultBytesLimitExceeded, demuxer.ShouldPushResultUnknownStream, demuxer.ShouldPushResultResponseRejected, demuxer.ShouldPushResultMessageTooBig:
				logNegativeDemuxResultInternal(header, demuxShouldPushResult)
				if !skipInternal(header.PayloadSize) {
					return
				}
			case demuxer.ShouldPushResultYes:
				payload := make([]byte, header.PayloadSize)
				if !readInternal(payload) {
					return
				}
				demuxPushInternal(
					header,
					InboundBinaryMessagePlain{payload},
				)
			}
		case frame.MessageRequestFrameHeader:
			demuxShouldPushResult := demux.ShouldPush(header.StreamID, header.PayloadSize)
			switch demuxShouldPushResult {
			case demuxer.ShouldPushResultMessagesLimitExceeded, demuxer.ShouldPushResultBytesLimitExceeded, demuxer.ShouldPushResultUnknownStream, demuxer.ShouldPushResultResponseRejected, demuxer.ShouldPushResultMessageTooBig:
				logNegativeDemuxResultInternal(header, demuxShouldPushResult)
				if !skipInternal(header.PayloadSize) {
					return
				}
			case demuxer.ShouldPushResultYes:
				payload := make([]byte, header.PayloadSize)
				if !readInternal(payload) {
					return
				}
				demuxPushInternal(
					header,
					InboundBinaryMessageRequest{RequestHandle(header.RequestID), payload},
				)
			}
		case frame.MessageResponseFrameHeader:
			demuxShouldPushResult := demux.ShouldPushResponse(header.StreamID, header.RequestID, header.PayloadSize)
			switch demuxShouldPushResult {
			case demuxer.ShouldPushResultMessagesLimitExceeded, demuxer.ShouldPushResultBytesLimitExceeded, demuxer.ShouldPushResultUnknownStream, demuxer.ShouldPushResultResponseRejected, demuxer.ShouldPushResultMessageTooBig:
				logNegativeDemuxResultInternal(header, demuxShouldPushResult)
				if !skipInternal(header.PayloadSize) {
					return
				}
			case demuxer.ShouldPushResultYes:
				limitsExceededTaper.Reset(func(oldCount uint64) {
					logWithHeaderInternal(header).Info("authenticatedConnectionReadLoop: limits are no longer being exceeded", commontypes.LogFields{
						"droppedCount": oldCount,
					})
				})

				payload := make([]byte, header.PayloadSize)
				if !readInternal(payload) {
					return
				}
				demuxPushInternal(
					header,
					InboundBinaryMessageResponse{payload},
				)
			}
		default:
			panic("unknown type of frame.FrameHeader")
		}
	}
}

func authenticatedConnectionWriteLoop(
	ctx context.Context,
	conn net.Conn,
	chSelfStreamStateNotification <-chan streamStateNotification,
	mux *muxer.Muxer,
	demux *demuxer.Demuxer,
	chWriteTerminated chan<- struct{},
	logger loghelper.LoggerWithContext,
	metrics *peerMetrics,
) {
	writeInternal := func(buf []byte) bool {
		_, err := conn.Write(buf)
		if err != nil {
			logger.Warn("Error writing to connection", commontypes.LogFields{"error": err})
			// shut everything down
			if err := safeClose(conn); err != nil {
				logger.Warn("Failed to close connection", commontypes.LogFields{"error": err})
			}
			close(chWriteTerminated)
			return false
		}
		metrics.connWrittenBytesTotal.Add(float64(len(buf)))
		return true
	}

	// If two bufs are smallish, try to coalesce the two writes since each write
	// to the underlying TLS conn will emit a new TLS record.
	// 2048 bytes since TLS per-(application data)-record overhead for us is 22 bytes.
	// When we exceed that size we lose only ~1% to the additional record's overhead.
	writeTwoInternalBuf := make([]byte, 0, 2048)
	writeTwoInternal := func(buf1 []byte, buf2 []byte) bool {
		if len(buf1)+len(buf2) <= cap(writeTwoInternalBuf) {
			writeTwoInternalBuf = writeTwoInternalBuf[:0]
			writeTwoInternalBuf = append(writeTwoInternalBuf, buf1...)
			writeTwoInternalBuf = append(writeTwoInternalBuf, buf2...)
			return writeInternal(writeTwoInternalBuf)
		} else {
			return writeInternal(buf1) && writeInternal(buf2)
		}
	}

	sendInternal := func(streamID internaltypes.StreamID, message OutboundBinaryMessage) bool {
		if err := conn.SetWriteDeadline(time.Now().Add(netTimeout)); err != nil {
			logger.Warn("Closing connection, error during SetWriteDeadline", commontypes.LogFields{"error": err})
			return false
		}

		// The header's value is set based on the type of message in the switch statement below.
		var header []byte
		var payload []byte

		switch m := message.(type) {
		case OutboundBinaryMessageRequest:
			requestID, err := internaltypes.MakeRandomRequestID()
			if err != nil {
				logger.Error("Error while sending request (failed to generate random request id)", commontypes.LogFields{
					"error": err,
				})
				return false
			}
			demux.SetPolicy(streamID, requestID, m.ResponsePolicy)
			header = frame.MessageRequestFrameHeader{streamID, len(m.Payload), requestID}.Encode()
			payload = m.Payload

		case OutboundBinaryMessageResponse:
			header = frame.MessageResponseFrameHeader{
				streamID,
				len(m.Payload),
				stream2types.RequestIDOfOutboundBinaryMessageResponse(m),
			}.Encode()
			payload = m.Payload

		case OutboundBinaryMessagePlain:
			header = frame.MessagePlainFrameHeader{streamID, len(m.Payload)}.Encode()
			payload = m.Payload
		}

		if !writeTwoInternal(
			header,
			payload,
		) {
			return false
		}
		metrics.messageBytes.Observe(float64(len(payload)))
		return true
	}

	handleStreamStateNotificationsInternal := func(notification streamStateNotification) bool {
		if err := conn.SetWriteDeadline(time.Now().Add(netTimeout)); err != nil {
			logger.Warn("Closing connection, error during SetWriteDeadline", commontypes.LogFields{"error": err})
			return false
		}
		if notification.open {
			streamName := []byte(notification.streamName)
			if !writeTwoInternal(
				frame.OpenStreamFrameHeader{notification.streamID, len(streamName)}.Encode(),
				streamName,
			) {
				return false
			}
		} else {
			if !writeInternal(frame.CloseStreamFrameHeader{notification.streamID}.Encode()) {
				return false
			}
		}
		return true
	}

	for {
		select {
		case <-ctx.Done():
			return

		case notification := <-chSelfStreamStateNotification:
			if !handleStreamStateNotificationsInternal(notification) {
				return
			}
			continue

		case <-mux.SignalMaybePending():
		}

		for {
			select {
			case <-ctx.Done():
				return

			case notification := <-chSelfStreamStateNotification:
				if !handleStreamStateNotificationsInternal(notification) {
					return
				}

			default:
			}

			msg, sid := mux.Pop()
			if msg == nil {
				break
			}
			if !sendInternal(sid, msg) {
				return
			}
		}
	}
}

// gotta be careful about closing tls connections to make sure we don't get
// tarpitted
func safeClose(conn net.Conn) error {
	// This isn't needed in more recent versions of go, but better safe than sorry!
	errDeadline := conn.SetWriteDeadline(time.Now().Add(netTimeout))
	errClose := conn.Close()
	if errClose != nil {
		return errClose
	}
	if errDeadline != nil {
		return errDeadline
	}
	return nil
}

func incomingConnsRateLimit(durationBetweenDials time.Duration) ratelimit.MillitokensPerSecond {
	// 2 dials per DurationBetweenDials are okay
	result := ratelimit.MillitokensPerSecond(2.0 / durationBetweenDials.Seconds() * 1000.0)
	// dialing once every two seconds is always okay
	if result < 500 {
		result = 500
	}
	return result
}

// Discoverer is responsible for discovering the addresses of peers on the network.
type Discoverer interface {
	Start(host ragep2pwrapper.Host, keyring types.PeerKeyring, logger loghelper.LoggerWithContext) error
	Close() error
	FindPeer(peer types.PeerID) ([]types.Address, error)
}
