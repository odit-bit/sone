package domain

import (
	"context"
	"io"
)

// segment represent current segment streaming file
type Segment struct {
	Body io.ReadCloser
}

// type that can receive segment file
type SegmentRepository interface {
	Get(ctx context.Context, name string) (*Segment, error)
}
