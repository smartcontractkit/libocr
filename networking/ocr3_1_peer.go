package networking

import (
	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

var _ types.BinaryNetworkEndpoint2Factory = &ocr3_1BinaryNetworkEndpointFactory{}

type ocr3_1BinaryNetworkEndpointFactory struct {
	*concretePeerV2
}

func (o *ocr3_1BinaryNetworkEndpointFactory) NewEndpoint(
	configDigest types.ConfigDigest,
	pids []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	defaultPriorityConfig types.BinaryNetworkEndpoint2Config,
	lowPriorityConfig types.BinaryNetworkEndpoint2Config,
) (types.BinaryNetworkEndpoint2, error) {
	return o.newEndpoint3_1(
		configDigest,
		pids,
		v2bootstrappers,
		defaultPriorityConfig,
		lowPriorityConfig,
	)
}
