package inhousedisco

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/libocr/commontypes"
	nettypes "github.com/smartcontractkit/libocr/networking/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type incomingMessage struct {
	payload WrappableMessage
	from    ragetypes.PeerID
}

type outgoingMessage struct {
	payload WrappableMessage
	to      ragetypes.PeerID
}

type discoveryProtocol struct {
	deltaReconcile     time.Duration
	chIncomingMessages <-chan incomingMessage
	chOutgoingMessages chan<- outgoingMessage
	chConnectivity     chan<- connectivityMsg
	chInternalBump     chan Announcement
	privKey            ed25519.PrivateKey
	ownID              ragetypes.PeerID
	ownAddrs           []ragetypes.Address

	mu                      sync.RWMutex
	bestAnnouncement        map[ragetypes.PeerID]Announcement
	groups                  map[types.ConfigDigest]*group
	bootstrappers           map[ragetypes.PeerID]map[ragetypes.Address]int
	numGroupsByOracle       map[ragetypes.PeerID]int
	numGroupsByBootstrapper map[ragetypes.PeerID]int

	db nettypes.DiscovererDatabase

	processes subprocesses.Subprocesses
	ctx       context.Context
	ctxCancel context.CancelFunc
	logger    loghelper.LoggerWithContext
}

const (
	announcementVersionWarnThreshold = 100e6

	saveInterval       = 2 * time.Minute
	reportInitialDelay = 10 * time.Second
	reportInterval     = 5 * time.Minute
)

