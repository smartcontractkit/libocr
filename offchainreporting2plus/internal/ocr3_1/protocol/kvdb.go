package protocol

import (
	"github.com/RoSpaceDev/libocr/internal/jmt"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

type KeyValueDatabaseReadTransaction interface {
	// The only read part of the interface that the plugin might see. The rest
	// of the methods might only be called by protocol code.
	ocr3_1types.KeyValueStateReader
	KeyValueDatabaseSemanticRead
	Discard()
}

type KeyValueDatabaseSemanticRead interface {
	// ReadHighestCommittedSeqNr returns the sequence number of which the state the transaction
	// represents. Really read from the database here, no cached values allowed.
	ReadHighestCommittedSeqNr() (uint64, error)
	ReadLowestPersistedSeqNr() (uint64, error)

	ReadAttestedStateTransitionBlock(seqNr uint64) (AttestedStateTransitionBlock, error)
	ReadAttestedStateTransitionBlocks(minSeqNr uint64, maxItems int) (blocks []AttestedStateTransitionBlock, more bool, err error)

	ReadTreeSyncStatus() (TreeSyncStatus, error)
	// ReadTreeSyncChunk retrieves a chunk of undigested key-value pairs in the
	// range [startIndex, requestEndInclIndex] of the key digest space. It
	// returns a maximally sized chunk that fully covers the range [startIndex,
	// endInclIndex], where endInclIndex <= requestEndInclIndex, such that the
	// chunk respects the protocol.MaxTreeSyncChunkKeys and
	// protocol.MaxTreeSyncChunkKeysPlusValuesLength limits. It also includes in
	// boundingLeaves the subrange proof, proving inclusion of key-values in the
	// range [startIndex, endInclIndex] without omissions.
	ReadTreeSyncChunk(
		toSeqNr uint64,
		startIndex jmt.Digest,
		requestEndInclIndex jmt.Digest,
	) (
		endInclIndex jmt.Digest,
		boundingLeaves []jmt.BoundingLeaf,
		keyValues []KeyValuePair,
		err error,
	)
	// ReadBlobPayload returns the payload of the blob if it exists in full and
	// the blob has not expired. If the blob existed at some point but has since
	// expired, it returns an error. If the blob never existed, it returns nil.
	// If only some chunks are present, it returns an error.
	ReadBlobPayload(BlobDigest) ([]byte, error)
	ReadBlobMeta(BlobDigest) (*BlobMeta, error)
	ReadBlobChunk(BlobDigest, uint64) ([]byte, error)
	ReadStaleBlobIndex(maxStaleSinceSeqNr uint64, limit int) ([]StaleBlob, error)

	jmt.RootReader
	jmt.NodeReader
}

type KeyValueDatabaseReadWriteTransaction interface {
	KeyValueDatabaseReadTransaction
	// The only write part of the interface that the plugin might see. The rest
	// of the methods might only be called by protocol code.
	ocr3_1types.KeyValueStateReadWriter
	KeyValueDatabaseSemanticWrite
	// Commit writes the new highest committed sequence number to the magic key
	// (if the transaction is _not_ unchecked) and commits the transaction to
	// the key value store, then discards the transaction.
	Commit() error
}

type VerifyAndWriteTreeSyncChunkResult int

const (
	_ VerifyAndWriteTreeSyncChunkResult = iota
	VerifyAndWriteTreeSyncChunkResultOkNeedMore
	VerifyAndWriteTreeSyncChunkResultOkComplete
	VerifyAndWriteTreeSyncChunkResultByzantine
	VerifyAndWriteTreeSyncChunkResultUnrelatedError
)

type KeyValueDatabaseSemanticWrite interface {
	// GetWriteSet returns a slice the KeyValuePair entries that
	// have been written in this transaction. If the value of a key has been
	// deleted, it the value is mapped to nil.
	GetWriteSet() ([]KeyValuePairWithDeletions, error)

	// CloseWriteSet returns the state root, writes it to the KV store
	// and closes the transaction for writing: any future attempts for Writes or Deletes
	// on this transaction will fail.
	CloseWriteSet() (StateRootDigest, error)

	// ApplyWriteSet applies the write set to the transaction and returns the
	// state root digest. Useful for reproposals and state synchronization. Only
	// works on checked transactions where the postSeqNr is specified at
	// creation.
	ApplyWriteSet(writeSet []KeyValuePairWithDeletions) (StateRootDigest, error)

	WriteAttestedStateTransitionBlock(seqNr uint64, block AttestedStateTransitionBlock) error
	DeleteAttestedStateTransitionBlocks(maxSeqNrToDelete uint64, maxItems int) (done bool, err error)

	// WriteHighestCommittedSeqNr writes the given sequence number to the magic
	// key. It is called before Commit on checked transactions.
	WriteHighestCommittedSeqNr(seqNr uint64) error
	WriteLowestPersistedSeqNr(seqNr uint64) error
	// VerifyAndWriteTreeSyncChunk first verifies that the keyValues are fully
	// and without omissions included in the key digest range of [startIndex,
	// endInclIndex]. Only after doing so, it writes all keyValues into the tree
	// and flat representation.
	VerifyAndWriteTreeSyncChunk(
		targetRootDigest StateRootDigest,
		targetSeqNr uint64,
		startIndex jmt.Digest,
		endInclIndex jmt.Digest,
		boundingLeaves []jmt.BoundingLeaf,
		keyValues []KeyValuePair,
	) (VerifyAndWriteTreeSyncChunkResult, error)

	WriteTreeSyncStatus(state TreeSyncStatus) error
	WriteBlobMeta(BlobDigest, BlobMeta) error
	DeleteBlobMeta(BlobDigest) error
	WriteBlobChunk(BlobDigest, uint64, []byte) error
	DeleteBlobChunk(BlobDigest, uint64) error
	WriteStaleBlobIndex(StaleBlob) error
	DeleteStaleBlobIndex(StaleBlob) error

	jmt.RootWriter
	DeleteRoots(minVersionToKeep jmt.Version, maxItems int) (done bool, err error)

	jmt.NodeWriter
	jmt.StaleNodeWriter
	DeleteStaleNodes(maxStaleSinceVersion jmt.Version, maxItems int) (done bool, err error)

	DestructiveDestroyForTreeSync(n int) (done bool, err error)
}

type BlobMeta struct {
	PayloadLength uint64
	ChunksHave    []bool
	ExpirySeqNr   uint64
}

type StaleBlob struct {
	StaleSinceSeqNr uint64
	BlobDigest      BlobDigest
}

type KeyValueDatabase interface {
	// Must error if the key value store is not ready to apply state transition
	// for the given sequence number. Must update the highest committed sequence
	// number magic key upon commit. Convenience method for synchronization
	// between outcome generation & state sync.
	NewSerializedReadWriteTransaction(postSeqNr uint64) (KeyValueDatabaseReadWriteTransaction, error)
	// Must error if the key value store is not ready to apply state transition
	// for the given sequence number. Convenience method for synchronization
	// between outcome generation & state sync.
	NewReadTransaction(postSeqNr uint64) (KeyValueDatabaseReadTransaction, error)

	// Unchecked transactions are useful when you don't care that the
	// transaction state represents the kv state as of some particular sequence
	// number, mostly when writing auxiliary data to the kv store. Unchecked
	// transactions do not update the highest committed sequence number magic
	// key upon commit, as would checked transactions.
	NewSerializedReadWriteTransactionUnchecked() (KeyValueDatabaseReadWriteTransaction, error)
	// Unserialized transactions are guaranteed to commit.
	// The protocol should make sure that there are no conflicts across potentially concurrent unserialized transactions,
	// and if two unserialized transactions could actually have conflicts the protocol ensures that the are
	// never opened concurrently.
	NewUnserializedReadWriteTransactionUnchecked() (KeyValueDatabaseReadWriteTransaction, error)
	// Unchecked transactions are useful when you don't care that the
	// transaction state represents the kv state as of some particular sequence
	// number, mostly when reading auxiliary data from the kv store.
	NewReadTransactionUnchecked() (KeyValueDatabaseReadTransaction, error)

	// Deprecated: Kept for convenience/small diff, consider using
	// [KeyValueDatabaseSemanticRead.ReadHighestCommittedSeqNr] instead.
	HighestCommittedSeqNr() (uint64, error)
	Close() error
}

type KeyValueDatabaseFactory interface {
	NewKeyValueDatabase(configDigest types.ConfigDigest) (KeyValueDatabase, error)
}
