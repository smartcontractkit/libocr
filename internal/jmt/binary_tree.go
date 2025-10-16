package jmt

import (
	"crypto/sha256"
	"fmt"
)

const (
	MaxBinaryTreeDepth = MaxHexTreeDepth * 4
	MaxProofLength     = MaxBinaryTreeDepth
)

var SparseMerklePlaceholderDigest = Digest([]byte("SPARSE_MERKLE_PLACEHOLDER_HASH__"))
var LeafDomainSeparator = []byte("JMT::LeafNode")
var InternalDomainSeparator = []byte("JMT::IntrnalNode") // sic

func digestInternalBinary(leftChildDigest Digest, rightChildDigest Digest) Digest {
	hash := sha256.New()
	hash.Write(InternalDomainSeparator)
	hash.Write(leftChildDigest[:])
	hash.Write(rightChildDigest[:])
	return Digest(hash.Sum(nil))
}

func digestLeafBinary(keyDigest Digest, valueDigest Digest) Digest {
	hash := sha256.New()
	hash.Write(LeafDomainSeparator)
	hash.Write(keyDigest[:])
	hash.Write(valueDigest[:])
	return Digest(hash.Sum(nil))
}

type sparseMerkleNode struct {
	isLeaf bool
	digest Digest
}

func (n sparseMerkleNode) String() string {
	if n.isLeaf {
		return fmt.Sprintf("Leaf(%x)", n.digest)
	}
	return fmt.Sprintf("Internal(%x)", n.digest)
}

func sparseDigestInternalBinary(leftNode sparseMerkleNode, rightNode sparseMerkleNode) sparseMerkleNode {
	if leftNode.digest == SparseMerklePlaceholderDigest && rightNode.digest == SparseMerklePlaceholderDigest {
		return sparseMerkleNode{false, SparseMerklePlaceholderDigest}
	} else if leftNode.isLeaf && rightNode.digest == SparseMerklePlaceholderDigest {
		return leftNode
	} else if leftNode.digest == SparseMerklePlaceholderDigest && rightNode.isLeaf {
		return rightNode
	} else {
		return sparseMerkleNode{false, digestInternalBinary(leftNode.digest, rightNode.digest)}
	}
}

func evolveLayer(bottomLayer []sparseMerkleNode) []sparseMerkleNode {
	if len(bottomLayer)%2 != 0 {
		panic("")
	}

	nextLayer := make([]sparseMerkleNode, 0, len(bottomLayer)/2)

	for i, j := 0, 1; j < len(bottomLayer); i, j = i+2, j+2 {
		nextLayer = append(nextLayer, sparseDigestInternalBinary(bottomLayer[i], bottomLayer[j]))
	}
	return nextLayer
}

func sparseMerkleTreeDigest(bottomLayer []sparseMerkleNode) Digest {
	if len(bottomLayer) == 0 {
		return SparseMerklePlaceholderDigest
	}
	if len(bottomLayer) == 1 {
		return bottomLayer[0].digest
	}
	return sparseMerkleTreeDigest(evolveLayer(bottomLayer))
}

// Calling with a compactable i will return invalid proofs, because the node is really not part of the final tree.
func sparseMerkleProof(bottomLayer []sparseMerkleNode, i int) []Digest {
	// find the bottommost layer where the node at index i remains after compactions
	for len(bottomLayer) > 1 {
		bottomLayerCandidate := evolveLayer(bottomLayer)

		if len(bottomLayerCandidate) > i/2 && bottomLayerCandidate[i/2] == bottomLayer[i] {
			bottomLayer = bottomLayerCandidate
			i = i / 2
		} else {
			break
		}
	}

	if len(bottomLayer) <= 1 {

		return nil
	}

	mid := len(bottomLayer) / 2

	var (
		sibling            Digest
		descendantSiblings []Digest
	)
	if i < mid {
		sibling = sparseMerkleTreeDigest(bottomLayer[mid:])
		descendantSiblings = sparseMerkleProof(bottomLayer[:mid], i)
	} else { // i >= mid
		sibling = sparseMerkleTreeDigest(bottomLayer[:mid])
		descendantSiblings = sparseMerkleProof(bottomLayer[mid:], i-mid)
	}

	return append(descendantSiblings, sibling)
}
