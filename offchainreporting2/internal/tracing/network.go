package tracing

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type binaryMessageWithTwinSender struct {
	message commontypes.BinaryMessageWithSender
	sender  OracleID
}

// Network is a testing utility that simulates an oracle network.
// Features:
// - it can simulate network latencies by introducing delays in transmission between two oracles. These can be changed dynamically.
// - it can change network connectivity dynamically, which lets you introduce partitions or isolate nodes
// - it can introduce twin oracles into the network
// - it traces messages sent, dropped, broadcasted and received for each oracles.
// What it doesn't do.. yet:
// - it doesn't support multiple ConfigDigest endpoints
// - it can't add or remove oracles dynamically
// There may be confusion between the use of tracing.OracleID, commontypes.OracleID and int. The rules are:
// - use commontypes.OracleID when we implement the interfaces in the "types" package. Not much we can do about this!
// - use ints in the public interface of this network. This makes tests short and expressive.
// - use tracing.OracleID internally to eliminate confusion and simplify the implemenation.
type Network struct {
	ctx       context.Context
	cancelCtx context.CancelFunc

	tracer *Tracer

	ioMu sync.Mutex
	// inbounds are channels to read data sent to an oracle
	// these channels are exposed to the oracles in their endpoint's Receive() method.
	inbounds map[OracleID]chan binaryMessageWithTwinSender

	// chokes is a mapping of artificially introduced drops of connectivity between a pair of oracles.
	// If there is a choke AND a latency between two oracles, the latency is ignored.
	// Chokes are used to implement partitions and other types of network topologies
	chokesMu sync.Mutex
	chokes   map[OracleID]map[OracleID]struct{}

	// latencies is a mapping of artificially increased latencies in the communication between two oracles.
	// Latencies are used to emulate network congestion.
	latenciesMu sync.Mutex
	latencies   map[OracleID]map[OracleID]time.Duration

	// twinsRoundRobin has a double function:
	// 1. if there is a oracleID in the map it means it has a twin.
	// 2. the corresponding value is used to pick one of the twins to send a message to using round robin.
	twinsRoundRobinMu sync.Mutex
	twinsRoundRobin   map[commontypes.OracleID]int

	config NetworkConfig
}

type NetworkConfig struct {
	RoundRobinBetweenTwins bool
}

func NewNetwork(tracer *Tracer) *Network {
	localCtx, cancelCtx := context.WithCancel(context.Background())
	return &Network{
		ctx:       localCtx,
		cancelCtx: cancelCtx,

		tracer: tracer,

		inbounds: make(map[OracleID]chan binaryMessageWithTwinSender),

		chokes:    make(map[OracleID]map[OracleID]struct{}),
		latencies: make(map[OracleID]map[OracleID]time.Duration),

		twinsRoundRobin: make(map[commontypes.OracleID]int),

		config: NetworkConfig{},
	}
}

// SetConfig is not thread safe.. for now.
func (n *Network) SetConfig(newConfig NetworkConfig) {
	n.config = newConfig
}

// Close terminates all the goroutines used by the Network and its endpoints.
func (n *Network) Close() {
	n.cancelCtx()
}

func (n *Network) addOracle(id OracleID) {
	// initialize oid's channels
	n.ioMu.Lock()
	defer n.ioMu.Unlock()
	n.inbounds[id] = make(chan binaryMessageWithTwinSender)
	// initialize the inbound round robin for twin oracles.
	if !id.IsTwin {
		return
	}
	n.twinsRoundRobinMu.Lock()
	n.twinsRoundRobin[id.OracleID] = 0
	n.twinsRoundRobinMu.Unlock()
}

func (n *Network) AddChoke(src, dst int) {
	srcID := FromInt(src)
	dstID := FromInt(dst)
	n.chokesMu.Lock()
	defer n.chokesMu.Unlock()
	if _, found := n.chokes[srcID]; !found {
		n.chokes[srcID] = make(map[OracleID]struct{})
	}
	n.chokes[srcID][dstID] = struct{}{}
}

func (n *Network) RemoveChoke(src, dst int) {
	srcID := FromInt(src)
	dstID := FromInt(dst)
	n.chokesMu.Lock()
	defer n.chokesMu.Unlock()
	if _, found := n.chokes[srcID]; !found {
		return
	}
	delete(n.chokes[srcID], dstID)
}

func (n *Network) isChoked(src, dst OracleID) bool {
	n.chokesMu.Lock()
	defer n.chokesMu.Unlock()
	chokes, found := n.chokes[src]
	if !found { // src can talk to everybody
		return false
	}
	_, found = chokes[dst] // if dst in the list of choked connections from src.
	return found
}

