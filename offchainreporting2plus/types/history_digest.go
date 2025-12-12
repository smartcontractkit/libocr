package types

import "fmt"

type HistoryDigest [32]byte

var _ fmt.Stringer = HistoryDigest{}

func (h HistoryDigest) String() string {
	return fmt.Sprintf("%x", h[:])
}

func BytesToHistoryDigest(b []byte) (HistoryDigest, error) {
	historyDigest := HistoryDigest{}

	if len(b) != len(historyDigest) {
		return HistoryDigest{}, fmt.Errorf("cannot convert bytes to HistoryDigest. bytes have wrong length %v", len(b))
	}

	if len(historyDigest) != copy(historyDigest[:], b) {
		// assertion
		panic("copy returned wrong length")
	}

	return historyDigest, nil
}
