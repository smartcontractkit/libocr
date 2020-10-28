package knockingtls

import (
	"context"
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"sync"

	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/sec"
	p2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/smartcontractkit/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"golang.org/x/crypto/ed25519"
)

const ID = "cl_knockingtls/1.0.0"
const domainSeparator = "knockknock" + ID

type KnockingTLSTransport struct {
	tls            *p2ptls.Transport 
	allowlistMutex sync.RWMutex      
	allowlist      []peer.ID         
	privateKey     *p2pcrypto.Ed25519PrivateKey
	myId           peer.ID
	logger         types.Logger
}

var (
	errInvalidConnection = errors.New("invalid connection")
	errInvalidSignature  = errors.New("invalid signature in knock")
)

func buildKnockMessage(p peer.ID) ([]byte, error) {
	
	if len(p.Pretty()) > 128 {
		return nil, errors.New("too big id. looks suspicious")
	}
	h := crypto.SHA256.New()
	h.Write([]byte(domainSeparator))
	h.Write([]byte(p.Pretty()))

	return h.Sum(nil), nil
}

func (c *KnockingTLSTransport) SecureInbound(ctx context.Context, insecure net.Conn) (sec.SecureConn, error) {
	
	shouldClose := true
	defer func() {
		if shouldClose {
			insecure.Close()
		}
	}()

	
	
	const knockSize = ed25519.PublicKeySize + ed25519.SignatureSize
	knock := make([]byte, knockSize)

	logger := loghelper.MakeLoggerWithContext(c.logger, types.LogFields{
		"remoteAddr": insecure.RemoteAddr(),
		"localAddr":  insecure.LocalAddr(),
	})

	n, err := insecure.Read(knock)
	if err != nil {
		return nil, fmt.Errorf("can't read sig: %w", err)
	}

	if n < knockSize {
		
		
		
		return nil, fmt.Errorf("can't read sig: %w", err)
	}

	pk, err := p2pcrypto.UnmarshalEd25519PublicKey(knock[:ed25519.PublicKeySize])
	if err != nil {
		return nil, err
	}

	peerId, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return nil, err
	}

	inAllowList := false
	
	func() {
		c.allowlistMutex.RLock()
		defer c.allowlistMutex.RUnlock()

		logger.Trace("verifying a knock", types.LogFields{
			"allowlist": c.allowlist,
			"fromId":    peerId.Pretty(),
			"knock":     hex.EncodeToString(knock),
		})

		for i := range c.allowlist {
			if peerId == c.allowlist[i] {
				inAllowList = true
				break
			}
		}
	}()

	if !inAllowList {
		return nil, errors.New(fmt.Sprintf("remote peer %s not in the list", peerId.Pretty()))
	}

	knockMsg, err := buildKnockMessage(c.myId)
	if err != nil {
		return nil, err
	}

	verified, err := pk.Verify(knockMsg, knock[ed25519.PublicKeySize:])
	if err != nil {
		return nil, err
	}

	if !verified {
		return nil, errInvalidSignature
	}

	
	shouldClose = false

	return c.tls.SecureInbound(ctx, insecure)
}

func (c *KnockingTLSTransport) SecureOutbound(ctx context.Context, insecure net.Conn, p peer.ID) (sec.SecureConn, error) {
	
	shouldClose := true
	defer func() {
		if shouldClose {
			insecure.Close()
		}
	}()

	pk, err := c.privateKey.GetPublic().Raw()
	if err != nil || len(pk) != ed25519.PublicKeySize {
		return nil, errors.New("can't get PK")
	}

	knockMsg, err := buildKnockMessage(p)
	if err != nil {
		return nil, err
	}
	sig, err := c.privateKey.Sign(knockMsg)
	if err != nil || len(sig) != ed25519.SignatureSize {
		return nil, errors.New("can't sign")
	}

	
	knock := append(pk, sig...)

	n, err := insecure.Write(knock)
	if err != nil {
		return nil, err
	}
	if n != len(pk)+len(sig) {
		return nil, errors.New("can't send all tag")
	}

	
	shouldClose = false
	return c.tls.SecureOutbound(ctx, insecure, p)
}

func (c *KnockingTLSTransport) UpdateAllowlist(allowlist []peer.ID) {
	c.allowlistMutex.Lock()
	defer c.allowlistMutex.Unlock()

	c.logger.Debug("allowlist updated", types.LogFields{
		"old": c.allowlist,
		"new": allowlist,
	})
	c.allowlist = allowlist
}


func NewKnockingTLS(logger types.Logger, myPrivKey p2pcrypto.PrivKey, allowlist ...peer.ID) (*KnockingTLSTransport, error) {
	ed25515Key, ok := myPrivKey.(*p2pcrypto.Ed25519PrivateKey)
	if !ok {
		return nil, errors.New("only support ed25519 key")
	}
	if allowlist == nil {
		allowlist = []peer.ID{}
	}

	tls, err := p2ptls.New(myPrivKey)
	if err != nil {
		return nil, err
	}

	id, err := peer.IDFromPrivateKey(myPrivKey)
	if err != nil {
		return nil, err
	}

	return &KnockingTLSTransport{
		tls:            tls,
		allowlistMutex: sync.RWMutex{},
		allowlist:      allowlist,
		privateKey:     ed25515Key,
		myId:           id,
		logger: loghelper.MakeLoggerWithContext(logger, types.LogFields{
			"id": "KnockingTLS",
		}),
	}, nil
}

var _ sec.SecureTransport = &KnockingTLSTransport{}
