package buffer

import (
	"bytes"
	"sync"
)

var Pool = pool{
	pool: sync.Pool{
		New: func() any {
			b := &bytes.Buffer{}
			return b
		},
	},
}

type pool struct {
	pool sync.Pool
}

func (b *pool) Get() *bytes.Buffer {
	v := b.pool.Get()
	buf, _ := v.(*bytes.Buffer)
	return buf
}

func (b *pool) Put(buf *bytes.Buffer) {
	buf.Reset()
	b.pool.Put(buf)
}
