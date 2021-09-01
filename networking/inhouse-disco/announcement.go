//go:generate protoc -I. --go_out=./serialization  ./serialization/peer_discovery_announcement.proto

package inhousedisco

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	rageaddress "github.com/smartcontractkit/libocr/ragep2p/address"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/networking/inhouse-disco/serialization"
	"google.golang.org/protobuf/proto"
)

type unsignedAnnouncement struct {
	Addrs   []ragetypes.Address // addresses of a peer
	Counter uint64              // counter
}

type Announcement struct {
	unsignedAnnouncement
	PublicKey ed25519.PublicKey // PublicKey used to verify Sig
	Sig       []byte            // sig over unsignedAnnouncement
}

type reconcile struct {
	Anns []Announcement
}

const (
	// The maximum number of addr an Announcement may broadcast
	maxAddrInAnnouncements = 10
	// Domain separator for signatures
	announcementDomainSeparator = "announcement for chainlink peer discovery v2.0.0"
	// Maximum message size over all message types. Should be able to
	// handle the equivalent of a reconcile with 1000 announcements of 1
	// address each. Considering our committees are typically of size 32
	// and we limit to 10 addresses per announcement we are overshooting by
	// (at the very least) ~3x. We have a test which asserts the tightness
	// of this bound.
	maxMessageLength = 110000
)

func serdeError(field string) error {
	return fmt.Errorf("invalid pm: %s", field)
}

// Validate and serialize an Announcement. Return error on invalid announcements.
func (ann Announcement) serialize() ([]byte, error) {
	pm, err := ann.toProto()
	if err != nil {
		return nil, err
	}
	return proto.Marshal(pm)
}

func (ann Announcement) toProto() (*serialization.SignedAnnouncement, error) {
	// Require all fields to be non-nil and addrs shorter than maxAddrInAnnouncements
	if ann.Addrs == nil || ann.PublicKey == nil || ann.Sig == nil || len(ann.Addrs) > maxAddrInAnnouncements {
		return nil, errors.New("invalid announcement")
	}

	// verify the signature
	err := ann.verify()
	if err != nil {
		return nil, err
	}

	// addr
	var addrs [][]byte
	for _, a := range ann.Addrs {
		addrs = append(addrs, []byte(a))
	}

	pm := serialization.SignedAnnouncement{
		Addrs:     addrs,
		Counter:   ann.Counter,
		PublicKey: ann.PublicKey,
		Sig:       ann.Sig,
	}
	return &pm, nil
}

func signedAnnouncementFromProto(pm *serialization.SignedAnnouncement) (Announcement, error) {
	// public key
	const expectedPublicKeySize = ed25519.PublicKeySize
	if len(pm.PublicKey) != expectedPublicKeySize {
		return Announcement{}, fmt.Errorf("unknown key size detected (expected %d, actual %d)", expectedPublicKeySize, len(pm.PublicKey))
	}

	// addrs
	if len(pm.Addrs) == 0 {
		return Announcement{}, serdeError("addrs is empty array")
	}

	if len(pm.Addrs) > maxAddrInAnnouncements {
		return Announcement{}, serdeError("more addr than maxAddrInAnnouncements")
	}

	var addrs []ragetypes.Address
	for _, addr := range pm.Addrs {
		raddr := ragetypes.Address(addr)
		if !rageaddress.IsValid(raddr) {
			return Announcement{}, fmt.Errorf("contained invalid address (%s)", addr)
		}
		addrs = append(addrs, raddr)
	}

	return Announcement{
		unsignedAnnouncement{
			addrs,
			pm.Counter,
		},
		pm.PublicKey,
		pm.Sig,
	}, nil
}

func deserializeSignedAnnouncement(binary []byte) (Announcement, error) {
	pm := serialization.SignedAnnouncement{}
	err := proto.Unmarshal(binary, &pm)
	if err != nil {
		return Announcement{}, err
	}
	return signedAnnouncementFromProto(&pm)
}

func (ann unsignedAnnouncement) String() string {
	return fmt.Sprintf("<counter=%d, addrs=%s>",
		ann.Counter,
		ann.Addrs)
}

func (ann Announcement) PeerID() (ragetypes.PeerID, error) {
	return ragetypes.PeerIDFromPublicKey(ann.PublicKey)
}

func (ann Announcement) String() string {
	pkStr := base64.StdEncoding.EncodeToString(ann.PublicKey)
	return fmt.Sprintf("<counter=%d, addrs=%s, pk=%s, sig=%s>",
		ann.Counter,
		ann.Addrs,
		pkStr,
		base64.StdEncoding.EncodeToString(ann.Sig))
}

