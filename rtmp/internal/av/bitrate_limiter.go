package av

import (
	"fmt"
	"io"
	"time"
)

type LimitedBitrate struct {
	reader         io.Reader
	maxBitrateKbps uint32

	bitrateKbps float64
	readSize    uint64
	now         func() time.Time // for mock
	last        time.Time
}

func NewLimitedBitrate(r io.Reader, maxBitrateKbps uint32) *LimitedBitrate {
	return &LimitedBitrate{
		reader:         r,
		maxBitrateKbps: maxBitrateKbps,

		now: time.Now,
	}
}

func (r *LimitedBitrate) Read(b []byte) (int, error) {
	// Check bitrate first
	cur := r.now()
	if r.last.IsZero() {
		r.last = cur
	}
	diff := cur.Sub(r.last)
	if diff >= 1*time.Second {
		r.bitrateKbps = (float64(r.readSize) / float64(diff/time.Second)) * 8 / 1024.0
		// reset
		r.readSize = 0
		r.last = cur

		if r.bitrateKbps > float64(r.maxBitrateKbps) {
			return 0, &BitrateExceededError{
				MaxKbps:     float64(r.maxBitrateKbps),
				CurrentKbps: r.bitrateKbps,
			}
		}
	}

	n, err := r.reader.Read(b)
	if err != nil {
		return 0, err
	}

	r.readSize += uint64(n)

	return n, nil
}

func (r *LimitedBitrate) Close() error {
	if c, ok := r.reader.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

func (r *LimitedBitrate) BitrateKbps() float64 {
	return r.bitrateKbps
}

type BitrateExceededError struct {
	MaxKbps     float64
	CurrentKbps float64
}

func (err *BitrateExceededError) Error() string {
	return fmt.Sprintf(
		"Bitrate exceeded: Limit = %vkbps, Value = %vkbps",
		err.MaxKbps,
		err.CurrentKbps,
	)
}
