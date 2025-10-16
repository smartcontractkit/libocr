package networking

import (
	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

type ocr2BinaryNetworkEndpointFactory struct {
	*concretePeerV2
}

type ocr2BootstrapperFactory struct {
	*concretePeerV2
}

func (o *ocr2BinaryNetworkEndpointFactory) NewEndpoint(
	configDigest types.ConfigDigest,
	pids []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
	limits types.BinaryNetworkEndpointLimits,
) (commontypes.BinaryNetworkEndpoint, error) {
	return o.newEndpoint(
		configDigest,
		pids,
		v2bootstrappers,
		BinaryNetworkEndpointLimits(limits),
	)
}

func (o *ocr2BootstrapperFactory) NewBootstrapper(
	configDigest types.ConfigDigest,
	peerIDs []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
) (commontypes.Bootstrapper, error) {
	return o.newBootstrapper(
		configDigest,
		peerIDs,
		v2bootstrappers,
	)
}
