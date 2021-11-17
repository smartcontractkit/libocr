package testimplementations

import (
	"math"
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type InterceptedBinaryMessageWithSender struct {
	commontypes.BinaryMessageWithSender
	to           commontypes.OracleID
	configDigest types.ConfigDigest
}

type MITMNetwork struct {
	mutex    sync.Mutex
	chs      map[types.ConfigDigest][]chan commontypes.BinaryMessageWithSender
	n        int
	Listener chan InterceptedBinaryMessageWithSender
}

func NewMITMNetwork(n int) *MITMNetwork {
	return &MITMNetwork{
		sync.Mutex{},
		map[types.ConfigDigest][]chan commontypes.BinaryMessageWithSender{},
		n,
		make(chan InterceptedBinaryMessageWithSender, 1000),
	}
}

func (n *MITMNetwork) EndpointFactory(
	id commontypes.OracleID,
	peerID string,
) types.BinaryNetworkEndpointFactory {
	return MITMEndpointFactory{n, id, peerID}
}

type MITMEndpointFactory struct {
	Net     *MITMNetwork
	ID      commontypes.OracleID
	PeerID_ string
}

var _ types.BinaryNetworkEndpointFactory = (*MITMEndpointFactory)(nil)

func (epf MITMEndpointFactory) NewEndpoint(
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

	return MITMEndpoint{
		epf.Net,
		configDigest,
		epf.ID,
		chs,
	}, nil
}

func (epf MITMEndpointFactory) PeerID() string {
	return epf.PeerID_
}

type MITMEndpoint struct {
	Net          *MITMNetwork
	ConfigDigest types.ConfigDigest
	ID           commontypes.OracleID
	Chs          []chan commontypes.BinaryMessageWithSender
}

var _ commontypes.BinaryNetworkEndpoint = (*MITMEndpoint)(nil)

func (ep MITMEndpoint) SendTo(
	payload []byte, to commontypes.OracleID,
) {
	msg := commontypes.BinaryMessageWithSender{
		Msg:    payload,
		Sender: ep.ID,
	}
	ep.Chs[to] <- msg
	ep.Net.Listener <- InterceptedBinaryMessageWithSender{
		msg,
		to,
		ep.ConfigDigest,
	}
}

func (ep MITMEndpoint) Broadcast(payload []byte) {
	msg := commontypes.BinaryMessageWithSender{
		Msg:    payload,
		Sender: ep.ID,
	}
	for _, ch := range ep.Chs {
		ch <- msg
	}
	ep.Net.Listener <- InterceptedBinaryMessageWithSender{
		msg,
		commontypes.OracleID(math.MaxUint8),
		ep.ConfigDigest,
	}

}

func (ep MITMEndpoint) Receive() <-chan commontypes.BinaryMessageWithSender {
	return ep.Chs[ep.ID]
}

func (ep MITMEndpoint) Start() error { return nil }
func (ep MITMEndpoint) Close() error { return nil }

func (MITMEndpointFactory) MakeBootstrapper(types.ConfigDigest, []string) (commontypes.Bootstrapper, error) {
	return &bootstrapper{}, nil
}
