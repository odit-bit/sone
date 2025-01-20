package domain

import "context"

type StreamGetOption struct {
	ID  string
	Key string
}

type StreamRepository interface {
	Save(ctx context.Context, stream Stream) error
	Exist(ctx context.Context, key string) bool
	Get(ctx context.Context, opt StreamGetOption) (Stream, bool, error)
	List(ctx context.Context, limit int) <-chan Stream
}
