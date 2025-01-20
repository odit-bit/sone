package domain

import (
	"context"
	"io"
)

type OpError struct {
	message string
}

func (op *OpError) Error() string {
	return op.message
}

func NewOpError(msg string) *OpError {
	return &OpError{message: msg}
}

type HLSRepository interface {
	//convert flv stream data to hls segment file and persisted it
	TranscodeFLV(ctx context.Context, stream_key string, r io.Reader) <-chan *OpError
}

type StreamKeyRepository interface {
	Exist(key string) bool
}