func (n *Network) AddLatency(src, dst int, latency time.Duration) {
	srcID := FromInt(src)
	dstID := FromInt(dst)
	n.latenciesMu.Lock()
	defer n.latenciesMu.Unlock()
	if _, found := n.latencies[srcID]; !found {
		n.latencies[srcID] = make(map[OracleID]time.Duration)
	}
	n.latencies[srcID][dstID] = latency
}

func (n *Network) RemoveLatency(src, dst int) {
	srcID := FromInt(src)
	dstID := FromInt(dst)
	n.latenciesMu.Lock()
	defer n.latenciesMu.Unlock()
	if _, found := n.latencies[srcID]; !found {
		return
	}
	delete(n.latencies[srcID], dstID)
}

// simulateLatency will sleep if there is a "latency" configured between src and dst oracles.
func (n *Network) simulateLatency(src, dst OracleID) {
	var found bool
	var latency time.Duration
	func() {
		n.latenciesMu.Lock()
		defer n.latenciesMu.Unlock()
		if _, found := n.latencies[src]; !found {
			return
		}
		latency, found = n.latencies[src][dst]
	}()
	if !found {
		return
	}
	select {
	case <-time.After(latency):
	case <-n.ctx.Done():
	}
}

// sendTo checks if the oracle has a twin and sends the message to the twin as well.
func (n *Network) sendTo(src, dst OracleID, msg commontypes.BinaryMessageWithSender) {
	if n.config.RoundRobinBetweenTwins {
		dst := n.roundRobin(src, dst)
		n.asyncSendTo(src, dst, msg)
		return
	}
	n.asyncSendTo(src, dst, msg)
	// Check if there dst has a twin, if it does send msg to it. send() takes care of any latencies or chokes.
	n.twinsRoundRobinMu.Lock()
	_, hasTwin := n.twinsRoundRobin[dst.OracleID]
	n.twinsRoundRobinMu.Unlock()
	if hasTwin {
		n.asyncSendTo(src, dst.Twin(), msg)
	}
}

// roundRobin decides whether to return dst or its twin.
// There are four cases:
// 1. dst doesnt have a twin the src's partition => just return dst
// 2. dst has a twin, but it's not reachable from src (because of separate partitions) => return dst
// 3. dst has a reachable twin but dst itself is not reachable from src => return dst's twin
// 4. both dst and it's twin are reachable from src => round robin from the two.
func (n *Network) roundRobin(src, dst OracleID) OracleID {
	n.twinsRoundRobinMu.Lock()
	_, hasTwin := n.twinsRoundRobin[dst.OracleID]
	n.twinsRoundRobinMu.Unlock()
	// case 1:
	if !hasTwin {
		return dst
	}
	canConnectToDst := n.isChoked(src, dst)
	canConnectToTwinDst := n.isChoked(src, dst.Twin())
	// case 2:
	if canConnectToDst && !canConnectToTwinDst {
		return dst
	}
	// case 3:
	if !canConnectToDst && canConnectToTwinDst {
		return dst.Twin()
	}
	// case 4:
	options := []OracleID{dst, dst.Twin()}
	n.twinsRoundRobinMu.Lock()
	defer n.twinsRoundRobinMu.Unlock()
	newDst := options[n.twinsRoundRobin[dst.OracleID]]
	n.twinsRoundRobin[dst.OracleID] = (n.twinsRoundRobin[dst.OracleID] + 1) % len(options)
	return newDst
}

// asyncSendTo will launch a goroutine to send msg on the inbound channel corresponding to dst.
// It also records appropriate tracing frames.
func (n *Network) asyncSendTo(src, dst OracleID, msg commontypes.BinaryMessageWithSender) {
	n.ioMu.Lock()
	inbound := n.inbounds[dst]
	n.ioMu.Unlock()
	go func(src, dst OracleID, inbound chan binaryMessageWithTwinSender) {
		n.simulateLatency(src, dst)
		n.tracer.Append(NewSendTo(src, src, dst, msg.Msg))
		msgWithTwinSender := binaryMessageWithTwinSender{msg, src}
		select {
		case inbound <- msgWithTwinSender:
		case <-n.ctx.Done():
		}
	}(src, dst, inbound)
}

