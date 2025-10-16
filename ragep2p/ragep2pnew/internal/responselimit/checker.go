package responselimit

import (
	"math/rand"
	"sync"
	"time"

	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/internaltypes"
)

type ResponseCheckResult byte

// Enum specifying the list of return values for responseChecker.CheckResponse(...).
const (
	// A response is rejected if the policy
	// (1) was not found, or
	// (2) was expired, or
	// (3) was found but decided to reject the request.
	//
	// As policies are automatically cleaned up (in some non-deterministic manner), there is no way to distinguish
	// cases (1) and (2), and for simplicity also case (3) is handled identically.
	//
	// We intentionally use 0 as the first enum value for Reject as a safe default here.
	ResponseCheckResultReject ResponseCheckResult = iota

	// A (non-expired) policy was found, and the policy did decide that the response should be allowed.
	ResponseCheckResultAllow
)

type responseCheckerMapEntry struct {
	index    int
	policy   ResponsePolicy
	streamID internaltypes.StreamID
}

// Data structure for keeping track of open requests until a set expiry date.
//
// Cleanup of expired entries is performed automatically. Whenever a new entry is added, two random entries are checked
// and removed if expired. This ensures that, on expectation, the number of tracked entries is approx. 2x the number
// of non-expired entries.
//
// SetPolicy(...) and CheckResponse(...) are O(1) operations.
type ResponseChecker struct {
	mutex    sync.Mutex
	rids     []internaltypes.RequestID
	policies map[internaltypes.RequestID]responseCheckerMapEntry
	rng      *rand.Rand
}

func NewResponseChecker() *ResponseChecker {
	return &ResponseChecker{
		sync.Mutex{},
		make([]internaltypes.RequestID, 0),
		make(map[internaltypes.RequestID]responseCheckerMapEntry),
		rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Sets the policy for a given (fresh) request ID. After setting the policy, calling Pop(...) for the same ID before the
// policy expires returns the policy Set with this function. If a policy with the provided ID is already present, it
// will be overwritten.
func (c *ResponseChecker) SetPolicy(sid internaltypes.StreamID, rid internaltypes.RequestID, policy ResponsePolicy) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Lookup an existing policy for the provided request ID.
	// If it exists, we override the policy, keeping its location at the prior index.
	// Otherwise, we need use a new index and also track the request ID in the c.rids list.
	entry, exists := c.policies[rid]
	if exists {
		entry = responseCheckerMapEntry{entry.index, policy, sid}
	} else {
		// We set entry.index = len(c.rids) to let it point to the request ID we will append to c.rids list.
		entry = responseCheckerMapEntry{len(c.rids), policy, sid}
		c.rids = append(c.rids, rid)
	}

	// Actually save the policy update back to the c.policies map.
	c.policies[rid] = entry

	// If the number of tracked policies increased, we check 2 random policies and remove them if expired. This way
	// the number of tracked policies only grows to 2x the number of non-expired policies in expectation.
	if !exists {
		c.cleanupExpired()
	}
}

// Lookup the policy for a given response and check if it should be allowed or rejected.
// See responseCheckResult for additional documentation on the potential return values of this function.
func (c *ResponseChecker) CheckResponse(sid internaltypes.StreamID, rid internaltypes.RequestID, size int) ResponseCheckResult {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.policies[rid]
	if !exists {
		return ResponseCheckResultReject
	}
	if entry.streamID != sid {
		return ResponseCheckResultReject
	}

	now := time.Now()
	if entry.policy.isPolicyExpired(now) {
		c.removeEntry(rid, entry.index)
		return ResponseCheckResultReject
	}

	policyResult := entry.policy.checkResponse(rid, size, now)

	// Recheck the policy of expiry, useful to cleanup one-time-use policies immediately.
	if entry.policy.isPolicyExpired(now) {
		c.removeEntry(rid, entry.index)
	}

	return policyResult
}

// Removes all currently tracked policies for the given stream ID. To ensure that responses sent to a stream cannot be
// accepted after this stream is closed and reopened, this function is called when the Stream is closed (and removed
// from the demuxer).
func (c *ResponseChecker) ClearPoliciesForStream(sid internaltypes.StreamID) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := 0; i < len(c.rids); i++ {
		rid := c.rids[i]
		policy := c.policies[rid]

		if policy.streamID == sid {
			// We found a policy which matches the given stream ID.
			// So we remove the entry from the list of request IDs and policies.
			c.removeEntry(rid, i)

			// The above removeEntry(...) removes c.rids[i], thus in the next iteration index its value is replaced
			// by a different request ID. We decrement index i to ensure that we don't skip the new value at index i.
			i--
		}
	}
}

// Check two random policies. A checked policy is removed if it is found to be expired.
func (c *ResponseChecker) cleanupExpired() {
	now := time.Now()

	// At most 2 iterations, enter loop body only if c.rids is non empty.
	for i := 0; i < 2 && len(c.rids) > 0; i++ {
		// Select a random policy.
		index := c.rng.Intn(len(c.rids))
		id := c.rids[index]
		policy := c.policies[id].policy

		// Remove it if it is expired.
		if policy.isPolicyExpired(now) {
			c.removeEntry(id, index)
		}
	}
}

// Remove the policy for a given request ID from (1) the map of policies and (2) the list of request IDs.
func (c *ResponseChecker) removeEntry(id internaltypes.RequestID, index int) {
	// Remove the entry from the map of polices.
	delete(c.policies, id)

	// Handle the "index == last-index" corner case separately.
	// This avoids wrongfully reinserting the deleted policy.
	if index == len(c.rids)-1 {
		c.rids = c.rids[0 : len(c.rids)-1]
		return
	}

	// Swap the last entry's id to the position of the to be removed id, and remove the last value from the rids list.
	lastID := c.rids[len(c.rids)-1]
	c.rids[index] = lastID
	c.rids = c.rids[0 : len(c.rids)-1]

	// Update the index point for the c.policies[lastId] to point to the now changed position.
	lastEntry := c.policies[lastID]
	lastEntry.index = index
	c.policies[lastID] = lastEntry
}
