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

func sparseMerkleTreeDigest(bottomLayer []sparseMerkleNode) Digest {
	return sparseMerkleTreeRoot(bottomLayer).digest
}

func sparseMerkleTreeRoot(bottomLayer []sparseMerkleNode) sparseMerkleNode {
	if len(bottomLayer) == 0 {
		return sparseMerkleNode{false, SparseMerklePlaceholderDigest}
	}
	if !isPowerOfTwo(len(bottomLayer)) {
		panic("bottomLayer length must be a power of two if not zero")
	}
	if len(bottomLayer) == 1 {
		return bottomLayer[0]
	}
	return sparseDigestInternalBinary(
		sparseMerkleTreeRoot(bottomLayer[:len(bottomLayer)/2]),
		sparseMerkleTreeRoot(bottomLayer[len(bottomLayer)/2:]),
	)
}

func sparseMerkleIsLeafAndProof(bottomLayer []sparseMerkleNode, i int) (singleLeaf bool, proof []Digest) {
	if !isPowerOfTwo(len(bottomLayer)) {
		panic("bottomLayer length must be a power of two")
	}
	if i < 0 || i >= len(bottomLayer) {
		panic("i is out of bounds")
	}
	if bottomLayer[i].digest == SparseMerklePlaceholderDigest {
		panic("bottomLayer[i] is a placeholder")
	}
	if len(bottomLayer) == 1 {
		return bottomLayer[0].isLeaf, nil
	}

	mid := len(bottomLayer) / 2
	var (
		sibling             Digest
		subtreeSiblings     []Digest
		subtreeIsSingleLeaf bool
	)
	if i < mid {
		sibling = sparseMerkleTreeDigest(bottomLayer[mid:])
		subtreeIsSingleLeaf, subtreeSiblings = sparseMerkleIsLeafAndProof(bottomLayer[:mid], i)
	} else { // i >= mid
		sibling = sparseMerkleTreeDigest(bottomLayer[:mid])
		subtreeIsSingleLeaf, subtreeSiblings = sparseMerkleIsLeafAndProof(bottomLayer[mid:], i-mid)
	}

	rootIsSingleLeaf := subtreeIsSingleLeaf && sibling == SparseMerklePlaceholderDigest
	if rootIsSingleLeaf {
		return true, nil
	}
	return false, append(subtreeSiblings, sibling)
}

func sparseMerkleProof(bottomLayer []sparseMerkleNode, i int) []Digest {
	_, proof := sparseMerkleIsLeafAndProof(bottomLayer, i)
	return proof
}

func isPowerOfTwo(num int) bool {
	return num > 0 && (num&(num-1)) == 0
}
