package networking

import "golang.org/x/time/rate"


type limiter interface {
	Allow() bool
}


func newLimiter(tokenBucketRefillRate float64, tokenBucketSize int) limiter {
	return rate.NewLimiter(rate.Limit(tokenBucketRefillRate), tokenBucketSize)
}
