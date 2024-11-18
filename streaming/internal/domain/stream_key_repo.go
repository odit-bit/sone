package domain

import "context"

type StreamKeyRepository interface {
	Save(ctx context.Context, key Key) error
	Exist(ctx context.Context, key Key) bool
}
