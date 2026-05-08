package util

import (
	"maps"

	"golang.org/x/exp/constraints"
)

func PointerTo[T any](v T) *T {
	return &v
}

func PointerIntegerCast[U constraints.Integer, T constraints.Integer](p *T) *U {
	if p == nil {
		return nil
	}
	v := U(*p)
	return &v
}

func SaturatingSub[T constraints.Unsigned](a T, b T) T {
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
