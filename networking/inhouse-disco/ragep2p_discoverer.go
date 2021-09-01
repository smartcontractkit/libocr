package inhousedisco

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/libocr/internal/loghelper"
	nettypes "github.com/smartcontractkit/libocr/networking/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type Ragep2pDiscoverer struct {
	logger         loghelper.LoggerWithContext
	proc           subprocesses.Subprocesses
	ctx            context.Context
	ctxCancel      context.CancelFunc
	deltaReconcile time.Duration
	db             nettypes.DiscovererDatabase
	host           *ragep2p.Host
	proto          *discoveryProtocol

	streamsMu sync.Mutex
	streams   map[ragetypes.PeerID]*ragep2p.Stream

	chIncomingMessages chan incomingMessage
	chOutgoingMessages chan outgoingMessage
	chConnectivity     chan connectivityMsg
}

func NewRagep2pDiscoverer(deltaReconcile time.Duration, db nettypes.DiscovererDatabase) *Ragep2pDiscoverer {
	ctx, ctxCancel := context.WithCancel(context.Background())
	return &Ragep2pDiscoverer{
		nil, // filled on Start()
		subprocesses.Subprocesses{},
		ctx,
		ctxCancel,
		deltaReconcile,
		db,
		nil, // filled on Start()
		nil, // filled on Start()
		sync.Mutex{},
		make(map[ragetypes.PeerID]*ragep2p.Stream),
		make(chan incomingMessage),
		make(chan outgoingMessage),
		make(chan connectivityMsg),
	}
}

func (r *Ragep2pDiscoverer) Start(h *ragep2p.Host, privKey ed25519.PrivateKey, ownAddrs []ragetypes.Address, logger loghelper.LoggerWithContext) error {
	r.host = h
	r.logger = loghelper.MakeRootLoggerWithContext(logger)
	proto, err := newDiscoveryProtocol(
		r.deltaReconcile,
		r.chIncomingMessages,
		r.chOutgoingMessages,
		r.chConnectivity,
		privKey,
		ownAddrs,
		r.db,
		logger,
	)
	if err != nil {
		return errors.Wrap(err, "failed to construct underlying discovery protocol")
	}
	r.proto = proto
	err = r.proto.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start underlying discovery protocol")
	}
	r.proc.Go(r.connectivityLoop)
	r.proc.Go(r.writeLoop)
	return nil
}

func (r *Ragep2pDiscoverer) connectivityLoop() {
	var subs subprocesses.Subprocesses
	defer subs.Wait()
	for {
		select {
		case c := <-r.chConnectivity:
			logger := r.logger.MakeChild(commontypes.LogFields{
				"remotePeerID": c.peerID,
			})
			if c.peerID == r.host.ID() {
				break
			}
			r.streamsMu.Lock()
			if c.msgType == connectivityAdd {
				if _, exists := r.streams[c.peerID]; exists {
					r.streamsMu.Unlock()
					break
				}
				// no point in keeping very large buffers, since only
				// the latest messages matter anyways.
				bufferSize := 2
				messagesLimit := ragep2p.TokenBucketParams{
					// we expect one message every deltaReconcile seconds, let's double it
					// for good measure
					2 / r.deltaReconcile.Seconds(),
					// twice the buffer size should be plenty
					2 * uint32(bufferSize),
				}
				// bytesLimit is messagesLimit * maxMessageLength
				bytesLimit := ragep2p.TokenBucketParams{
					messagesLimit.Rate * maxMessageLength,
					messagesLimit.Capacity * maxMessageLength,
				}
				s, err := r.host.NewStream(
					c.peerID,
					"peer-discovery",
					bufferSize,
					bufferSize,
					maxMessageLength,
					messagesLimit,
					bytesLimit,
				)
				if err != nil {
					logger.Warn("NewStream failed!", reason(err))
					r.streamsMu.Unlock()
					break
				}
				r.streams[c.peerID] = s
				r.streamsMu.Unlock()
				pid := c.peerID
				subs.Go(func() {
					chDone := r.ctx.Done()
					for {
						select {
						case m, ok := <-s.ReceiveMessages():
							if !ok { // stream Close() will signal us when it's time to go
								return
							}
							w, err := fromProtoWrappedBytes(m)
							if err != nil {
								logger.Warn("Failed to unwrap incoming message", reason(err))
								break
							}
							select {
							case r.chIncomingMessages <- incomingMessage{w, pid}:
							case <-chDone:
								return
							}
						case <-chDone:
							return
						}
					}
				})
			} else {
				if _, exists := r.streams[c.peerID]; !exists {
					logger.Warn("Asked to remove connectivity with peer we don't have a stream for", nil)
					r.streamsMu.Unlock()
					break
				}
				if err := r.streams[c.peerID].Close(); err != nil {
					logger.Warn("Failed to close stream", reason(err))
				}
				delete(r.streams, c.peerID)
				r.streamsMu.Unlock()
			}
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *Ragep2pDiscoverer) writeLoop() {
	for {
		select {
		case m := <-r.chOutgoingMessages:
			r.streamsMu.Lock()
			s, exists := r.streams[m.to]
			if !exists {
				r.logger.Warn("write message to peer we don't have a stream open for", commontypes.LogFields{
					"remotePeerID": m.to,
				})
				r.streamsMu.Unlock()
				break
			}
			r.streamsMu.Unlock()
			bs, err := toBytesWrapped(m.payload)
			if err != nil {
				r.logger.Warn("failed to convert message to bytes", commontypes.LogFields{"message": m.payload})
				break
			}
			s.SendMessage(bs)
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *Ragep2pDiscoverer) Close() {
	r.proto.Close()
	r.ctxCancel()
	r.proc.Wait()
}

func (r *Ragep2pDiscoverer) AddGroup(digest types.ConfigDigest, onodes []ragetypes.PeerID, bnodes []ragetypes.PeerInfo) error {
	r.logger.Trace("Ragep2pDiscoverer::AddGroup()", commontypes.LogFields{
		"digest":     hex.EncodeToString(digest[:]),
		"oracles":    onodes,
		"bootstraps": bnodes,
	})
	return r.proto.addGroup(digest, onodes, bnodes)
}

func (r *Ragep2pDiscoverer) RemoveGroup(digest types.ConfigDigest) error {
	r.logger.Trace("Ragep2pDiscoverer::RemoveGroup()", commontypes.LogFields{
		"digest": hex.EncodeToString(digest[:]),
	})
	return r.proto.removeGroup(digest)
}

func (r *Ragep2pDiscoverer) FindPeer(peer ragetypes.PeerID) ([]ragetypes.Address, error) {
	return r.proto.FindPeer(peer)
}

var _ ragep2p.Discoverer = &Ragep2pDiscoverer{}
