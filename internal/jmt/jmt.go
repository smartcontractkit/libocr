package jmt

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"slices"
	"sort"
)

type KeyValue struct {
	// Key length must be greater than 0.
	Key []byte
	// Value of nil is permitted and indicates a desire to delete Key.
	Value []byte
}

type digestedKeyValue struct {
	key         []byte
	keyDigest   Digest
	value       []byte
	valueDigest Digest
}

func (dkv digestedKeyValue) String() string {
	return fmt.Sprintf("%x->%x", dkv.keyDigest, dkv.value)
}

// BatchUpdate performs the updates indicated by keyValueUpdates on top of the
// tree at oldVersion. If oldVersion coincides with newVersion, the updates are
// performed in place and no stale nodes are written. If oldVersion is less than
// newVersion, we perform copy-on-write, and stale nodes might be written that
// might need to be reaped in the future using ReapStaleNodes. Using an
// oldVersion that is greater than newVersion is an error.
//
// keyValueUpdates must be unique in Key.
//
// BatchUpdate returns the NodeKey of the new root node at newVersion.
func BatchUpdate(
	rootReadWriter RootReadWriter,
	nodeReadWriter NodeReadWriter,
	staleNodeWriter StaleNodeWriter,
	oldVersion Version,
	newVersion Version,
	keyValueUpdates []KeyValue,
) (NodeKey, error) {
	if oldVersion > newVersion {
		return NodeKey{}, fmt.Errorf("old version %d is greater than new version %d", oldVersion, newVersion)
	}

	oldRootNodeKey, err := rootReadWriter.ReadRoot(oldVersion)
	if err != nil {
		return NodeKey{}, fmt.Errorf("error reading root node with version %d: %w", oldVersion, err)
	}

	oldRootNode, err := nodeReadWriter.ReadNode(oldRootNodeKey)
	if err != nil {
		return NodeKey{}, fmt.Errorf("error reading root node with node key %v: %w", oldRootNodeKey, err)
	}

	digestedInserts := make([]digestedKeyValue, 0, len(keyValueUpdates))
	digestedDeletes := make([]digestedKeyValue, 0, len(keyValueUpdates))

	{

		seenDigestedKeys := make(map[Digest]struct{}, len(keyValueUpdates))

		for i, keyValue := range keyValueUpdates {
			if len(keyValue.Key) == 0 {
				return NodeKey{}, fmt.Errorf("%d-th keyValueUpdate: key is empty", i)
			}

			keyDigest := DigestKey(keyValue.Key)
			if _, ok := seenDigestedKeys[keyDigest]; ok {
				return NodeKey{}, fmt.Errorf("%d-th keyValueUpdate: duplicate key %v in keyValueUpdates", i, keyValue.Key)
			}
			seenDigestedKeys[keyDigest] = struct{}{}

			var valueDigest Digest
			if keyValue.Value != nil {
				valueDigest = DigestValue(keyValue.Value)
			}

			dkv := digestedKeyValue{
				keyValue.Key,
				keyDigest,
				keyValue.Value,
				valueDigest,
			}

			if keyValue.Value == nil {
				digestedDeletes = append(digestedDeletes, dkv)
			} else {
				digestedInserts = append(digestedInserts, dkv)
			}
		}
	}

	// in order of leaf insertion
	sort.Slice(digestedInserts, func(i, j int) bool {
		return bytes.Compare(digestedInserts[i].keyDigest[:], digestedInserts[j].keyDigest[:]) < 0
	})
	sort.Slice(digestedDeletes, func(i, j int) bool {
		return bytes.Compare(digestedDeletes[i].keyDigest[:], digestedDeletes[j].keyDigest[:]) < 0
	})

	newRootNodeKey, _, err := batchUpdate(
		nodeReadWriter,
		staleNodeWriter,
		newVersion,
		oldRootNodeKey,
		oldRootNode,
		digestedInserts,
		digestedDeletes,
	)
	if err != nil {
		return NodeKey{}, fmt.Errorf("error performing batch update %v: %w", keyValueUpdates, err)
	}

	err = rootReadWriter.WriteRoot(newVersion, newRootNodeKey)
	if err != nil {
		return NodeKey{}, fmt.Errorf("error writing root node with node key %v: %w", newRootNodeKey, err)
	}

	return newRootNodeKey, nil
}

