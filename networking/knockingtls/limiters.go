package knockingtls

import (
	"fmt"
	"strings"
	"sync"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"golang.org/x/time/rate"
)

type refCountLimiter struct {
	limiter    *rate.Limiter // int64(math.Ceiling(float64):float64) -> SetLimit(float64(int64))
	refCount   int
	refillRate int64
}

// Limiters is an indexed collection such that each peer connection has a bandwidth rate limiter.
type Limiters struct {
	limiters map[p2ppeer.ID]refCountLimiter
	// Mutex for accessing all the properties of this collection!
	mu     sync.Mutex
	logger types.Logger
}

func NewLimiters(logger types.Logger) *Limiters {
	return &Limiters{
		make(map[p2ppeer.ID]refCountLimiter),
		sync.Mutex{},
		logger,
	}
}

// IncreaseLimits bumps the refill rate and bucket size for the specified peer ids.
// deltaTokenBucketRefillRate and deltaTokenBucketSize need to be either both positive or both negative. Otherwise they will be ignored.
func (ls *Limiters) IncreaseLimits(peerIDs []p2ppeer.ID, deltaTokenBucketRefillRate int64, deltaTokenBucketSize int) {
	if !((deltaTokenBucketRefillRate >= 0 && deltaTokenBucketSize >= 0) ||
		(deltaTokenBucketRefillRate <= 0 && deltaTokenBucketSize <= 0)) {
		ls.logger.Error("invariant violation: deltaTokenBucketRefillRate and deltaTokenBucketSize need to have the same sign", types.LogFields{
			"deltaTokenBucketRefillRate": deltaTokenBucketRefillRate,
			"deltaTokenBucketSize":       deltaTokenBucketSize,
		})
		return
	}

	positiveDeltas := deltaTokenBucketRefillRate > 0 || deltaTokenBucketSize > 0
	ls.mu.Lock()
	defer ls.mu.Unlock()
	for _, peerID := range peerIDs {

		// Figure out if there is a limiter for this peer. If there isn't and deltas are positive, add one. Othwerise, log error.
		rcLimiter, found := ls.limiters[peerID]
		if !found {
			if positiveDeltas {
				ls.limiters[peerID] = refCountLimiter{
					rate.NewLimiter(rate.Limit(float64(deltaTokenBucketRefillRate)), deltaTokenBucketSize),
					1,
					deltaTokenBucketRefillRate,
				}
			} else {
				ls.logger.Error("invariant violation: trying to decrease parameters for a rate limiter which doesn't exist", types.LogFields{
					"peerID":                     peerID.Pretty(),
					"deltaTokenBucketRefillRate": deltaTokenBucketRefillRate,
					"deltaTokenBucketSize":       deltaTokenBucketSize,
				})
			}
			continue
		}

		// Invariant at this point: the limiter for peerID exists and was not created just now.

		// Calculate and update new parameters for the existing limiter.
		// If the new parameters are negative, something went wrong, so log error and remove limiter.
		newLimit := rcLimiter.refillRate + deltaTokenBucketRefillRate
		newSize := rcLimiter.limiter.Burst() + deltaTokenBucketSize
		if newLimit < 0 || newSize < 0 {
			ls.logger.Error("incorrect new bandwith limiter params", types.LogFields{
				"peerID":         peerID.Pretty(),
				"newLimit":       newLimit,
				"newSize":        newSize,
				"referenceCount": rcLimiter.refCount,
			})
			delete(ls.limiters, peerID)
			continue
		} else {
			rcLimiter.limiter.SetLimit(rate.Limit(float64(newLimit)))
			rcLimiter.limiter.SetBurst(newSize)
			rcLimiter.refillRate = newLimit
		}

		// Invariant at this point: the limiter for peerID exists and has updated non-negative params.

		// Update reference count for the current limiter. If it's zero, log and remove the limiter.
		if positiveDeltas {
			rcLimiter.refCount += 1
		} else {
			rcLimiter.refCount -= 1
		}
		ls.limiters[peerID] = rcLimiter // We need to reassign because you can't change values associated with keys in a map!

		if rcLimiter.refCount == 0 {
			delete(ls.limiters, peerID)
			ls.logger.Info("removed bandwith limiter for peer connection as it's no longer used", types.LogFields{
				"peerID":         peerID.Pretty(),
				"referenceCount": rcLimiter.refCount,
				"currentLimit":   rcLimiter.limiter.Limit(),
				"currentSize":    rcLimiter.limiter.Burst(),
			})
		}
	}
}

// Find returns the limiter corresponding to the given peerID.
func (ls *Limiters) Find(peerID p2ppeer.ID) (*rate.Limiter, bool) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	l, ok := ls.limiters[peerID]
	if !ok {
		return nil, false
	}
	return l.limiter, true
}

func (ls *Limiters) String() string {
	var b strings.Builder
	b.WriteString("Limiters {")
	for key, l := range ls.limiters {
		b.WriteString(fmt.Sprintf("%s={refillRate=%f, size=%d, refcount=%d, refillRate=%d},",
			key.Pretty(), l.limiter.Limit(), l.limiter.Burst(), l.refCount, l.refillRate))
	}
	b.WriteString("}")
	return b.String()
}
