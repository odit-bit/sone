package domain

import "context"

type Entry struct {
	Endpoint string
	Name     string
}

type EntryRepository interface {
	List(ctx context.Context, offset, limit int) <-chan Entry
}