func filterSortedDigestedKeyValues(
	sortedDigestedKeyValues []digestedKeyValue,
	digestPrefix NibblePath,
) []digestedKeyValue {
	matchIndex := digestPrefix.NumNibbles() - 1
	if matchIndex < 0 || matchIndex >= len(Digest{})*2 {
		panic(fmt.Errorf("match index %v out of bounds", matchIndex))
	}
	matchNibble := digestPrefix.Get(matchIndex)

	// assume filterSortedDigestedKeyValues has been recursively called with digestPrefix[:-1] to produce our input sortedDigestedKeyValues
	firstMatchingIndex := sort.Search(len(sortedDigestedKeyValues), func(i int) bool {
		return NibblePathFromDigest(sortedDigestedKeyValues[i].keyDigest).Get(matchIndex) >= matchNibble
	})
	firstNonMatchingIndex := sort.Search(len(sortedDigestedKeyValues), func(i int) bool {
		return NibblePathFromDigest(sortedDigestedKeyValues[i].keyDigest).Get(matchIndex) > matchNibble
	})
	return sortedDigestedKeyValues[firstMatchingIndex:firstNonMatchingIndex]
}

func batchUpdate(
	nodeReadWriter NodeReadWriter,
	staleNodeWriter StaleNodeWriter,
	version Version, // version of new or altered nodes
	rootNodeKey NodeKey,
	rootNode Node,
	sortedDigestedInserts []digestedKeyValue,
	sortedDigestedDeletes []digestedKeyValue,
) (newRootNodeKey NodeKey, newRootNode Node, err error) {
	oldNodeKey := rootNodeKey
	oldNodeWasNil := rootNode == nil
	replaceRootNode := func(newNodeKey NodeKey, newNode Node) error {
		var err error
		if !oldNodeWasNil {
			// there is something to delete
			if oldNodeKey.Version < version {
				err = errors.Join(err, staleNodeWriter.WriteStaleNode(StaleNode{version, oldNodeKey}))
			} else if oldNodeKey.Version == version {
				err = errors.Join(err, nodeReadWriter.WriteNode(oldNodeKey, nil))
			} else {
				return fmt.Errorf("assumption violation: old node key version %d is greater than new node key version %d", oldNodeKey.Version, version)
			}
		}
		err = errors.Join(err, nodeReadWriter.WriteNode(newNodeKey, newNode))
		return err
	}

	if len(sortedDigestedInserts) == 0 && len(sortedDigestedDeletes) == 0 {
		return rootNodeKey, rootNode, nil
	}

	// len(sortedDigestedInserts) >= 1 or len(sortedDigestedDeletes) >= 1 beyond this point

	if rootNode == nil {
		// we can ignore deletions, because we can't delete nodes that don't exist
		sortedDigestedDeletes = nil

		if len(sortedDigestedInserts) == 0 {
			// nothing to do, nil -> nil
			return rootNodeKey, rootNode, nil
		}

		if len(sortedDigestedInserts) == 1 {
			// nil -> leaf
			dkv := sortedDigestedInserts[0]
			leafNode := &LeafNode{
				dkv.keyDigest,
				dkv.key,
				dkv.valueDigest,
				dkv.value,
			}
			leafNodeKey := NodeKey{
				version,
				rootNodeKey.NibblePath,
			}
			err := replaceRootNode(leafNodeKey, leafNode)
			if err != nil {
				return NodeKey{}, nil, err
			}
			return leafNodeKey, leafNode, nil
		}

		// nil -> internal
		// more than one insertions, definitely need an internal node

		rootNode = &InternalNode{} // pretend it's an empty internal node
	}

	if rootLeafNode, ok := rootNode.(*LeafNode); ok {
		leafNodeKey := NodeKey{
			version,
			rootNodeKey.NibblePath,
		}

		if len(sortedDigestedInserts) == 1 {
			dkv := sortedDigestedInserts[0]
			if dkv.keyDigest == rootLeafNode.KeyDigest {
				// leaf -> leaf
				leafNode := &LeafNode{
					dkv.keyDigest,
					dkv.key,
					dkv.valueDigest,
					dkv.value,
				}
				err := replaceRootNode(leafNodeKey, leafNode)
				if err != nil {
					return NodeKey{}, nil, err
				}
				return leafNodeKey, leafNode, nil
			}
		}

		// could contain many spurious deletes, so search for our leaf key digest
		_, deleteRootLeafNode := sort.Find(len(sortedDigestedDeletes), func(i int) int {
			return bytes.Compare(rootLeafNode.KeyDigest[:], sortedDigestedDeletes[i].keyDigest[:])
		})

		if len(sortedDigestedInserts) == 0 {
			if deleteRootLeafNode {
				// leaf -> nil
				err := replaceRootNode(leafNodeKey, nil)
				if err != nil {
					return NodeKey{}, nil, err
				}
				return NodeKey{}, nil, nil
			} else {
				// noop
				return rootNodeKey, rootNode, nil
			}
		}

		// leaf -> internal
		insertionPoint := sort.Search(len(sortedDigestedInserts), func(i int) bool {
			return bytes.Compare(sortedDigestedInserts[i].keyDigest[:], rootLeafNode.KeyDigest[:]) >= 0
		})
		if !deleteRootLeafNode && (insertionPoint >= len(sortedDigestedInserts) || sortedDigestedInserts[insertionPoint].keyDigest != rootLeafNode.KeyDigest) {
			// we didn't already have an update in mind for this leaf key digest
			// carry it by inserting (in key digest sorted order) to sortedDigestedInserts

			carryDigestedKeyValue := digestedKeyValue{
				rootLeafNode.Key,
				rootLeafNode.KeyDigest,
				rootLeafNode.Value,
				rootLeafNode.ValueDigest,
			}

			sortedDigestedInserts = slices.Clone(sortedDigestedInserts)
			sortedDigestedInserts = slices.Insert(sortedDigestedInserts, insertionPoint, carryDigestedKeyValue)
		}
		rootNode = &InternalNode{} // pretend it's an empty internal node
	}

	if rootInternalNode, ok := rootNode.(*InternalNode); ok {
		newChildrenNonNilCount := 0
		var (
			newChildrenLastLeafNodeKey NodeKey
			newChildrenLastLeaf        *LeafNode
		)

		internalNode := &InternalNode{}

		anyChildChanged := false

		for nibble := range 16 {
			child := rootInternalNode.Children[nibble]
			childNibblePath := rootNodeKey.NibblePath.Append(byte(nibble))

			var (
				childNodeKey NodeKey
				childNode    Node
			)
			if child == nil {

				childNodeKey = NodeKey{version, childNibblePath}
			} else {
				childNodeKey = NodeKey{child.Version, childNibblePath}
				var err error
				childNode, err = nodeReadWriter.ReadNode(childNodeKey)
				if err != nil {
					return NodeKey{}, nil, err
				}
			}

			sortedDigestedInsertsBelowNibble := filterSortedDigestedKeyValues(sortedDigestedInserts, childNibblePath)
			sortedDigestedDeletesBelowNibble := filterSortedDigestedKeyValues(sortedDigestedDeletes, childNibblePath)
			newChildNodeKey, newChildNode, err := batchUpdate(
				nodeReadWriter,
				staleNodeWriter,
				version,
				childNodeKey,
				childNode,
				sortedDigestedInsertsBelowNibble,
				sortedDigestedDeletesBelowNibble,
			)

			if err != nil {
				return NodeKey{}, nil, err
			}

			if childNode != newChildNode {

				anyChildChanged = true
			}

			if newChildNode == nil {
				internalNode.Children[nibble] = nil
				// deletion has been taken care of by batchUpdate with the child node as root
			} else {
				newChildrenNonNilCount++
				isLeaf := false
				if newChildNodeLeaf, ok := newChildNode.(*LeafNode); ok {
					isLeaf = true

					newChildrenLastLeafNodeKey = newChildNodeKey
					newChildrenLastLeaf = newChildNodeLeaf
				}

				internalNode.Children[nibble] = &Child{
					newChildNodeKey.Version,
					digestNode(newChildNode),
					isLeaf,
				}
			}
		}

		if !anyChildChanged {
			// noop
			return rootNodeKey, rootNode, nil
		}

		newRootNodeKey := NodeKey{
			version,
			rootNodeKey.NibblePath,
		}

		// if no children anymore, delete this node

		if newChildrenNonNilCount == 0 {
			err := replaceRootNode(newRootNodeKey, nil)
			if err != nil {
				return NodeKey{}, nil, err
			}
			return NodeKey{}, nil, nil
		}

		// if only child is a leaf, make the root node that leaf, while keeping the same nibble path

		if newChildrenNonNilCount == 1 && newChildrenLastLeaf != nil {
			// the leaf is already written more deeply, so we need to delete it

			if newChildrenLastLeafNodeKey.Version == version {
				// we can remove it in place
				err := nodeReadWriter.WriteNode(newChildrenLastLeafNodeKey, nil)
				if err != nil {
					return NodeKey{}, nil, err
				}
			} else {
				// it's an older node so we mark it as stale
				err := staleNodeWriter.WriteStaleNode(StaleNode{version, newChildrenLastLeafNodeKey})
				if err != nil {
					return NodeKey{}, nil, err
				}
			}

			err := replaceRootNode(newRootNodeKey, newChildrenLastLeaf)
			if err != nil {
				return NodeKey{}, nil, err
			}
			return newRootNodeKey, newChildrenLastLeaf, nil
		}

		err := replaceRootNode(newRootNodeKey, internalNode)
		if err != nil {
			return NodeKey{}, nil, err
		}
		return newRootNodeKey, internalNode, nil
	}

	panic("unreachable")
}

