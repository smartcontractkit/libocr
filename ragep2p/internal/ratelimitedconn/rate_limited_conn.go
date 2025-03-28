package ratelimitedconn

import (
	"fmt"
	"net"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
)

type Limiter interface {
	Allow(n int) bool
}

// TODO: would it make sense to merge this with the connRateLimiter?
type RateLimitedConn struct {
	net.Conn
	bandwidthLimiter    Limiter
	logger              commontypes.Logger
	readBytesTotal      prometheus.Counter
	writtenBytesTotal   prometheus.Counter
	rateLimitingEnabled bool
	transientCapacity   int
}

var _ net.Conn = (*RateLimitedConn)(nil)

func NewRateLimitedConn(
	conn net.Conn,
	bandwidthLimiter Limiter,
	logger commontypes.Logger,
	readBytesTotal prometheus.Counter,
	writtenBytesTotal prometheus.Counter,
) *RateLimitedConn {
	return &RateLimitedConn{
		conn,
		bandwidthLimiter,
		logger,
		readBytesTotal,
		writtenBytesTotal,
		false,
		0,
	}
}

// EnableRateLimiting is not thread-safe!
func (r *RateLimitedConn) EnableRateLimiting() {
	r.rateLimitingEnabled = true
}

// Allow a message of the specified size to be read during the subsequent calls to Read(...). The caller should specify
// the message size directly, and may use addTLSOverhead=true to account for TLS overhead. After fully reading the
// message, ClearTransientCapacity() must be called.
func (r *RateLimitedConn) AllowTransientCapacity(expectedPayloadSize int, addTLSOverhead bool) {
	r.transientCapacity = expectedPayloadSize
	if addTLSOverhead {
		// In principle, TLS 1.3 uses record of size 16 KiB (when transmitting sufficiently large messages).
		// Per record, we have an overhead of 21 bytes.
		//
		// Lets be a bit more generous here and assume:
		//  - a max record size of 1400 byte which does fit into all MTUs, and
		//  - a per record overhead of 64 bytes.
		//
		// This comes out to: 64/1400  <  5% (= 1/20).

		r.transientCapacity += max(64, expectedPayloadSize/20)
	}
}

func (r *RateLimitedConn) ClearTransientCapacity() {
	r.transientCapacity = 0
}

func (r *RateLimitedConn) Read(b []byte) (n int, err error) {
	if !r.rateLimitingEnabled {
		n, err = r.Conn.Read(b)
		r.readBytesTotal.Add(float64(n))
		return n, err
	}

	// Check if we have (remaining) transient capacity.
	if r.transientCapacity > 0 {
		// If so, we are processing a response.
		// We read up to the number of remaining transient capacity bytes from the underlying connection for that
		// response and deduct the number of bytes read from the transient capacity.

		// Setting a upper bound (i.e., r.transientCapacity) for the number of bytes to read, and not using the
		// internally [and dynamically adjusted] TLS buffer size len(b), is important to ensure responses fit within the
		// calculated transient capacity. Without that consideration, extra available data (non-belong to the response)
		// could be wrongfully counted towards the responses size due to the way this Read() function is called from
		// the TLS code - wrongfully exceeding the rate limit set for the response (if the normal bandwidth limiter's
		// capacity is close to zero.)
		numBytesToRead := min(r.transientCapacity, len(b))

		n, err = r.Conn.Read(b[:numBytesToRead])
		r.readBytesTotal.Add(float64(n))
		r.transientCapacity -= n
		return n, err
	}

	// If we have no (remaining) transient capacity, we apply the normal rate limits.

	n, err = r.Conn.Read(b)
	r.readBytesTotal.Add(float64(n))
	nBytesAllowed := r.bandwidthLimiter.Allow(n)
	if nBytesAllowed {
		return n, err
	}

	// kill the conn: close it and emit an error
	_ = r.Conn.Close() // ignore error, there's not much we can with it here

	// TODO: log the limits here
	r.logger.Error("inbound data exceeded rate limit, connection closed", commontypes.LogFields{
		// "tokenBucketRefillRate": r.bandwidthLimiter.Limit(),
		// "tokenBucketSize":       r.bandwidthLimiter.Burst(),
		"bytesRead": n,
		"readError": err, // This error may not be null, we're adding it here to not miss it.
	})

	return 0, fmt.Errorf("inbound data exceeded rate limit, connection closed")
}

// func (r *RateLimitedConn) Read(b []byte) (n int, err error) {
// 	n, err = r.Conn.Read(b)
// 	r.readBytesTotal.Add(float64(n))
// 	if !r.rateLimitingEnabled {
// 		return n, err
// 	}

// 	// Hold the number of bytes read minus any potentially deducted allowance from the transient capacity.
// 	// We use a separate variable for this and keep the number of read bytes unmodified in `n`.
// 	nd := n

// 	// Check if a transient capacity was set.
// 	if r.transientCapacity != 0 {
// 		// If so, we use it, but we have to distinguish two cases: either we have enough capacity left or not.
// 		// Handle case 1: The (remaining) transient capacity is sufficient, so we deduct the number of read bytes from
// 		// it and return immediately.
// 		if n <= r.transientCapacity {
// 			r.transientCapacity -= n
// 			return n, err
// 		}

// 		// Handle case 2: The (remaining) transient capacity is not sufficient, so we use up all the remaining capacity,
// 		// and fallthrough to apply the normal bandwidth limiter for the remaining number of read bytes.
// 		nd = n - r.transientCapacity
// 		r.transientCapacity = 0
// 	}

// 	ndBytesAllowed := r.bandwidthLimiter.Allow(nd)
// 	if ndBytesAllowed {
// 		return n, err
// 	}

// 	// kill the conn: close it and emit an error
// 	_ = r.Conn.Close() // ignore error, there's not much we can with it here

// 	// TODO: log the limits here
// 	r.logger.Error("inbound data exceeded rate limit, connection closed", commontypes.LogFields{
// 		// "tokenBucketRefillRate": r.bandwidthLimiter.Limit(),
// 		// "tokenBucketSize":       r.bandwidthLimiter.Burst(),
// 		"bytesRead": n,
// 		"readError": err, // This error may not be null, we're adding it here to not miss it.
// 	})

// 	return 0, fmt.Errorf("inbound data exceeded rate limit, connection closed")
// }

func (r *RateLimitedConn) Write(b []byte) (n int, err error) {
	n, err = r.Conn.Write(b)
	r.writtenBytesTotal.Add(float64(n))
	return n, err
}
