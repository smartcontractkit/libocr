package ragedisco

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/libocr/networking/ragedisco/serialization"
	"google.golang.org/protobuf/proto"
)

type unsignedAnnouncement struct {
	Addrs   []ragetypes.Address // addresses of a peer
	Counter uint64              // counter
}

// Announcement is a signed message in which a peer attests to their network
// addresses. An Announcement needs to adhere to some syntactic validity rules,
// found in checkWellFormed(), which are enforced on calls to verify().
type Announcement struct {
	unsignedAnnouncement
	PublicKey ed25519.PublicKey // PublicKey used to verify Sig
	Sig       []byte            // sig over unsignedAnnouncement
}

type reconcile struct {
	Anns []Announcement
}

const (
	// The maximum number of addresses an Announcement may contain
	maxAddrsInAnnouncement = 10
	// Domain separator for signatures
	announcementDomainSeparator = "announcement for chainlink peer discovery v2.0.0"
	// Maximum message size over all message types. The worst case message is a
	// reconcile message with the maximum number of announcements (one per
	// oracle, capped at MaxOracles), with each announcement containing the
	// maximum number (maxAddrsInAnnouncement) of the maximum length addresses
	// (maxAddrPortValidForAnnouncementSize). We have a test which asserts this
	// bound.
	maxMessageLength = 110_000
)

// Does NOT check if the announcement is well-formed.
func (ann Announcement) serialize() ([]byte, error) {
	pm, err := ann.toProto()
	if err != nil {
		return nil, err
	}
	return proto.Marshal(pm)
}

func (ann Announcement) toProto() (*serialization.SignedAnnouncement, error) {
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

func (uann unsignedAnnouncement) checkWellFormed() error {
	if len(uann.Addrs) == 0 || len(uann.Addrs) > maxAddrsInAnnouncement {
		return fmt.Errorf("invalid length of addresses (was %d, min is 1, max is %d)", len(uann.Addrs), maxAddrsInAnnouncement)
	}
	for _, addr := range uann.Addrs {
		if !isValidForAnnouncement(addr) {
			return fmt.Errorf("invalid address (%s)", addr)
		}
	}
	return nil
}

func (ann Announcement) checkWellFormed() error {
	if err := ann.unsignedAnnouncement.checkWellFormed(); err != nil {
		return err
	}
	const expectedPublicKeySize = ed25519.PublicKeySize
	if len(ann.PublicKey) != expectedPublicKeySize {
		return fmt.Errorf("unknown key size detected (expected %d, actual %d)", expectedPublicKeySize, len(ann.PublicKey))
	}
	if ann.Sig == nil {
		return fmt.Errorf("nil sig")
	}
	return nil
}

// Does NOT guarantee that the announcements in the reconcile contain signatures
// that verify. This must be checked independently for each announcement.
func (r reconcile) checkWellFormed() error {
	if len(r.Anns) > MaxOracles {
		return fmt.Errorf("unexpectedly many announcements (expect at most %d, actual %d)", MaxOracles, len(r.Anns))
	}
	for i, ann := range r.Anns {
		if err := ann.checkWellFormed(); err != nil {
			return fmt.Errorf("failed to validate announcement %v with index %d in reconcile: %w", ann, i, err)
		}
	}
	return nil
}

func signedAnnouncementFromProto(pm *serialization.SignedAnnouncement) (Announcement, error) {
	addrs := make([]ragetypes.Address, len(pm.Addrs))
	for i, addr := range pm.Addrs {
		addrs[i] = ragetypes.Address(addr)
	}

	ann := Announcement{
		unsignedAnnouncement{
			addrs,
			pm.Counter,
		},
		pm.PublicKey,
		pm.Sig,
	}
	return ann, nil
}

func deserializeSignedAnnouncement(binary []byte) (Announcement, error) {
	pm := serialization.SignedAnnouncement{}
	err := proto.Unmarshal(binary, &pm)
	if err != nil {
		return Announcement{}, err
	}
	return signedAnnouncementFromProto(&pm)
}

func (ann Announcement) PeerID() (ragetypes.PeerID, error) {
	return ragetypes.PeerIDFromPublicKey(ann.PublicKey)
}

func (ann Announcement) String() string {
	var identityPart string
	if pid, err := ragetypes.PeerIDFromPublicKey(ann.PublicKey); err == nil {
		identityPart = fmt.Sprintf("PeerID:%s", pid.String())
	} else {
		identityPart = fmt.Sprintf("InvalidPublicKey:%x", ann.PublicKey)
	}
	return fmt.Sprintf("{%s Counter:%d Addrs:%s Sig:%s}",
		identityPart,
		ann.Counter,
		ann.Addrs,
		base64.StdEncoding.EncodeToString(ann.Sig))
}

func (r reconcile) String() string {
	return fmt.Sprintf("%s", r.Anns)
}

// digest returns a deterministic digest used for signing
// will return an error for an invalid unsignedAnnouncement
func (uann unsignedAnnouncement) digest() ([]byte, error) {
	// serialize only addrs and the counter
	if err := uann.checkWellFormed(); err != nil {
		return nil, err
	}

	hasher := sha256.New()
	hasher.Write([]byte(announcementDomainSeparator))

	// encode addr length
	err := binary.Write(hasher, binary.LittleEndian, uint32(len(uann.Addrs)))
	if err != nil {
		return nil, err
	}
	// encode addr
	for _, a := range uann.Addrs {
		ab := []byte(a)
		err = binary.Write(hasher, binary.LittleEndian, uint32(len(ab)))
		if err != nil {
			return nil, err
		}
		hasher.Write(ab)
	}

	// counter
	err = binary.Write(hasher, binary.LittleEndian, uann.Counter)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func (uann unsignedAnnouncement) sign(keyring ragetypes.PeerKeyring) (Announcement, error) {
	digest, err := uann.digest()
	if err != nil {
		return Announcement{}, err
	}

	sig, err := keyring.Sign(digest)
	if err != nil {
		return Announcement{}, fmt.Errorf("keyring Sign failed: %w", err)
	}

	epk := ragetypes.Ed25519PublicKeyFromPeerPublicKey(keyring.PublicKey())

	return Announcement{
		uann,
		epk,
		sig,
	}, nil
}

func (ann Announcement) verify() error {
	if err := ann.checkWellFormed(); err != nil {
		return err
	}
	msg, err := ann.digest()
	if err != nil {
		return err
	}

	if !ed25519.Verify(ann.PublicKey, msg, ann.Sig) {
		return fmt.Errorf("invalid signature")
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
		return nil, fmt.Errorf("unrecognised message type %T", msg)
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

var (
	_ WrappableMessage = &reconcile{}
	_ WrappableMessage = &Announcement{}
)