// ProveSubrange returns the bounding leafs required to verify inclusion of the
// vector of all the key-values of which the digests fall in the subrange
// [startIndex, endInclIndex], using VerifySubrange.
func ProveSubrange(
	rootReader RootReader,
	nodeReader NodeReader,
	version Version,
	startIndex Digest,
	endInclIndex Digest,
) ([]BoundingLeaf, error) {

	var boundingLeaves []BoundingLeaf

	left, _, err := ReadRangeAscOrDesc(rootReader, nodeReader, version, MinDigest, startIndex, math.MaxInt, 1, ReadRangeOrderDesc)
	if err != nil {
		return nil, fmt.Errorf("finding left bounding leaf for %x failed: %w", startIndex, err)
	}
	if len(left) > 0 {
		leftLeafKeyDigest := DigestKey(left[0].Key)
		leftLeafProof, err := proveInclusionDigested(rootReader, nodeReader, version, leftLeafKeyDigest)
		if err != nil {
			return nil, fmt.Errorf("proving inclusion of left bounding leaf %x failed: %w", leftLeafKeyDigest, err)
		}

		var leftSiblings []Digest
		for i, p := range leftLeafProof {
			if bitGet(leftLeafKeyDigest, len(leftLeafProof)-1-i) {
				// only keep left proof siblings
				leftSiblings = append(leftSiblings, p)
			}
		}
		boundingLeaves = append(boundingLeaves, BoundingLeaf{
			LeafKeyAndValueDigests{
				leftLeafKeyDigest,
				DigestValue(left[0].Value),
			},
			leftSiblings,
		})
	}

	right, _, err := ReadRangeAscOrDesc(rootReader, nodeReader, version, endInclIndex, MaxDigest, math.MaxInt, 1, ReadRangeOrderAsc)
	if err != nil {
		return nil, fmt.Errorf("finding right bounding leaf for %x failed: %w", endInclIndex, err)
	}
	if len(right) > 0 {
		rightLeafKeyDigest := DigestKey(right[0].Key)
		rightLeafProof, err := proveInclusionDigested(rootReader, nodeReader, version, rightLeafKeyDigest)
		if err != nil {
			return nil, fmt.Errorf("proving inclusion of right bounding leaf %x failed: %w", rightLeafKeyDigest, err)
		}

		var rightSiblings []Digest
		for i, p := range rightLeafProof {
			if !bitGet(rightLeafKeyDigest, len(rightLeafProof)-1-i) {
				// only keep right proof siblings
				rightSiblings = append(rightSiblings, p)
			}
		}
		boundingLeaves = append(boundingLeaves, BoundingLeaf{
			LeafKeyAndValueDigests{
				rightLeafKeyDigest,
				DigestValue(right[0].Value),
			},
			rightSiblings,
		})
	}

	return boundingLeaves, nil
}

