package util

import (
	"maps"
)

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func PointerTo[T any](v T) *T {
	return &v
}

func PointerIntegerCast[U integer, T integer](p *T) *U {
	if p == nil {
		return nil
	}
	v := U(*p)
	return &v
}

func SaturatingSub[T unsigned](a T, b T) T {
	if b > a {
		var zero T
		return zero
	}
	return a - b
}

func NilCoalesce[T any](maybe *T, default_ T) T {
	if maybe != nil {
		return *maybe
	} else {
		return default_
	}
}

func NilCoalesceSlice[T any](maybe []T) []T {
	if maybe != nil {
		return maybe
	} else {
		return []T{}
	}
}

// b has priority in case of key conflict with a
func MapsUnion[K comparable, V any](a map[K]V, b map[K]V) map[K]V {
	c := make(map[K]V)
	maps.Copy(c, a)
	maps.Copy(c, b)
	return c
}
