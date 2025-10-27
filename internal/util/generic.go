package util

import "golang.org/x/exp/constraints"

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