type LeafKeyAndValueDigests struct {
	KeyDigest   Digest
	ValueDigest Digest
}

// ProveInclusion returns the proof siblings for the inclusion proof of the
// undigested key in the tree at version.
func ProveInclusion(
	rootReader RootReader,
	nodeReader NodeReader,
	version Version,
	key []byte,
) ([]Digest, error) {
	return proveInclusionDigested(rootReader, nodeReader, version, DigestKey(key))
}

func proveInclusionDigested(
	rootReader RootReader,
	nodeReader NodeReader,
	version Version,
	digestedKey Digest,
) ([]Digest, error) {
	var siblings []Digest

	rootNodeKey, err := rootReader.ReadRoot(version)
	if err != nil {
		return nil, fmt.Errorf("error reading root node with version %d: %w", version, err)
	}

	rootNode, err := nodeReader.ReadNode(rootNodeKey)
	if err != nil {
		return nil, fmt.Errorf("error reading root node with node key %v: %w", rootNodeKey, err)
	}

	if rootNode == nil {
		return nil, fmt.Errorf("asked to inclusion prove in nil tree")
	}

	nibblePath := NibblePathFromDigest(digestedKey)

	for i := 0; i < nibblePath.NumNibbles() && rootNode != nil; i++ {
		nibble := nibblePath.Get(i)
		switch n := rootNode.(type) {
		case *InternalNode:
			child := n.Children[nibble]
			if child == nil {
				return nil, fmt.Errorf("child node with nibble %v not found", nibble)
			}

			siblings = append(proveInternalNodeChild(n, int(nibble)), siblings...)

			rootNodeKey = NodeKey{
				child.Version,
				rootNodeKey.NibblePath.Append(nibble),
			}
			rootNode, err = nodeReader.ReadNode(rootNodeKey)
			if err != nil {
				return nil, fmt.Errorf("error reading child node with node key %v: %w", rootNodeKey, err)
			}
		case *LeafNode:
			if n.KeyDigest == digestedKey {
				break
			}
			return nil, fmt.Errorf("leaf node with key digest %x not found", digestedKey)
		}
	}

	return siblings, nil
}

