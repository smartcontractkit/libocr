package util

func PointerTo[T any](v T) *T {
	return &v
}

func NilCoalesce[T any](maybe *T, default_ T) T {
	if maybe != nil {
		return *maybe
	} else {
		return default_
	}
}