func newDiscoveryProtocol(
	deltaReconcile time.Duration,
	chIncomingMessages <-chan incomingMessage,
	chOutgoingMessages chan<- outgoingMessage,
	chConnectivity chan<- connectivityMsg,
	privKey ed25519.PrivateKey,
	ownAddrs []ragetypes.Address,
	db nettypes.DiscovererDatabase,
	logger loghelper.LoggerWithContext,
) (*discoveryProtocol, error) {
	ownID, err := ragetypes.PeerIDFromPrivateKey(privKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain peer id from private key")
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	return &discoveryProtocol{
		deltaReconcile,
		chIncomingMessages,
		chOutgoingMessages,
		chConnectivity,
		make(chan Announcement),
		privKey,
		ownID,
		ownAddrs,
		sync.RWMutex{},
		make(map[ragetypes.PeerID]Announcement),
		make(map[types.ConfigDigest]*group),
		make(map[ragetypes.PeerID]map[ragetypes.Address]int),
		make(map[ragetypes.PeerID]int),
		make(map[ragetypes.PeerID]int),
		db,
		subprocesses.Subprocesses{},
		ctx,
		ctxCancel,
		logger.MakeChild(commontypes.LogFields{"struct": "discoveryProtocol"}),
	}, nil
}

func (p *discoveryProtocol) Start() error {
	_, _, err := p.unsafeBumpOwnAnnouncement()
	if err != nil {
		return errors.Wrap(err, "failed to bump own announcement")
	}
	p.processes.Go(p.recvLoop)
	p.processes.Go(p.sendLoop)
	p.processes.Go(p.saveLoop)
	p.processes.Go(p.statusReportLoop)
	return nil
}

func formatAnnouncementsForReport(allIDs map[ragetypes.PeerID]struct{}, baSigned map[ragetypes.PeerID]Announcement) (string, int) {
	// Would use json here but I want to avoid having quotes in logs as it would cause escaping all over the place.
	var sb strings.Builder
	sb.WriteRune('{')
	i := 0
	undetected := 0
	for id := range allIDs {
		ann, exists := baSigned[id]
		var s string
		if exists {
			s = ann.unsignedAnnouncement.String()
		} else {
			s = "<not found>"
			undetected++
		}

		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(id.String())
		sb.WriteString(": ")
		sb.WriteString(s)
		i++
	}
	sb.WriteRune('}')
	return sb.String(), undetected
}

func (p *discoveryProtocol) statusReportLoop() {
	chDone := p.ctx.Done()
	timer := time.After(reportInitialDelay)
	for {
		select {
		case <-timer:
			func() {
				p.mu.RLock()
				defer p.mu.RUnlock()
				uniquePeersToDetect := make(map[ragetypes.PeerID]struct{})
				for id, cnt := range p.numGroupsByOracle {
					if cnt == 0 {
						continue
					}
					uniquePeersToDetect[id] = struct{}{}
				}

				reportStr, undetected := formatAnnouncementsForReport(uniquePeersToDetect, p.bestAnnouncement)
				p.logger.Info("Discoverer status report", commontypes.LogFields{
					"statusByPeer":    reportStr,
					"peersToDetect":   len(uniquePeersToDetect),
					"peersUndetected": undetected,
					"peersDetected":   len(uniquePeersToDetect) - undetected,
				})
				timer = time.After(reportInterval)
			}()
		case <-chDone:
			return
		}
	}
}

// Peer A is allowed to learn about an Announcement by peer B if B is an oracle node in
// one of the groups A participates in.
func (p *discoveryProtocol) unsafeAllowedPeers(ann Announcement) (ps []ragetypes.PeerID) {
	annPeerID, err := ann.PeerID()
	if err != nil {
		p.logger.Warn("failed to obtain peer id from announcement", reason(err))
		return
	}
	peers := make(map[ragetypes.PeerID]struct{})
	for _, g := range p.groups {
		if !g.hasOracle(annPeerID) {
			continue
		}
		for _, pid := range g.peerIDs() {
			peers[pid] = struct{}{}
		}
	}
	for pid := range peers {
		if pid == p.ownID {
			continue
		}
		ps = append(ps, pid)
	}
	return
}

func (p *discoveryProtocol) addGroup(digest types.ConfigDigest, onodes []ragetypes.PeerID, bnodes []ragetypes.PeerInfo) error {
	var newPeerIDs []ragetypes.PeerID
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.groups[digest]; exists {
		return fmt.Errorf("asked to add group with digest we already have (digest: %s)", digest.Hex())
	}
	newGroup := group{oracleNodes: onodes, bootstrapperNodes: bnodes}
	p.groups[digest] = &newGroup
	for _, oid := range onodes {
		if p.numGroupsByOracle[oid] == 0 {
			newPeerIDs = append(newPeerIDs, oid)
		}
		p.numGroupsByOracle[oid]++
	}
	for _, bs := range bnodes {
		p.numGroupsByBootstrapper[bs.ID]++
		for _, addr := range bs.Addrs {
			if _, exists := p.bootstrappers[bs.ID]; !exists {
				p.bootstrappers[bs.ID] = make(map[ragetypes.Address]int)
			}
			p.bootstrappers[bs.ID][addr]++
		}
	}
	for _, pid := range newGroup.peerIDs() {
		// it's ok to send connectivityAdd messages multiple times
		select {
		case p.chConnectivity <- connectivityMsg{connectivityAdd, pid}:
		case <-p.ctx.Done():
			return nil
		}
	}

	// we hold mu here
	if err := p.unsafeLoadFromDB(newPeerIDs); err != nil {
		// db-level errors are not prohibitive
		p.logger.Warn("Failed to load announcements from db", reason(err))
	}
	return nil
}

func (p *discoveryProtocol) unsafeLoadFromDB(ragePeerIDs []ragetypes.PeerID) error {
	// The database may have been set to nil, and we don't necessarily need it to function.
	if len(ragePeerIDs) == 0 || p.db == nil {
		return nil
	}
	strPeerIDs := make([]string, len(ragePeerIDs))
	for i, pid := range ragePeerIDs {
		strPeerIDs[i] = pid.String()
	}
	annByID, err := p.db.ReadAnnouncements(p.ctx, strPeerIDs)
	if err != nil {
		return err
	}
	for _, dbannBytes := range annByID {
		dbann, err := deserializeSignedAnnouncement(dbannBytes)
		if err != nil {
			p.logger.Warn("failed to deserialize signed announcement from db", commontypes.LogFields{
				"error": err,
				"bytes": dbannBytes,
			})
			continue
		}
		err = p.unsafeProcessAnnouncement(dbann)
		if err != nil {
			p.logger.Warn("failed to process announcement from db", commontypes.LogFields{
				"error": err,
				"ann":   dbann,
			})
		}
	}
	return nil
}

func (p *discoveryProtocol) saveAnnouncementToDB(ann Announcement) error {
	if p.db == nil {
		return nil
	}
	ser, err := ann.serialize()
	if err != nil {
		return err
	}
	pid, err := ann.PeerID()
	if err != nil {
		return err
	}
	return p.db.StoreAnnouncement(p.ctx, pid.String(), ser)
}

func (p *discoveryProtocol) saveToDB() error {
	if p.db == nil {
		return nil
	}
	p.mu.RLock()
	defer p.mu.RUnlock()

	var allErrors error
	for _, ann := range p.bestAnnouncement {
		allErrors = multierr.Append(allErrors, p.saveAnnouncementToDB(ann))
	}
	return allErrors
}

func (p *discoveryProtocol) saveLoop() {
	if p.db == nil {
		return
	}
	for {
		select {
		case <-time.After(saveInterval):
		case <-p.ctx.Done():
			return
		}

		if err := p.saveToDB(); err != nil {
			p.logger.Warn("failed to save announcements to db", reason(err))
		}
	}
}

func (p *discoveryProtocol) removeGroup(digest types.ConfigDigest) error {
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "removeGroup"})
	logger.Trace("Called", nil)
	p.mu.Lock()
	defer p.mu.Unlock()

	goneGroup, exists := p.groups[digest]
	if !exists {
		return fmt.Errorf("can't remove group that is not registered (digest: %s)", digest.Hex())
	}

	delete(p.groups, digest)

	for _, oid := range goneGroup.oracleIDs() {
		p.numGroupsByOracle[oid]--
		if p.numGroupsByOracle[oid] == 0 {
			if ann, exists := p.bestAnnouncement[oid]; exists {
				if err := p.saveAnnouncementToDB(ann); err != nil {
					p.logger.Warn("Failed to save announcement from removed group to DB", reason(err))
				}
			}
			delete(p.bestAnnouncement, oid)
			delete(p.numGroupsByOracle, oid)
		}
	}

	for _, binfo := range goneGroup.bootstrapperNodes {
		bid := binfo.ID

		p.numGroupsByBootstrapper[bid]--
		if p.numGroupsByBootstrapper[bid] == 0 {
			delete(p.numGroupsByBootstrapper, bid)
			delete(p.bootstrappers, bid)
			continue
		}
		for _, addr := range binfo.Addrs {
			p.bootstrappers[bid][addr]--
			if p.bootstrappers[bid][addr] == 0 {
				delete(p.bootstrappers[bid], addr)
			}
		}
	}

	// Cleanup connections for peers we don't have in any group anymore.
	for _, pid := range goneGroup.peerIDs() {
		if p.numGroupsByOracle[pid]+p.numGroupsByBootstrapper[pid] == 0 {
			select {
			case p.chConnectivity <- connectivityMsg{connectivityRemove, pid}:
			case <-p.ctx.Done():
				return nil
			}
		}
	}

	return nil
}