// Read returns the undigested value for the undigested key in the tree at
// version. If the key does not exist, nil is returned, and no error.
func Read(
	rootReader RootReader,
	nodeReader NodeReader,
	version Version,
	key []byte,
) ([]byte, error) {
	rootNodeKey, err := rootReader.ReadRoot(version)
	if err != nil {
		return nil, fmt.Errorf("error reading root node with version %d: %w", version, err)
	}

	rootNode, err := nodeReader.ReadNode(rootNodeKey)
	if err != nil {
		return nil, fmt.Errorf("error reading root node with node key %v: %w", rootNodeKey, err)
	}

	if rootNode == nil {
		return nil, nil
	}

	nibblePath := NibblePathFromDigest(DigestKey(key))

	for i := 0; i < nibblePath.NumNibbles(); i++ {
		nibble := nibblePath.Get(i)

		if n, ok := rootNode.(*InternalNode); ok {
			child := n.Children[nibble]
			if child == nil {
				return nil, nil
			}
			childNodeKey := NodeKey{
				child.Version,
				rootNodeKey.NibblePath.Append(nibble),
			}
			childNode, err := nodeReader.ReadNode(childNodeKey)
			if err != nil {
				return nil, fmt.Errorf("error reading child node with node key %v: %w", childNodeKey, err)
			}
			if childNode == nil {
				return nil, fmt.Errorf("child node with node key %v unexpectedly nil", childNodeKey)
			}
			rootNode = childNode
			rootNodeKey = childNodeKey
		} else {
			break
		}
	}

	if n, ok := rootNode.(*LeafNode); ok {
		if n.KeyDigest == DigestKey(key) {
			return n.Value, nil
		} else {
			return nil, nil
		}
	}
	panic("unreachable")
}

