package testimplementations

import (
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type Network struct {
	mutex sync.Mutex
	chs   map[types.ConfigDigest][]chan commontypes.BinaryMessageWithSender
	n     int
}

func NewNetwork(n int) *Network {
	return &Network{
		sync.Mutex{},
		map[types.ConfigDigest][]chan commontypes.BinaryMessageWithSender{},
		n,
	}
}

func (n *Network) EndpointFactory(
	id commontypes.OracleID,
	peerID string,
) types.BinaryNetworkEndpointFactory {
	return EndpointFactory{n, id, peerID}
}

type EndpointFactory struct {
	Net     *Network
	ID      commontypes.OracleID
	PeerID_ string
}

var _ types.BinaryNetworkEndpointFactory = (*EndpointFactory)(nil)

func (epf EndpointFactory) NewEndpoint(
	configDigest types.ConfigDigest,
	_ []string,
	_ []commontypes.BootstrapperLocator,
	_ int,
	_ types.BinaryNetworkEndpointLimits,

) (
	commontypes.BinaryNetworkEndpoint, error,
) {
	epf.Net.mutex.Lock()
	defer epf.Net.mutex.Unlock()

	chs, ok := epf.Net.chs[configDigest]
	if !ok {
		chs = make([]chan commontypes.BinaryMessageWithSender, epf.Net.n)
		for i := 0; i < epf.Net.n; i++ {
			chs[i] = make(chan commontypes.BinaryMessageWithSender, 1000)
		}
		epf.Net.chs[configDigest] = chs
	}

	return EndPoint{
		epf.Net,
		configDigest,
		epf.ID,
		chs,
	}, nil
}

func (epf EndpointFactory) PeerID() string {
	return epf.PeerID_
}

type EndPoint struct {
	Net          *Network
	ConfigDigest types.ConfigDigest
	ID           commontypes.OracleID
	Chs          []chan commontypes.BinaryMessageWithSender
}

var _ commontypes.BinaryNetworkEndpoint = (*EndPoint)(nil)

func (ep EndPoint) SendTo(
	payload []byte, to commontypes.OracleID,
) {
	msg := commontypes.BinaryMessageWithSender{
		Msg:    payload,
		Sender: ep.ID,
	}
	ep.Chs[to] <- msg
}

func (ep EndPoint) Broadcast(payload []byte) {
	msg := commontypes.BinaryMessageWithSender{
		Msg:    payload,
		Sender: ep.ID,
	}
	for _, ch := range ep.Chs {
		ch <- msg
	}
}

func (ep EndPoint) Receive() <-chan commontypes.BinaryMessageWithSender {
	return ep.Chs[ep.ID]
}

func (ep EndPoint) Start() error { return nil }
func (ep EndPoint) Close() error { return nil }

func (epf EndpointFactory) MakeBootstrapper(types.ConfigDigest, []string) (commontypes.Bootstrapper, error) {
	return &bootstrapper{}, nil
}

type bootstrapper struct{}

func (*bootstrapper) Start() error { return nil }
func (*bootstrapper) Close() error { return nil }
