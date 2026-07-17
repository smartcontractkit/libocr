package jmt

import (
	"crypto/sha256"

	"github.com/smartcontractkit/libocr/internal/bigendianbytearray"
)

type Digest = [sha256.Size]byte

var (
	MinDigest = bigendianbytearray.Min32[Digest]()
	MaxDigest = bigendianbytearray.Max32[Digest]()
)

func DecrementDigest(digest Digest) (Digest, bool) {
	return bigendianbytearray.Decrement32(digest)
}

func IncrementDigest(digest Digest) (Digest, bool) {
	return bigendianbytearray.Increment32(digest)
}