// ReadRootDigest returns the digest of the root node of the tree at version. If
// digest is SparseMerklePlaceholderDigest, the tree is empty. The caller must
// be careful to only call this method for versions that exist and roots have
// indeed been written. Versions that are not written are implied to represent
// an empty tree.
func ReadRootDigest(
	rootReader RootReader,
	nodeReader NodeReader,
	version Version,
) (Digest, error) {
	rootNodeKey, err := rootReader.ReadRoot(version)
	if err != nil {
		return Digest{}, fmt.Errorf("error reading root node for version %v: %w", version, err)
	}

	rootNode, err := nodeReader.ReadNode(rootNodeKey)
	if err != nil {
		return Digest{}, fmt.Errorf("error reading node for version %v and node key %v: %w", version, rootNodeKey, err)
	}

	return digestNode(rootNode), nil
}

// ReadRange returns all key-value pairs in the range [minKeyDigest,
// maxKeyDigest] in the tree at version. This is useful for state
// synchronization where a verifier expects to build the tree from scratch in
// order from the leftmost to the rightmost leaf. We return the key-value pairs
// instead of the digests to give the verifier the most flexibility. If there
// were more key-value pairs in the range that we could not fit in the limits,
// we return truncated=true.
func ReadRange(
	rootReader RootReader,
	nodeReader NodeReader,
	version Version,
	minKeyDigest Digest,
	maxKeyDigest Digest,
	keysPlusValuesBytesLimit int,
	lenLimit int,
) (keyValues []KeyValue, truncated bool, err error) {
	return ReadRangeAscOrDesc(
		rootReader,
		nodeReader,
		version,
		minKeyDigest,
		maxKeyDigest,
		keysPlusValuesBytesLimit,
		lenLimit,
		ReadRangeOrderAsc,
	)
}

type ReadRangeOrder int

const (
	_ ReadRangeOrder = iota
	ReadRangeOrderAsc
	ReadRangeOrderDesc
)

// ReadRangeAscOrDesc is a variant of ReadRange that allows the caller to
// specify the order in which to read the key-value pairs.
func ReadRangeAscOrDesc(
	rootReader RootReader,
	nodeReader NodeReader,
	version Version,
	minKeyDigest Digest,
	maxKeyDigest Digest,
	keysPlusValuesBytesLimit int,
	lenLimit int,
	order ReadRangeOrder,
) (keyValues []KeyValue, truncated bool, err error) {
	rootNodeKey, err := rootReader.ReadRoot(version)
	if err != nil {
		return nil, false, fmt.Errorf("error reading root node for version %v: %w", version, err)
	}

	rootNode, err := nodeReader.ReadNode(rootNodeKey)
	if err != nil {
		return nil, false, fmt.Errorf("error reading node for node key %v: %w", rootNodeKey, err)
	}
	keyValues, _, truncated, err = readRange(
		nodeReader,
		rootNodeKey,
		rootNode,
		0,
		NibblePathFromDigest(minKeyDigest),
		NibblePathFromDigest(maxKeyDigest),
		minKeyDigest,
		maxKeyDigest,
		keysPlusValuesBytesLimit,
		lenLimit,
		order,
	)
	return keyValues, truncated, err
}

