package overheadawareconn

import (
	"fmt"
	"net"

	"github.com/prometheus/client_golang/prometheus"
)

// OverheadAwareConn keeps track of TLS overhead (vs application data consumed by ragep2p)
// and enables us to kill the connection is the overhead is suspiciously high.
// This is to defend against an adversary causing undue overhead by e.g. fragmenting their TLS
// records into tiny pieces and other similar shenanigans.
//
// Note that the overhead computation is specific to ragep2p's use of TLS. This is not useful
// for other applications.
//
// Not thread-safe, except for functions exposed directly from the embedded net.Conn
// such as Close() and SetDeadline().
type OverheadAwareConn struct {
	net.Conn
	readBytesTotal                prometheus.Counter
	writtenBytesTotal             prometheus.Counter
	setupComplete                 bool
	readPostSetupRawBytes         int64 // raw = including TLS record & encrypption overhead
	deliveredApplicationDataBytes int64
}

var _ net.Conn = (*OverheadAwareConn)(nil)

func NewOverheadAwareConn(
	conn net.Conn,
	readBytesTotal prometheus.Counter,
	writtenBytesTotal prometheus.Counter,
) *OverheadAwareConn {
	return &OverheadAwareConn{
		conn,
		readBytesTotal,
		writtenBytesTotal,

		false,
		0,
		0,
	}
}

func (r *OverheadAwareConn) SetupComplete() {
	r.setupComplete = true
}

func (r *OverheadAwareConn) Read(b []byte) (n int, err error) {
	n, err = r.Conn.Read(b)
	r.readBytesTotal.Add(float64(n))

	if r.setupComplete {
		r.readPostSetupRawBytes += int64(n)
	}

	return n, err
}

func (r *OverheadAwareConn) Write(b []byte) (n int, err error) {
	n, err = r.Conn.Write(b)
	r.writtenBytesTotal.Add(float64(n))
	return n, err
}

// Whenever ragep2p reads a message, that message starts with a frame header.
// The frame header's size is at least 37 bytes. It will always fit into a single TLS record.
// The per-record overhead of TLS 1.3 in our configuration is 22 bytes.
// So even for a worst-case scenario of the user only sending empty ragep2p messages,
// we should only have a factor of (37+22)/37 ≈ 1.6. We conservatively round up to 2.
//
// Intuitively, it is obvious that the overhead will go down for larger message, but let's make
// this a bit more concrete:
//   - If frame header size + payload size is less than 2048 bytes, ragep2p will coalesce the write
//     into a single TLS record, giving us a factor of (2048+22)/2048 ≈ 1.01
//   - Otherwise, ragep2p will send separate records for the frame header and the payload, giving us
//     a factor of (2049 + 2*22)/2049 ≈ 1.02 or smaller.
//   - Once messages reach the maximum TLS record size of 16348 bytes, they need to be fragmented.
//     Again, application data is so large at this point that the record overhead is negligible:
//     (37 + 16349 + 3*22)/(37 + 16349) ≈ 1.005
const MaximumAllowedApplicationDataToRawFactor = 2

// Golang's TLS implementation performs some read-ahead buffering, e.g. when buffering an incomplete
// record. (Decryption can only take place once the whole record has been read.)
// Let's *generously* overestimate how much.
const GenerouslyOverstimatedTLSReadAhead = 100 * 1024

// AddDeliveredApplicationDataBytes tells the OverheadAwareConn how many bytes of application data
// have been delivered to the user and returns an error if it seems that the overhead
// on the underlying net.Conn is too large.
func (r *OverheadAwareConn) AddDeliveredApplicationDataBytes(bytes int) error {
	if bytes < 0 {
		return fmt.Errorf("bytes must be non-negative")
	}

	r.deliveredApplicationDataBytes += int64(bytes)

	readPostSetupRawBytesAttributableToDeliveredApplicationData := r.readPostSetupRawBytes - GenerouslyOverstimatedTLSReadAhead

	if r.deliveredApplicationDataBytes*MaximumAllowedApplicationDataToRawFactor < readPostSetupRawBytesAttributableToDeliveredApplicationData {
		return fmt.Errorf("inbound read overhead on underlying TCP connection is too large, suspecting shenanigans: deliveredApplicationDataBytes=%d, readPostSetupRawBytes=%d, generouslyOverstimatedTLSReadAhead=%d",
			r.deliveredApplicationDataBytes,
			r.readPostSetupRawBytes,
			GenerouslyOverstimatedTLSReadAhead,
		)
	}

	return nil
}
