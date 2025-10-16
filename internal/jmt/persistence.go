package jmt

// When starting from scratch, you can safely use version 0 as the "old version"
// contains an empty tree.
type Version = uint64

type RootWriter interface {
	WriteRoot(version Version, nodeKey NodeKey) error
}

type RootReader interface {
	ReadRoot(version Version) (nodeKey NodeKey, err error)
}

type RootReadWriter interface {
	RootReader
	RootWriter
}

type NodeWriter interface {
	// Implementers should check whether node is nil first. nil indicates an
	// empty tree, and should effectively cause a deletion of the node key.
	// Implementers should make sure to not be affected by later mutations to
	// Node.
	WriteNode(nodeKey NodeKey, nodeOrNil Node) error
}

type NodeReader interface {
	ReadNode(nodeKey NodeKey) (nodeOrNil Node, err error)
}

type NodeReadWriter interface {
	NodeReader
	NodeWriter
}

type StaleNode struct {
	StaleSinceVersion Version
	NodeKey           NodeKey
}

type StaleNodeWriter interface {
	WriteStaleNode(staleNode StaleNode) error
}
