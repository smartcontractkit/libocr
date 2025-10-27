// In-memory Merkle Tree which implements an authenticated list.
package mt

import (
	"crypto/sha256"
	"fmt"
	"iter"
)

type Digest = [sha256.Size]byte

var MerklePlaceholderDigest = Digest([]byte("MERKLE_PLACEHOLDER_HASH_________"))
var LeafSeparator = []byte("MT::LeafNode")
var InternalSeparator = []byte("MT::InternalNode")

func digestInternal(leftChildDigest Digest, rightChildDigest Digest) Digest {
	hash := sha256.New()
	hash.Write(InternalSeparator)
	hash.Write(leftChildDigest[:])
	hash.Write(rightChildDigest[:])
	return Digest(hash.Sum(nil))
}

func digestLeaf(preimage []byte) Digest {
	hash := sha256.New()
	hash.Write(LeafSeparator)
	hash.Write(preimage)
	return Digest(hash.Sum(nil))
}

func buildTreeLevels(leafPreimages [][]byte) iter.Seq[[]Digest] {
	return func(yield func([]Digest) bool) {
		leafCount := len(leafPreimages)
		if leafCount == 0 {
			if !yield([]Digest{MerklePlaceholderDigest}) {
				return
			}
		}

		// Start with the leaf digests
		currentLayer := make([]Digest, leafCount)
		for i, leafPreimage := range leafPreimages {
			currentLayer[i] = digestLeaf(leafPreimage)
		}
		if !yield(currentLayer) {
			return
		}

		// Build the tree upwards, padding with placeholders when odd number of nodes
		for len(currentLayer) > 1 {
			nextLayerSize := (len(currentLayer) + 1) / 2 // Ceiling division
			nextLayer := make([]Digest, 0, nextLayerSize)

			for i := 0; i < len(currentLayer); i += 2 {
				leftDigest := currentLayer[i]
				rightDigest := MerklePlaceholderDigest
				if i+1 < len(currentLayer) {
					rightDigest = currentLayer[i+1]
				}
				nextLayer = append(nextLayer, digestInternal(leftDigest, rightDigest))
			}
			currentLayer = nextLayer
			if !yield(currentLayer) {
				return
			}
		}
	}
}

// Root computes the Merkle tree root from leaf preimages.
func Root(leafPreimages [][]byte) Digest {
	var rootDigest Digest
	for treeLevel := range buildTreeLevels(leafPreimages) {
		if len(treeLevel) == 1 {
			rootDigest = treeLevel[0]
		}
	}
	return rootDigest
}

// Prove generates a Merkle inclusion proof that the leafPreimage at index is
// included in the tree rooted at Root(leafPreimages).
func Prove(leafPreimages [][]byte, index uint64) ([]Digest, error) {
	leafCount := len(leafPreimages)
	if leafCount == 0 {
		return nil, fmt.Errorf("cannot prove inclusion in empty tree")
	}
	if index >= uint64(leafCount) {
		return nil, fmt.Errorf("index %d is out of bounds for %d leaves", index, leafCount)
	}

	var proof []Digest
	currentIndex := index

	for treeLevel := range buildTreeLevels(leafPreimages) {
		if len(treeLevel) <= 1 {
			break
		}
		siblingIndex := currentIndex ^ 1
		siblingDigest := MerklePlaceholderDigest
		if siblingIndex < uint64(len(treeLevel)) {
			siblingDigest = treeLevel[siblingIndex]
		}
		proof = append(proof, siblingDigest)
		currentIndex /= 2
	}

	return proof, nil
}

// Verify verifies that leafPreimage is preimage of the index-th leaf in the
// Merkle tree rooted at expectedRootDigest.
func Verify(expectedRootDigest Digest, index uint64, leafPreimage []byte, proof []Digest) error {
	currentDigest := digestLeaf(leafPreimage)
	currentIndex := index

	for _, siblingDigest := range proof {
		if currentIndex%2 == 0 {
			// Current node is left child, sibling is right
			currentDigest = digestInternal(currentDigest, siblingDigest)
		} else {
			// Current node is right child, sibling is left
			currentDigest = digestInternal(siblingDigest, currentDigest)
		}
		currentIndex /= 2
	}

	if currentDigest != expectedRootDigest {
		return fmt.Errorf("computed root digest mismatch: computed %x, expected %x", currentDigest, expectedRootDigest)
	}

	return nil
}
