package networking

import (
	"fmt"
	"io"
	"sync"

	"github.com/RoSpaceDev/libocr/commontypes"
	ocr2types "github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	ragetypes "github.com/RoSpaceDev/libocr/ragep2p/types"

	"github.com/RoSpaceDev/libocr/internal/loghelper"
)

var (
	_ commontypes.Bootstrapper = &bootstrapperV2{}
)

type bootstrapperState int

const (
	_ bootstrapperState = iota
	bootstrapperUnstarted
	bootstrapperStarted
	bootstrapperClosed
)

type bootstrapperV2 struct {
	logger       loghelper.LoggerWithContext
	configDigest ocr2types.ConfigDigest
	registration io.Closer
	state        bootstrapperState

	stateMu *sync.Mutex
}

func newBootstrapperV2(
	logger loghelper.LoggerWithContext,
	configDigest ocr2types.ConfigDigest,
	v2peerIDs []ragetypes.PeerID,
	v2bootstrappers []ragetypes.PeerInfo,
	registration io.Closer,
) (*bootstrapperV2, error) {
	logger = logger.MakeChild(commontypes.LogFields{
		"id":           "bootstrapperV2",
		"configDigest": configDigest.Hex(),
	})

	logger.Info("BootstrapperV2: Initialized", commontypes.LogFields{
		"bootstrappers": v2bootstrappers,
		"oracles":       v2peerIDs,
	})

	return &bootstrapperV2{
		logger,
		configDigest,
		registration,
		bootstrapperUnstarted,
		new(sync.Mutex),
	}, nil
}

func (b *bootstrapperV2) Start() error {
	succeeded := false
	defer func() {
		if !succeeded {
			b.Close()
		}
	}()

	b.stateMu.Lock()
	defer b.stateMu.Unlock()

	if b.state != bootstrapperUnstarted {
		return fmt.Errorf("cannot start bootstrapperV2 that is not unstarted, state was: %d", b.state)
	}

	b.state = bootstrapperStarted

	b.logger.Info("BootstrapperV2: Started listening", nil)
	succeeded = true
	return nil
}

func (b *bootstrapperV2) Close() error {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()
	if b.state != bootstrapperStarted {
		return fmt.Errorf("cannot close bootstrapperV2 that is not started, state was: %d", b.state)
	}
	b.state = bootstrapperClosed

	if err := b.registration.Close(); err != nil {
		return fmt.Errorf("could not unregister bootstrapperV2: %w", err)
	}
	return nil
}
