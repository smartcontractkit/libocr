package types

import "context"

// LOOPPContext is passed to functions that run in [local out-of-process
// plugins] and that were originally executed in the same
// process as OCR protocol logic.
// In the LOOPP architecture, these functions are executed in a different
// process (with very low round-trip latency, on the order of a millisecond in
// the typical case).
//
// Implementers of an interface function receiving such a context are
// responsible for ensuring that their function:
//   - returns quickly. This is important so as not to block protocol logic.
//   - doesn't perform remote network calls or other "slow" operations that are
//     allowed to run until context expiry.
//   - passes the LOOPPContext to other functions called down the stack. This is
//     to enable OpenTelemetry Tracing.
//
// Note that not all functions running in LOOPPs receive a LOOPPContext. Other
// functions that are expected to potentially run until context expiry receive
// a regular context.Context.
//
// [local out-of-process plugins]: https://github.com/smartcontractkit/chainlink-common/blob/main/pkg/loop/README.md?rgh-link-date=2024-03-19T16%3A48%3A39Z
type LOOPPContext = context.Context