// digest returns a deterministic digest used for signing
func (ann unsignedAnnouncement) digest() ([]byte, error) {
	// serialize only addrs and the counter
	if ann.Addrs == nil || len(ann.Addrs) > maxAddrInAnnouncements {
		return nil, errors.New("invalid announcement")
	}

	hasher := sha256.New()
	hasher.Write([]byte(announcementDomainSeparator))

	// encode addr length
	err := binary.Write(hasher, binary.LittleEndian, uint32(len(ann.Addrs)))
	if err != nil {
		return nil, err
	}
	// encode addr
	for _, a := range ann.Addrs {
		ab := []byte(a)
		err = binary.Write(hasher, binary.LittleEndian, uint32(len(ab)))
		if err != nil {
			return nil, err
		}
		hasher.Write(ab)
	}

	// counter
	err = binary.Write(hasher, binary.LittleEndian, ann.Counter)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func (ann *unsignedAnnouncement) sign(sk ed25519.PrivateKey) (Announcement, error) {
	digest, err := ann.digest()
	if err != nil {
		return Announcement{}, err
	}

	sig := ed25519.Sign(sk, digest)

	epk, ok := sk.Public().(ed25519.PublicKey)
	if !ok {
		return Announcement{}, errors.New("public key is not ed25519 public key")
	}

	return Announcement{
		*ann,
		epk,
		sig,
	}, nil
}

func (ann Announcement) verify() error {
	if ann.Sig == nil {
		return errors.New("nil sig")
	}

	msg, err := ann.digest()
	if err != nil {
		return err
	}

	verified := ed25519.Verify(ann.PublicKey, msg, ann.Sig)
	if !verified {
		return errors.New("invalid signature")
	}

	return nil
}

func (r reconcile) toProto() (*serialization.Reconcile, error) {
	serAnns := make([]*serialization.SignedAnnouncement, len(r.Anns))
	for i, ann := range r.Anns {
		protoAnn, err := ann.toProto()
		if err != nil {
			return nil, err
		}
		serAnns[i] = protoAnn
	}

	ser := serialization.Reconcile{
		Anns: serAnns,
	}
	return &ser, nil
}

func (r reconcile) toProtoWrapped() (*serialization.MessageWrapper, error) {
	rProto, err := r.toProto()
	if err != nil {
		return nil, err
	}
	msgWrapper := serialization.MessageWrapper{}
	msgWrapper.Msg = &serialization.MessageWrapper_MessageReconcile{rProto}
	return &msgWrapper, nil
}

func reconcileFromProto(pr *serialization.Reconcile) (*reconcile, error) {
	anns := make([]Announcement, len(pr.Anns))
	for i, protoAnn := range pr.Anns {
		ann, err := signedAnnouncementFromProto(protoAnn)
		if err != nil {
			return nil, err
		}
		anns[i] = ann
	}
	return &reconcile{Anns: anns}, nil
}

func (ann Announcement) toProtoWrapped() (*serialization.MessageWrapper, error) {
	annProto, err := ann.toProto()
	if err != nil {
		return nil, err
	}
	msgWrapper := serialization.MessageWrapper{}
	msgWrapper.Msg = &serialization.MessageWrapper_MessageSignedAnnouncement{annProto}
	return &msgWrapper, nil
}

func fromProtoWrappedBytes(b []byte) (WrappableMessage, error) {
	wrapper := &serialization.MessageWrapper{}
	err := proto.Unmarshal(b, wrapper)
	if err != nil {
		return nil, err
	}

	switch msg := wrapper.Msg.(type) {
	case *serialization.MessageWrapper_MessageReconcile:
		rec, err := reconcileFromProto(wrapper.GetMessageReconcile())
		if err != nil {
			return nil, err
		}
		return rec, nil
	case *serialization.MessageWrapper_MessageSignedAnnouncement:
		ann, err := signedAnnouncementFromProto(wrapper.GetMessageSignedAnnouncement())
		if err != nil {
			return nil, err
		}
		return &ann, nil
	default:
		return nil, errors.Errorf("Unrecognised Msg type %T", msg)
	}
}

type WrappableMessage interface {
	toProtoWrapped() (*serialization.MessageWrapper, error)
}

func toBytesWrapped(m WrappableMessage) ([]byte, error) {
	p, err := m.toProtoWrapped()
	if err != nil {
		return nil, err
	}
	return proto.Marshal(p)
}

func (ann Announcement) toLogFields() commontypes.LogFields {
	pid, err := ragetypes.PeerIDFromPublicKey(ann.PublicKey)
	if err != nil {
		return commontypes.LogFields{
			"ann_publicKey": ann.PublicKey,
			"ann_ver":       ann.Counter,
			"ann_addrs":     ann.Addrs,
		}
	}
	return commontypes.LogFields{
		"ann_pid":   pid,
		"ann_ver":   ann.Counter,
		"ann_addrs": ann.Addrs,
	}
}

func (r reconcile) toLogFields() commontypes.LogFields {
	var annsLogFields []commontypes.LogFields
	for _, ann := range r.Anns {
		annsLogFields = append(annsLogFields, ann.toLogFields())
	}
	return commontypes.LogFields{
		"anns": annsLogFields,
	}
}

var (
	_ WrappableMessage = &reconcile{}
	_ WrappableMessage = &Announcement{}
)
