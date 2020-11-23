

package dhtrouter

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/smartcontractkit/libocr/networking/dht-router/serialization"
)


type announcementCounter struct {
	userPrefix uint32 
	value      uint64 
}

func (n announcementCounter) Gt(other announcementCounter) bool {
	if n.userPrefix > other.userPrefix {
		return true
	}
	return n.userPrefix == other.userPrefix && n.value > other.value
}

type announcement struct {
	Addrs   []ma.Multiaddr      
	Counter announcementCounter 
}

type signedAnnouncement struct {
	announcement
	PublicKey p2pcrypto.PubKey 
	Sig       []byte           
}

const (
	
	maxAddrInAnnouncements = 10
	
	announcementDomainSeparator = "announcement OCR v1.0.0"
)

func serdeError(field string) error {
	return fmt.Errorf("invalid pm: %s", field)
}


func (ann signedAnnouncement) serialize() ([]byte, error) {
	
	if ann.Addrs == nil || ann.PublicKey == nil || ann.Sig == nil || len(ann.Addrs) > maxAddrInAnnouncements {
		return nil, errors.New("invalid announcement")
	}

	
	err := ann.verify()
	if err != nil {
		return nil, err
	}

	
	var addrs [][]byte
	for _, a := range ann.Addrs {
		addrBytes, err := a.MarshalBinary()
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addrBytes)
	}

	pkBytes, err := p2pcrypto.MarshalPublicKey(ann.PublicKey)
	if err != nil {
		return nil, err
	}

	pm := serialization.SignedAnnouncement{
		Addrs:      addrs,
		UserPrefix: ann.Counter.userPrefix,
		Counter:    ann.Counter.value,
		PublicKey:  pkBytes,
		Sig:        ann.Sig,
	}

	return proto.Marshal(&pm)
}

func deserializeSignedAnnouncement(binary []byte) (signedAnnouncement, error) {
	pm := serialization.SignedAnnouncement{}
	err := proto.Unmarshal(binary, &pm)
	if err != nil {
		return signedAnnouncement{}, err
	}

	
	if len(pm.Addrs) == 0 {
		return signedAnnouncement{}, serdeError("addrs is empty array")
	}

	if len(pm.Addrs) > maxAddrInAnnouncements {
		return signedAnnouncement{}, serdeError("more addr than maxAddrInAnnouncements")
	}

	var addrs []ma.Multiaddr
	for _, addr := range pm.Addrs {
		mAddr, err := ma.NewMultiaddrBytes(addr)
		if err != nil {
			return signedAnnouncement{}, err
		}
		addrs = append(addrs, mAddr)
	}

	publicKey, err := p2pcrypto.UnmarshalPublicKey(pm.PublicKey)
	if err != nil {
		return signedAnnouncement{}, err
	}

	return signedAnnouncement{
		announcement{
			addrs,
			announcementCounter{
				pm.UserPrefix,
				pm.Counter,
			},
		},
		publicKey,
		pm.Sig,
	}, nil
}

func (ann signedAnnouncement) String() string {
	pkStr := "can't stringify PK"
	if b, err := ann.PublicKey.Bytes(); err == nil {
		pkStr = base64.StdEncoding.EncodeToString(b)
	}
	return fmt.Sprintf("addrs=%s, pk=%s, sig=%s",
		ann.Addrs,
		pkStr,
		base64.StdEncoding.EncodeToString(ann.Sig))
}


func (ann announcement) digest() ([]byte, error) {
	
	if ann.Addrs == nil || len(ann.Addrs) > maxAddrInAnnouncements {
		return nil, errors.New("invalid announcement")
	}

	hasher := sha256.New()
	hasher.Write([]byte(announcementDomainSeparator))

	
	err := binary.Write(hasher, binary.LittleEndian, uint32(len(ann.Addrs)))
	if err != nil {
		return nil, err
	}
	
	for _, a := range ann.Addrs {
		addr, err := a.MarshalBinary()
		if err != nil {
			return nil, err
		}
		err = binary.Write(hasher, binary.LittleEndian, uint32(len(addr)))
		if err != nil {
			return nil, err
		}
		hasher.Write(addr)
	}

	
	err = binary.Write(hasher, binary.LittleEndian, ann.Counter.userPrefix)
	if err != nil {
		return nil, err
	}
	err = binary.Write(hasher, binary.LittleEndian, ann.Counter.value)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func (ann *announcement) sign(sk p2pcrypto.PrivKey) (signedAnnouncement, error) {
	digest, err := ann.digest()
	if err != nil {
		return signedAnnouncement{}, err
	}

	sig, err := sk.Sign(digest)
	if err != nil {
		return signedAnnouncement{}, err
	}

	return signedAnnouncement{
		*ann,
		sk.GetPublic(),
		sig,
	}, nil
}

func (ann signedAnnouncement) verify() error {
	if ann.Sig == nil {
		return errors.New("nil sig")
	}

	msg, err := ann.digest()
	if err != nil {
		return err
	}

	verified, err := ann.PublicKey.Verify(msg, ann.Sig)
	if err != nil {
		return err
	}

	if !verified {
		return errors.New("invalid signature")
	}

	return nil
}
