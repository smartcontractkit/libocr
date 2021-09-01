package networking

import (
	"fmt"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

func newOCREndpoint(
	networkingStack NetworkingStack,
	logger loghelper.LoggerWithContext,
	configDigest ocr2types.ConfigDigest,
	peer *concretePeer,
	v1peerIDs []p2ppeer.ID,
	v2peerIDs []ragetypes.PeerID,
	v1bootstrappers []p2ppeer.AddrInfo,
	v2bootstrappers []ragetypes.PeerInfo,
	config EndpointConfig,
	failureThreshold int,
	limits BinaryNetworkEndpointLimits,
) (commontypes.BinaryNetworkEndpoint, error) {
	if !networkingStack.subsetOf(peer.networkingStack) {
		return nil, fmt.Errorf("newOCREndpoint called with incompatible networking stack (peer has: %s, you want: %s)", peer.networkingStack, networkingStack)
	}
	v1, v2 := networkingStack.needsv1(), networkingStack.needsv2()
	if v1 && v2 {
		return newOCREndpointV1V2(
			logger,
			configDigest,
			peer,
			v1peerIDs,
			v2peerIDs,
			v1bootstrappers,
			v2bootstrappers,
			config,
			failureThreshold,
			limits,
		)
	}
	if v1 {
		return newOCREndpointV1(
			logger,
			configDigest,
			peer,
			v1peerIDs,
			v1bootstrappers,
			config,
			failureThreshold,
			limits,
		)
	}
	if v2 {
		return newOCREndpointV2(
			logger,
			configDigest,
			peer,
			v2peerIDs,
			v2bootstrappers,
			EndpointConfigV2{
				config.IncomingMessageBufferSize,
				config.OutgoingMessageBufferSize,
			},
			failureThreshold,
			limits,
		)
	}
	return nil, fmt.Errorf("peer is neither v1 nor v2")
}
