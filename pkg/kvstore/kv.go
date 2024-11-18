package kvstore

import "sync"

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
