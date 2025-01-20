package kvstore

import (
	"context"
	"sync"
)

var kv = &sync.Map{}

type Client struct {
	m *sync.Map
}

func Open() *Client {
	return &Client{kv}
}

func (kv *Client) Set(key string, value string) error {
	kv.m.Store(key, value)
	return nil
}

func (kv *Client) Exist(key string) bool {
	_, ok := kv.m.Load(key)
	return ok
}

func (kv *Client) SetByte(key string, value []byte) error {
	kv.m.Store(key, value)
	return nil
}

func (kv *Client) Load(key string) ([]byte, bool) {
	val, ok := kv.m.Load(key)
	if !ok {
		return nil, false
	}
	b, ok := val.([]byte)
	if !ok {
		return nil, false
	}
	return b, true
}

func (kv *Client) List(ctx context.Context, limit int) <-chan []byte {
	c := make(chan []byte)
	go func() {
		kv.m.Range(func(key, value any) bool {
			select {
			case <-ctx.Done():
				return false
			default:
				if limit == 0 {
					return false
				}
				b, ok := value.([]byte)
				if ok {
					c <- b
				}
				limit--
			}
			return true
		})
		close(c)
	}()
	return c
}
