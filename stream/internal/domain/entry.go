package domain

import (
	"context"
	"io"
)

// represent stream entry
type Entry struct {
	Path string
	Name string
}




// represent segment file
type Segment struct {
	Body io.ReadCloser
}

// stream data is store as transport stream format (.ts)
type EntryRepository interface {
	// underlying type will populate the segment arg
	GetSegment(name string, segment *Segment) error
	// get list of current stream
	List(offset, limit int) <-chan Entry
}

type IngestErr struct {
	Err error
}

// represent type that ingest serialized stream for further process
type RTMPIngester interface {
	IngestFLV(ctx context.Context, name string, r io.Reader) <-chan IngestErr
}
