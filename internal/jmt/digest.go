package jmt

import (
	"bytes"
	"crypto/sha256"
)

type Digest = [sha256.Size]byte

var (
	MinDigest = Digest{}
	MaxDigest = Digest(bytes.Repeat([]byte{0xff}, len(Digest{})))
)

func DecrementDigest(digest Digest) (Digest, bool) {
	decDigest := digest
	for i := len(decDigest) - 1; i >= 0; i-- {
		if decDigest[i] == 0 {
			decDigest[i] = 0xff
		} else {
			decDigest[i]--
			return decDigest, true
		}
	}
	return Digest{}, false
}

func IncrementDigest(digest Digest) (Digest, bool) {
	incDigest := digest
	for i := len(incDigest) - 1; i >= 0; i-- {
		if incDigest[i] == 0xff {
			incDigest[i] = 0
		} else {
			incDigest[i]++
			return incDigest, true
		}
	}
	return Digest{}, false
}
