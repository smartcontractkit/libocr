package networking

import "golang.org/x/time/rate"

// limiter is a subset of the x/time/rate.Limiter api
type limiter interface {
	Allow() bool
}

// newLimiter builds a limiter with a specified number of allowed message per second and a burst value.
func newLimiter(tokenBucketRefillRate float64, tokenBucketSize int) limiter {
	return rate.NewLimiter(rate.Limit(tokenBucketRefillRate), tokenBucketSize)
}
