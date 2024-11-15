package domain

import "context"

type StreamKeyRepo interface {
	Save(ctx context.Context, key Key) error
	Exist(ctx context.Context, key Key) bool
}
