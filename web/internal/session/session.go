package session

import (
	"context"
	"net/http"
	"time"
)

const (
	_ int = iota
	Login
)

const (
	_isLogin   = "is_login"
	_username  = "username"
	_lastLogin = "last_login"
)

type Manager interface {
	Load(r *http.Request) (*Session, error)
	Save(data *Session, w http.ResponseWriter) error
}

type Store interface {
	GetBool(ctx context.Context, key string) bool
	GetString(ctx context.Context, key string) string
	GetTime(ctx context.Context, key string) time.Time
	Put(ctx context.Context, key string, data any)
}

type Session struct {
	store Store
	ctx   context.Context
}

func NewSession(ctx context.Context, store Store) *Session {
	return &Session{store: store, ctx: ctx}
}

func (s *Session) Context() context.Context {
	return s.ctx
}

func (s *Session) IsLogin() bool {
	return s.store.GetBool(s.ctx, _isLogin)
}

func (s *Session) Name() string {
	return s.store.GetString(s.ctx, _username)
}

func (s *Session) LastModified() time.Time {
	return s.store.GetTime(s.ctx, _lastLogin)
}

func (s *Session) SetLogin(b bool) {
	s.store.Put(s.ctx, _isLogin, b)
}

func (s *Session) SetName(name string) {
	s.store.Put(s.ctx, _username, name)
}

func (s *Session) Put(key string, val any) {
	s.store.Put(s.ctx, key, val)
}

func (s *Session) GetString(key string) string {
	return s.store.GetString(s.ctx, key)
}
