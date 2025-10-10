package jmt

import (
	"crypto/sha256"
	"fmt"
)

const MaxHexTreeDepth = len(Digest{}) * 2

// go-sumtype:decl Node

type Node interface {
	isNode()
}

type InternalNode struct {
	Children [16]*Child
}

func (n *InternalNode) isNode() {}

func (n *InternalNode) String() string {
	return fmt.Sprintf("InternalNode{children: %v}", n.Children)
}

type Child struct {
	Version Version
	Digest  Digest
	IsLeaf  bool
}

func (c *Child) String() string {
	return fmt.Sprintf("Child{version: %d, digest: %x, isLeaf: %t}", c.Version, c.Digest, c.IsLeaf)
}

type LeafNode struct {
	KeyDigest   Digest
	Key         []byte
	ValueDigest Digest
	Value       []byte
}

func (n *LeafNode) isNode() {}

func (n *LeafNode) String() string {
	return fmt.Sprintf("LeafNode{keyDigest: %x, valueDigest: %x, value: %x}", n.KeyDigest, n.ValueDigest, n.Value)
}

type NodeKey struct {
	Version    Version
	NibblePath NibblePath
}

func (nk NodeKey) Equal(nk2 NodeKey) bool {
	return nk.Version == nk2.Version && nk.NibblePath.Equal(nk2.NibblePath)
}

func digestNode(node Node) Digest {

	if node == nil {
		return SparseMerklePlaceholderDigest
	}

	switch n := node.(type) {
	case *InternalNode:
		return digestInternalNode(n)
	case *LeafNode:
		return digestLeafNode(n)
	default:
		panic("")
	}
}

func hexInternalNodeToBinaryBottomLayer(node *InternalNode) []sparseMerkleNode {
	bottomLayer := make([]sparseMerkleNode, len(node.Children))
	for i, child := range node.Children {
		if child != nil {
			bottomLayer[i] = sparseMerkleNode{
				child.IsLeaf,
				child.Digest,
			}
		} else {
			bottomLayer[i] = sparseMerkleNode{
				false,
				SparseMerklePlaceholderDigest,
			}
		}
	}
	return bottomLayer
}

func digestInternalNode(node *InternalNode) Digest {
	return sparseMerkleTreeDigest(hexInternalNodeToBinaryBottomLayer(node))
}

func proveInternalNodeChild(node *InternalNode, childIndex int) []Digest {
	bottomLayer := hexInternalNodeToBinaryBottomLayer(node)
	return sparseMerkleProof(bottomLayer, childIndex)
}

func digestLeafNode(node *LeafNode) Digest {
	return digestLeafBinary(node.KeyDigest, node.ValueDigest)
}

func DigestKey(key []byte) Digest {
	hash := sha256.New()
	hash.Write(key)
	return Digest(hash.Sum(nil))
}

func DigestValue(value []byte) Digest {
	hash := sha256.New()
	hash.Write(value)
	return Digest(hash.Sum(nil))
}