func (p *discoveryProtocol) FindPeer(peer ragetypes.PeerID) (addrs []ragetypes.Address, err error) {
	allAddrs := make(map[ragetypes.Address]struct{})
	p.mu.RLock()
	defer p.mu.RUnlock()
	if baddrs, ok := p.bootstrappers[peer]; ok {
		for a := range baddrs {
			allAddrs[a] = struct{}{}
		}
	}
	if ann, ok := p.bestAnnouncement[peer]; ok {
		for _, a := range ann.Addrs {
			allAddrs[a] = struct{}{}
		}
	}
	for a := range allAddrs {
		addrs = append(addrs, a)
	}
	return
}

func (p *discoveryProtocol) recvLoop() {
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "recvLoop"})
	logger.Debug("Entering", nil)
	defer logger.Debug("Exiting", nil)
	for {
		select {
		case <-p.ctx.Done():
			return
		case msg := <-p.chIncomingMessages:
			logger := logger.MakeChild(commontypes.LogFields{"remotePeerID": msg.from})
			switch v := msg.payload.(type) {
			case *Announcement:
				logger.Trace("Received announcement", v.toLogFields())
				if err := p.processAnnouncement(*v); err != nil {
					logger := logger.MakeChild(reason(err))
					logger = logger.MakeChild(v.toLogFields())
					logger.Warn("Failed to process announcement", nil)
				}
			case *reconcile:
				logger.Trace("Received reconcile", commontypes.LogFields{"reconcile": v.toLogFields()})
				for _, ann := range v.Anns {
					if err := p.processAnnouncement(ann); err != nil {
						logger := logger.MakeChild(reason(err))
						logger = logger.MakeChild(v.toLogFields())
						logger = logger.MakeChild(ann.toLogFields())
						logger.Warn("Failed to process announcement which was part of a reconcile", nil)
					}
				}
			default:
				logger.Warn("Received unknown message type", commontypes.LogFields{"msg": v})
			}
		}
	}
}

// processAnnouncement locks mu for its whole lifetime.
func (p *discoveryProtocol) processAnnouncement(ann Announcement) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.unsafeProcessAnnouncement(ann)
}

