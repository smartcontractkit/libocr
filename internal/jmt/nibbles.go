package jmt

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
)

type NibblePath struct {
	numNibbles int
	bytes      []byte
}

func NewNibblePath(numNibbles int, bytes []byte) (NibblePath, bool) {
	expectedBytesLen := (numNibbles + 1) / 2
	if len(bytes) != expectedBytesLen {
		return NibblePath{}, false
	}
	if numNibbles%2 == 1 && bytes[len(bytes)-1]&0x0f != 0 {
		return NibblePath{}, false
	}
	return NibblePath{numNibbles, slices.Clone(bytes)}, true
}

func (np NibblePath) Compare(np2 NibblePath) int {

	bytesCmp := bytes.Compare(np.bytes, np2.bytes)
	if bytesCmp != 0 {
		return bytesCmp
	}
	return cmp.Compare(np.numNibbles, np2.numNibbles)
}

func (np NibblePath) Equal(np2 NibblePath) bool {
	return np.numNibbles == np2.numNibbles && bytes.Equal(np.bytes, np2.bytes)
}

func (np NibblePath) Get(index int) byte {
	if index >= np.numNibbles {

		panic("index out of bounds")
	}
	b := np.bytes[index/2]
	var n byte
	if index%2 == 0 {
		n = (b & 0xf0) >> 4
	} else {
		n = b & 0x0f
	}
	return n
}

func (np NibblePath) NumNibbles() int {
	return np.numNibbles
}

func (np NibblePath) Bytes() []byte {
	return np.bytes
}

func (np NibblePath) Append(nibble byte) NibblePath {
	np2 := NibblePath{
		np.numNibbles + 1,
		bytes.Clone(np.bytes),
	}
	if np.numNibbles%2 == 1 {
		np2.bytes[len(np2.bytes)-1] = np2.bytes[len(np2.bytes)-1]&0xf0 | nibble
		return np2
	} else {
		np2.bytes = append(np2.bytes, (nibble&0x0f)<<4)
		return np2
	}
}

func (np NibblePath) Prefix(count int) NibblePath {
	nbytes := (count + 1) / 2
	keepbytes := bytes.Clone(np.bytes[:nbytes])
	if count%2 == 1 {
		keepbytes[len(keepbytes)-1] = keepbytes[len(keepbytes)-1] & 0xf0
	}
	return NibblePath{
		count,
		keepbytes,
	}
}

func (np NibblePath) String() string {
	r := fmt.Sprintf("%x", np.bytes)
	if np.numNibbles%2 == 1 {
		r = r[:len(r)-1]
	}
	return r
}

func NibblePathFromDigest(digest Digest) NibblePath {
	return NibblePath{
		len(digest) * 2,
		digest[:],
	}
}
