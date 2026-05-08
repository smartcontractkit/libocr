// Package bigendianbytearray provides arithmetic operations on fixed-size byte
// arrays interpreted as unsigned big-endian integers.
package bigendianbytearray

func Min32[T ~[32]byte]() T {
	return T{}
}

func Max32[T ~[32]byte]() T {
	var t T
	for i := range t {
		t[i] = 0xff
	}
	return t
}

func Decrement32[T ~[32]byte](t T) (T, bool) {
	tdec := t
	for i := len(tdec) - 1; i >= 0; i-- {
		if tdec[i] == 0 {
			tdec[i] = 0xff
		} else {
			tdec[i]--
			return tdec, true
		}
	}
	return T{}, false
}

func WrappingDecrement32[T ~[32]byte](t T) T {
	decr, ok := Decrement32(t)
	if ok {
		return decr
	} else {
		return Max32[T]()
	}
}

func Increment32[T ~[32]byte](t T) (T, bool) {
	tincr := t
	for i := len(tincr) - 1; i >= 0; i-- {
		if tincr[i] == 0xff {
			tincr[i] = 0
		} else {
			tincr[i]++
			return tincr, true
		}
	}
	return T{}, false
}

func WrappingIncrement32[T ~[32]byte](t T) T {
	incr, ok := Increment32(t)
	if ok {
		return incr
	} else {
		return Min32[T]()
	}
}