// unsafeProcessAnnouncement requires mu to be held.
func (p *discoveryProtocol) unsafeProcessAnnouncement(ann Announcement) error {
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "processAnnouncement"}).MakeChild(ann.toLogFields())
	pid, err := ann.PeerID()
	if err != nil {
		return errors.Wrap(err, "failed to obtain peer id from announcement")
	}

	if p.numGroupsByOracle[pid] == 0 {
		return fmt.Errorf("got announcement for an oracle we don't share a group with (%s)", pid)
	}

	err = ann.verify()
	if err != nil {
		return errors.Wrap(err, "failed to verify announcement")
	}

	if localann, exists := p.bestAnnouncement[pid]; !exists || localann.Counter <= ann.Counter {
		if exists && pid != p.ownID && localann.Counter == ann.Counter {
			return nil
		}
		p.bestAnnouncement[pid] = ann
		if pid == p.ownID {
			bumpedann, better, err := p.unsafeBumpOwnAnnouncement()
			if err != nil {
				return errors.Wrap(err, "failed to bump own announcement")
			}

			if !better {
				return nil
			}

			logger.Info("Received better announcement for us - bumped", nil)
			select {
			case p.chInternalBump <- *bumpedann:
			case <-p.ctx.Done():
				return nil
			}
		} else {
			logger.Info("Received better announcement for peer", nil)
			select {
			case p.chConnectivity <- connectivityMsg{connectivityAdd, pid}:
			case <-p.ctx.Done():
				return nil
			}
		}
	}

	return nil
}

func (p *discoveryProtocol) sendToAllowedPeers(ann Announcement) {
	p.mu.RLock()
	allowedPeers := p.unsafeAllowedPeers(ann)
	p.mu.RUnlock()
	for _, pid := range allowedPeers {
		select {
		case p.chOutgoingMessages <- outgoingMessage{ann, pid}:
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *discoveryProtocol) sendLoop() {
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "sendLoop"})
	logger.Debug("Entering", nil)
	defer logger.Debug("Exiting", nil)
	tick := time.After(0)
	for {
		select {
		case <-p.ctx.Done():
			return
		case ourann := <-p.chInternalBump:
			logger.Info("Our announcement was bumped - broadcasting", ourann.toLogFields())
			p.sendToAllowedPeers(ourann)
		case <-tick:
			logger.Debug("Starting reconciliation", nil)
			reconcileByPeer := make(map[ragetypes.PeerID]*reconcile)
			func() {
				p.mu.RLock()
				defer p.mu.RUnlock()
				for _, ann := range p.bestAnnouncement {
					for _, pid := range p.unsafeAllowedPeers(ann) {
						if _, exists := reconcileByPeer[pid]; !exists {
							reconcileByPeer[pid] = &reconcile{Anns: []Announcement{}}
						}
						r := reconcileByPeer[pid]
						r.Anns = append(r.Anns, ann)
					}
				}
			}()

			for pid, rec := range reconcileByPeer {
				select {
				case p.chOutgoingMessages <- outgoingMessage{rec, pid}:
					logger.Trace("Sending reconcile", commontypes.LogFields{"remotePeerID": pid, "reconcile": rec.toLogFields()})
				case <-p.ctx.Done():
					return
				}
			}
			tick = time.After(p.deltaReconcile)
		}
	}
}

// unsafeBumpOwnAnnouncement requires mu to be held by the caller.
func (p *discoveryProtocol) unsafeBumpOwnAnnouncement() (*Announcement, bool, error) {
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "unsafeBumpOwnAnnouncement"})
	oldann, exists := p.bestAnnouncement[p.ownID]
	newctr := uint64(0)

	if exists {
		if equalAddrs(oldann.Addrs, p.ownAddrs) {
			return nil, false, nil
		}
		// Counter is uint64, and it only changes when a peer's
		// addresses change. We assume a peer will not change addresses
		// more than 2**64 times.
		newctr = oldann.Counter + 1
	}
	if newctr > announcementVersionWarnThreshold {
		logger.Warn("New announcement version too big!", commontypes.LogFields{"version": newctr})
	}
	newann := unsignedAnnouncement{Addrs: p.ownAddrs, Counter: newctr}
	sann, err := newann.sign(p.privKey)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to sign own announcement")
	}
	logger.Info("Replacing our own announcement", sann.toLogFields())
	p.bestAnnouncement[p.ownID] = sann
	return &sann, true, nil
}

func (p *discoveryProtocol) Close() {
	logger := p.logger.MakeChild(commontypes.LogFields{"in": "Close"})
	logger.Debug("Exiting", nil)
	defer logger.Debug("Exited", nil)
	p.ctxCancel()
	p.processes.Wait()
}
