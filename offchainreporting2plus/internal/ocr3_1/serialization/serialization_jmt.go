package serialization

import (
	"encoding/binary"
	"fmt"

	"github.com/RoSpaceDev/libocr/internal/jmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// SerializeJmtNode converts a jmt.Node to bytes
func SerializeJmtNode(node jmt.Node) ([]byte, error) {
	if node == nil {
		return nil, fmt.Errorf("cannot serialize nil jmt node")
	}

	pb := &Node{}
	switch n := node.(type) {
	case *jmt.InternalNode:
		internalNode, err := serializeInternalNode(n)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize internal node: %w", err)
		}
		pb.Node = &Node_InternalNode{internalNode}
	case *jmt.LeafNode:
		leafNode, err := serializeLeafNode(n)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize leaf node: %w", err)
		}
		pb.Node = &Node_LeafNode{leafNode}
	default:
		return nil, fmt.Errorf("unknown jmt node type: %T", node)
	}

	return proto.Marshal(pb)
}

// DeserializeJmtNode converts bytes to a jmt.Node
func DeserializeJmtNode(b []byte) (jmt.Node, error) {
	pb := &Node{}
	if err := proto.Unmarshal(b, pb); err != nil {
		return nil, fmt.Errorf("could not unmarshal protobuf: %w", err)
	}

	switch n := pb.Node.(type) {
	case *Node_InternalNode:
		return deserializeInternalNode(n.InternalNode)
	case *Node_LeafNode:
		return deserializeLeafNode(n.LeafNode)
	default:
		return nil, fmt.Errorf("unknown protobuf node type: %T", pb.Node)
	}
}

func serializeInternalNode(node *jmt.InternalNode) (*InternalNode, error) {
	var bitmap uint32
	children := make([]*InternalNodeChild, 0, 16)

	for i, child := range node.Children {
		if child != nil {
			bitmap |= (1 << i) // set bit for this nibble position
			children = append(children, &InternalNodeChild{
				// zero-initialize protobuf built-ins
				protoimpl.MessageState{},
				0,
				nil,
				// fields
				uint64(child.Version),
				child.Digest[:],
				child.IsLeaf,
			})
		}
		// nil children are omitted from the array
	}

	return &InternalNode{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		bitmap,
		children,
	}, nil
}

func deserializeInternalNode(pb *InternalNode) (*jmt.InternalNode, error) {
	var children [16]*jmt.Child
	childIndex := 0

	// Use bitmap to determine which nibble positions have children
	for nibble := range children {
		if (pb.ChildrenBitmap & (1 << nibble)) != 0 {
			// This nibble position has a child
			if childIndex >= len(pb.Children) {
				return nil, fmt.Errorf("bitmap indicates child at nibble %d but children array is too short", nibble)
			}

			pbChild := pb.Children[childIndex]
			var digest jmt.Digest
			if len(pbChild.Digest) != len(digest) {
				return nil, fmt.Errorf("child digest must be %d bytes, got %d", len(digest), len(pbChild.Digest))
			}
			copy(digest[:], pbChild.Digest)
			children[nibble] = &jmt.Child{
				jmt.Version(pbChild.Version),
				digest,
				pbChild.IsLeaf,
			}
			childIndex++
		}
		// nibble positions with unset bitmap bits remain nil
	}

	if childIndex != len(pb.Children) {
		return nil, fmt.Errorf("bitmap indicates %d children but array has %d", childIndex, len(pb.Children))
	}

	return &jmt.InternalNode{children}, nil
}

func serializeLeafNode(node *jmt.LeafNode) (*LeafNode, error) {
	return &LeafNode{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		node.KeyDigest[:],
		node.Key,
		node.ValueDigest[:],
		node.Value,
	}, nil
}

func deserializeLeafNode(pb *LeafNode) (*jmt.LeafNode, error) {
	var keyDigest, valueDigest jmt.Digest
	expectedDigestLen := len(jmt.Digest{})
	if len(pb.KeyDigest) != expectedDigestLen {
		return nil, fmt.Errorf("key digest must be %d bytes, got %d", expectedDigestLen, len(pb.KeyDigest))
	}
	if len(pb.ValueDigest) != expectedDigestLen {
		return nil, fmt.Errorf("value digest must be %d bytes, got %d", expectedDigestLen, len(pb.ValueDigest))
	}
	copy(keyDigest[:], pb.KeyDigest)
	copy(valueDigest[:], pb.ValueDigest)

	return &jmt.LeafNode{
		keyDigest,
		pb.Key,
		valueDigest,
		pb.Value,
	}, nil
}

// Important: Serialization must preserve ordering of NodeKeys.
func AppendSerializeNodeKey(buffer []byte, nodeKey jmt.NodeKey) []byte {
	buffer = binary.BigEndian.AppendUint64(buffer, nodeKey.Version)
	buffer = append(buffer, byte(nodeKey.NibblePath.NumNibbles()))
	buffer = append(buffer, nodeKey.NibblePath.Bytes()...)
	return buffer
}

func DeserializeNodeKey(enc []byte) (jmt.NodeKey, error) {
	if len(enc) < 8 {
		return jmt.NodeKey{}, fmt.Errorf("encoding has no version")
	}
	version := binary.BigEndian.Uint64(enc[:8])
	enc = enc[8:]
	if len(enc) == 0 {
		return jmt.NodeKey{}, fmt.Errorf("encoding has no num nibbles")
	}
	numNibbles := int(enc[0])
	enc = enc[1:]
	if len(enc) != (numNibbles+1)/2 {
		return jmt.NodeKey{}, fmt.Errorf("encoding has less bytes than expected for nibbles")
	}
	nibblePath, ok := jmt.NewNibblePath(numNibbles, enc[:])
	if !ok {
		return jmt.NodeKey{}, fmt.Errorf("encoding has invalid nibble path")
	}
	return jmt.NodeKey{version, nibblePath}, nil
}

// Important: Serialization must preserve ordering of StaleNodes.
func AppendSerializeStaleNode(buffer []byte, staleNode jmt.StaleNode) []byte {
	buffer = binary.BigEndian.AppendUint64(buffer, staleNode.StaleSinceVersion)
	buffer = AppendSerializeNodeKey(buffer, staleNode.NodeKey)
	return buffer
}

func DeserializeStaleNode(enc []byte) (jmt.StaleNode, error) {
	if len(enc) < 8 {
		return jmt.StaleNode{}, fmt.Errorf("encoding too short")
	}
	version := binary.BigEndian.Uint64(enc[:8])
	enc = enc[8:]
	nodeKey, err := DeserializeNodeKey(enc)
	if err != nil {
		return jmt.StaleNode{}, fmt.Errorf("error decoding node key: %w", err)
	}
	return jmt.StaleNode{version, nodeKey}, nil
}
