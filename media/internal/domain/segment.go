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
type LiveRepository interface {
	GetVideoSegment(ctx context.Context, key string, name string) (*Segment, error)
}

type SegmentStorer interface {
	InsertVideoSegment(ctx context.Context, path string, r io.Reader) error
}

type SegmentFetcher interface {
	GetVideoSegment(ctx context.Context, path string) (*Segment, error)
}

type SegmentRepository interface {
	SegmentFetcher
	SegmentStorer
}
