package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/odit-bit/sone/streaming/internal/scheme"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Sqlite struct {
	LiveStream interface {
		Find(ctx context.Context, key string) (*scheme.LiveStream, bool, error)
		FindFilter(ctx context.Context, f scheme.Filter) (*scheme.LiveStreamsErr, error)
		Save(ctx context.Context, stream *scheme.LiveStream) (interface {
			Commit() error
			Cancel() error
		}, error)
	}
}

func NewSqlite(db *sql.DB) (Sqlite, error) {
	ls, err := NewSqliteLiveStream(db)
	if err != nil {
		return Sqlite{}, err
	}

	sq := Sqlite{
		ls,
	}
	return sq, nil
}

// func NewMockSqlite() Sqlite {
// 	return Sqlite{
// 		NewMockLiveStreamSQLite(),
// 	}
// }