func (n *Network) broadcast(src OracleID, msg commontypes.BinaryMessageWithSender) {
	n.tracer.Append(NewBroadcast(src, src, msg.Msg))
	for dst := range n.inbounds {
		if dst.IsTwin {
			continue // asyncSendTo will handle sending messages to twins in a round-robin pattern. Here we just skip twins.
		}
		n.sendTo(src, dst, msg)
	}
}

func (n *Network) getReceiver(dst OracleID) chan commontypes.BinaryMessageWithSender {
	inbound := n.inbounds[dst]
	ch := make(chan commontypes.BinaryMessageWithSender)
	go func() {
		for {
			select {
			case wrappedMessage := <-inbound:
				src, msg := wrappedMessage.sender, wrappedMessage.message
				if n.isChoked(src, dst) {
					// src can't reach dst.
					n.tracer.Append(NewDrop(src, src, dst, msg.Msg))
					continue
				}
				select {
				case ch <- msg:
					n.tracer.Append(NewReceive(dst, src, dst, msg.Msg))
				case <-n.ctx.Done():
					return
				}
			case <-n.ctx.Done():
				return
			}
		}
	}()
	return ch
}

func (n *Network) NewEndpointFactory(oracleID OracleID, peerID string) *EndpointFactory {
	n.addOracle(oracleID)
	return &EndpointFactory{
		n,
		oracleID,
		peerID,
	}
}

// EndpointFactory produces endpoints that wrap Network's channels.
type EndpointFactory struct {
	network *Network
	id      OracleID
	peerID  string
}

var _ types.BinaryNetworkEndpointFactory = (*EndpointFactory)(nil)

func (e *EndpointFactory) PeerID() string {
	return e.peerID
}

func (e *EndpointFactory) NewEndpoint(
	_ types.ConfigDigest,
	_ []string, // peerID
	_ []commontypes.BootstrapperLocator, // v2bootstrappers
	_ int, // failureThreshold
	_ types.BinaryNetworkEndpointLimits, // limits
) (commontypes.BinaryNetworkEndpoint, error) {
	return &Endpoint{
		e.network,
		e.id,
	}, nil
}

// Endpoint
type Endpoint struct {
	network *Network
	id      OracleID
}

var _ commontypes.BinaryNetworkEndpoint = (*Endpoint)(nil)

// SendTo(msg, to) sends msg to "to"
func (e *Endpoint) SendTo(payload []byte, to commontypes.OracleID) {
	msg := commontypes.BinaryMessageWithSender{
		payload,       // Msg
		e.id.OracleID, // Sender: we send messages with the shared id for both the oracle and its twin.
	}
	e.network.sendTo(e.id, OracleID{to, false}, msg)
}

// Broadcast(msg) sends msg to all oracles
func (e *Endpoint) Broadcast(payload []byte) {
	msg := commontypes.BinaryMessageWithSender{
		payload,       // Msg
		e.id.OracleID, // Sender: we send messages with the shared id for both the oracle and its twin.
	}
	e.network.broadcast(e.id, msg)
}

// Receive returns channel which carries all messages sent to this oracle.
func (e *Endpoint) Receive() <-chan commontypes.BinaryMessageWithSender {
	return e.network.getReceiver(e.id)
}

// Start starts the endpoint
func (e *Endpoint) Start() error {
	e.network.tracer.Append(NewEndpointStart(e.id))
	return nil
}

// Close stops the endpoint
func (e *Endpoint) Close() error {
	e.network.tracer.Append(NewEndpointClose(e.id))
	return nil
}

// Public Utitlities

// SplitNetwork introduces partitions in the network layer.
// The partitions can overlap!
func SplitNetwork(network *Network, partition1, partition2 []int) {
	p1MinusP2 := aMinusB(partition1, partition2)
	p2MinusP1 := aMinusB(partition2, partition1)
	for _, id1 := range p1MinusP2 {
		for _, id2 := range p2MinusP1 {
			network.AddChoke(id1, id2)
			network.AddChoke(id2, id1)
		}
	}
}

// IsolateNode configures the network to choke all connections to and from a given node.
func IsolateNode(network *Network, this int) {
	thisOID := FromInt(this)
	for thatOID := range network.inbounds {
		if thisOID == thatOID {
			continue // A node can always reach itself?!
		}
		that := int(thatOID.OracleID)
		network.AddChoke(this, that)
		network.AddChoke(that, this)
	}
}

// Helpers

func aMinusB(as, bs []int) []int {
	out := []int{}
	for _, a := range as {
		found := false
		for _, b := range bs {
			if a == b {
				found = true
				break
			}
		}
		if !found {
			out = append(out, a)
		}
	}
	return out
}