func readRange(
	nodeReader NodeReader,
	rootNodeKey NodeKey,
	rootNode Node,
	depth int,
	minKeyDigestNibblePath NibblePath,
	maxKeyDigestNibblePath NibblePath,
	minKeyDigest Digest,
	maxKeyDigest Digest,
	keysPlusValuesBytesLimit int,
	lenLimit int,
	order ReadRangeOrder,
) (keyValues []KeyValue, keysPlusValuesBytes int, truncated bool, err error) {
	if rootNode == nil {
		return nil, 0, false, nil
	}

	switch n := rootNode.(type) {
	case *InternalNode:
		var (
			initial int
			delta   int
		)

		switch order {
		case ReadRangeOrderAsc:
			initial, delta = 0, +1
		case ReadRangeOrderDesc:
			initial, delta = len(n.Children)-1, -1
		}
		for i := initial; i >= 0 && i < len(n.Children); i += delta {
			child := n.Children[i]
			if child == nil {
				continue
			}
			childNibblePath := rootNodeKey.NibblePath.Append(byte(i))
			{
				minKeyDigestNibblePathPrefix := minKeyDigestNibblePath.Prefix(depth + 1)
				maxKeyDigestNibblePathPrefix := maxKeyDigestNibblePath.Prefix(depth + 1)
				if !(childNibblePath.Compare(minKeyDigestNibblePathPrefix) >= 0 && childNibblePath.Compare(maxKeyDigestNibblePathPrefix) <= 0) {
					continue
				}
			}
			childNodeKey := NodeKey{
				child.Version,
				childNibblePath,
			}
			childNode, err := nodeReader.ReadNode(childNodeKey)
			if err != nil {
				return nil, 0, false, fmt.Errorf("error reading child node with node key %v: %w", childNodeKey, err)
			}
			if childNode == nil {
				return nil, 0, false, fmt.Errorf("child node with node key %v unexpectedly nil", childNodeKey)
			}
			moreKeyValues, moreKeysPlusValuesBytes, moreTruncated, err := readRange(
				nodeReader,
				childNodeKey,
				childNode,
				depth+1,
				minKeyDigestNibblePath,
				maxKeyDigestNibblePath,
				minKeyDigest,
				maxKeyDigest,
				keysPlusValuesBytesLimit-keysPlusValuesBytes,
				lenLimit-len(keyValues),
				order,
			)
			if err != nil {
				return nil, 0, false, err
			}
			keyValues = append(keyValues, moreKeyValues...)
			keysPlusValuesBytes += moreKeysPlusValuesBytes
			if moreTruncated {
				truncated = true
				break
			}
		}
	case *LeafNode:
		if bytes.Compare(n.KeyDigest[:], minKeyDigest[:]) >= 0 && bytes.Compare(n.KeyDigest[:], maxKeyDigest[:]) <= 0 {
			lenLimitOk := len(keyValues)+1 <= lenLimit
			keysPlusValuesBytesLimitOk := keysPlusValuesBytes+len(n.Key)+len(n.Value) <= keysPlusValuesBytesLimit
			if !(lenLimitOk && keysPlusValuesBytesLimitOk) {
				truncated = true
			} else {
				keysPlusValuesBytes += len(n.Key) + len(n.Value)
				keyValues = append(keyValues, KeyValue{n.Key, n.Value})
			}
		}
	}

	return keyValues, keysPlusValuesBytes, truncated, nil
}
