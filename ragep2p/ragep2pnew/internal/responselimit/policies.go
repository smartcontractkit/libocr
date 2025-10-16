package responselimit

import (
	"time"

	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/internaltypes"
)

//go-sumtype:decl ResponsePolicy

// Interface for specifying rate-limit exceptions for responses.
//
// When a request is made, a response policy is used to specify if a response (or in principle multiple responses)
// should be allowed or rejected.
//
// Policies have to be tracked internally. Therefore, to allow proper cleanup of resources, it is critical that
// IsPolicyExpired(...) returns true when no (additional) response is expected anymore.
type ResponsePolicy interface {
	// Must return true as soon as the policy is no longer required, i.e., it must return true if the policy is expired
	// at the provided timestamp (or any later point in time).
	isPolicyExpired(timestamp time.Time) bool

	// Specifies whether a response for the given request ID should be allowed or rejected.
	// Before and after checkResponse(...) is called internally, a policy is always checked for expiry.
	// checkResponse(...) is never called on an expired policy.
	checkResponse(requestID internaltypes.RequestID, responseSize int, responseTimestamp time.Time) ResponseCheckResult
}

var _ ResponsePolicy = &SingleUseSizedLimitedResponsePolicy{}

// A response policy, allowing at most one response subject to the following constraints:
//   - the response's payload size is at most `MaxSize`
//   - the response is received before `ExpiryTimestamp`
type SingleUseSizedLimitedResponsePolicy struct {
	MaxSize         int
	ExpiryTimestamp time.Time
}

func (p *SingleUseSizedLimitedResponsePolicy) isPolicyExpired(timestamp time.Time) bool {
	return !timestamp.Before(p.ExpiryTimestamp)
}

func (p *SingleUseSizedLimitedResponsePolicy) checkResponse(
	requestID internaltypes.RequestID,
	responseSize int,
	responseTimestamp time.Time,
) ResponseCheckResult {
	// As this is intended to be a single use policy only, we set the timestamp to its zero value. This is ensuring that
	// any subsequent call to `isPolicyExpired(...)` returns false. As consequence `checkResponse(...)`, will not be
	// called on again by the response checker, and the policy will be removed from the checker.
	p.ExpiryTimestamp = time.Time{}

	if responseSize > p.MaxSize {
		return ResponseCheckResultReject
	}
	return ResponseCheckResultAllow
}
