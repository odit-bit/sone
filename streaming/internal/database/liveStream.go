package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/odit-bit/sone/streaming/internal/scheme"
)

type SQLiteLiveStream struct {
	// tableName string
	DB *sql.DB
}

func (s *SQLiteLiveStream) migrate(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS live_streams (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		is_live INTEGER NOT NULL,

		key TEXT NOT NULL
	);
`
	_, err := s.DB.ExecContext(ctx, query)
	return err
}

func NewSqliteLiveStream(db *sql.DB) (*SQLiteLiveStream, error) {
	sq := SQLiteLiveStream{
		DB: db,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := sq.migrate(ctx); err != nil {
		return nil, err
	}
	return &sq, nil
}

// FindKey implements StreamRepo.
func (s *SQLiteLiveStream) Find(ctx context.Context, id string) (*scheme.LiveStream, bool, error) {
	ls := scheme.LiveStream{}
	query := "SELECT id, title, is_live, key FROM live_streams WHERE id = ? LIMIT 1"
	if err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&ls.ID,
		&ls.Title,
		&ls.IsLive,
		&ls.Key,
	); err != nil {
		return nil, false, ErrRecordNotFound
	}

	return &ls, true, nil
}

func (s *SQLiteLiveStream) FindFilter(ctx context.Context, f scheme.Filter) (*scheme.LiveStreamsErr, error) {
	var query string
	col, val, ok := s.filter(f)
	if ok {
		query = fmt.Sprintf("SELECT id, title, is_live, key FROM live_streams WHERE %s = %s LIMIT 1;", col.query(), col.placeholders())
	} else {
		query = "SELECT id, title, is_live, key FROM live_streams limit 1000"
	}
	rows, err := s.DB.QueryContext(ctx, query, val...)
	if err != nil {
		return nil, ErrRecordNotFound
	}

	c := make(chan *scheme.LiveStream)
	cErr := make(chan error, 1)
	go func() {
		defer close(cErr)
		defer close(c)
		defer rows.Close()

		for rows.Next() {
			var ls scheme.LiveStream
			err := rows.Scan(
				&ls.ID,
				&ls.Title,
				&ls.IsLive,
				&ls.Key,
			)
			if err != nil {
				select {
				case <-ctx.Done():
				case cErr <- err:
				}
				return
			}
			select {
			case c <- &ls:
			case <-ctx.Done():
				return
			}
		}

	}()

	list := scheme.LiveStreamIterator(c, cErr)
	return list, nil
}

type columns []string

func (c *columns) query() string {
	s := strings.Join(*c, ",")
	return s
}

func (c *columns) placeholders() string {
	const ph = "?,"
	s := strings.Repeat(ph, len(*c))
	return s[:len(s)-1]
}

func (s *SQLiteLiveStream) filter(f scheme.Filter) (columns, []any, bool) {
	var useFilter bool
	field := columns{}
	values := []any{}
	if f.Id != "" {
		field = append(field, "id")
		values = append(values, f.Id)
		useFilter = true
	}
	if f.Key != "" {
		field = append(field, "key")
		values = append(values, f.Key)
		useFilter = true

	}
	if f.Title != "" {
		field = append(field, "title")
		values = append(values, f.Title)
		useFilter = true

	}

	if f.IsLive {
		field = append(field, "is_live")
		values = append(values, 1)
		useFilter = true

	}

	return field, values, useFilter

}

// Save implements StreamRepo.
func (s *SQLiteLiveStream) Save(ctx context.Context, stream *scheme.LiveStream) (interface {
	Commit() error
	Cancel() error
}, error) {
	query := `
	INSERT INTO live_streams (id, title, is_live, key) 
	VAlUES (?, ?, ?, ?)
	ON CONFLICT (id) DO
	UPDATE
	SET 
		title = excluded.title,
		is_live = excluded.is_live,
		key = excluded.key;
	`

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	cErr := make(chan error)
	go func() {
		if _, err := tx.ExecContext(ctx, query, stream.ID, stream.Title, stream.IsLive, stream.Key); err != nil {
			cErr <- err
		}
		close(cErr)
	}()

	return &TX{tx: tx, ctx: ctx, cErr: cErr}, nil
}

type TX struct {
	ctx  context.Context
	cErr chan error
	tx   *sql.Tx
}

func (t *TX) Commit() error {
	var commitErr error
	select {
	case err := <-t.cErr:
		commitErr = errors.Join(err)
	case <-t.ctx.Done():
		commitErr = errors.Join(t.ctx.Err())
	}
	return errors.Join(commitErr, t.tx.Commit())
}

func (t *TX) Cancel() error {
	var commitErr error
	select {
	case err := <-t.cErr:
		commitErr = errors.Join(err)
	case <-t.ctx.Done():
		commitErr = errors.Join(t.ctx.Err())
	}
	return errors.Join(commitErr, t.tx.Rollback())
}

//// MOCKING

type mockLiveStreamSQLite struct{}

func NewMockLiveStreamSQLite() *mockLiveStreamSQLite {
	return &mockLiveStreamSQLite{}
}

func (mock *mockLiveStreamSQLite) Find(ctx context.Context, key string) (*scheme.LiveStream, bool, error) {
	return &scheme.LiveStream{}, true, nil
}
func (mock *mockLiveStreamSQLite) Save(ctx context.Context, stream *scheme.LiveStream) (interface {
	Commit() error
	Cancel() error
}, error) {
	return nil, nil
}
