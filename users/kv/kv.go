package kv

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odit-bit/sone/pkg/kvstore"
	"github.com/odit-bit/sone/users/internal/domain"
)

var _ domain.StreamRepository = (*streamKV)(nil)

type streamKV struct {
	db *kvstore.Client
}

func New(kv *kvstore.Client) *streamKV {
	return &streamKV{db: kv}
}

// Exist implements domain.StreamRepository.
func (s *streamKV) Exist(ctx context.Context, key string) bool {
	return s.db.Exist(string(key))
}

// Save implements domain.StreamRepository.
func (s *streamKV) Save(ctx context.Context, stream domain.Stream) error {
	b, err := json.Marshal(stream)
	if err != nil {
		return err
	}
	s.db.SetByte(stream.ID, b)
	return nil
}

func (s *streamKV) Get(ctx context.Context, opt domain.StreamGetOption) (domain.Stream, bool, error) {

	var filter string
	if opt.Key != "" {
		filter = opt.Key
		return s.getKey(ctx, filter)
	}
	if opt.ID != "" {
		filter = opt.ID
		return s.getDefault(ctx, filter)
	}
	return domain.Stream{}, false, fmt.Errorf("no option is choose ")
}

func (s *streamKV) getDefault(_ context.Context, id string) (domain.Stream, bool, error) {
	b, ok := s.db.Load(id)
	if !ok {
		return domain.Stream{}, false, nil
	}
	var stream domain.Stream
	if err := json.Unmarshal(b, &stream); err != nil {
		return stream, false, err
	}
	return stream, true, nil
}

func (s *streamKV) getKey(ctx context.Context, key string) (domain.Stream, bool, error) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	c := s.db.List(ctx2, 1000)
	for b := range c {
		var stream domain.Stream
		if err := json.Unmarshal(b, &stream); err != nil {
			return stream, false, err
		}
		if stream.Key == key {
			return stream, true, nil
		}
	}
	return domain.Stream{}, false, nil
}

func (s *streamKV) List(ctx context.Context, limit int) <-chan domain.Stream {
	cStream := make(chan domain.Stream)
	go func() {
		c := s.db.List(ctx, limit)
		for b := range c {
			var stream domain.Stream
			if err := json.Unmarshal(b, &stream); err != nil {
				//log the error
				continue
			}
			cStream <- stream
		}
		close(cStream)
	}()
	return cStream
}
