package web

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/odit-bit/sone/web/internal/session"
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

func init() {
	gob.Register(time.Time{})
}

var _ session.Manager = (*sessionManager)(nil)

type sessionManager struct {
	store *scs.SessionManager
}

func NewSessionManager() *sessionManager {
	store := scs.New()
	store.Store = memstore.New()
	return &sessionManager{store: store}
}

func (sm *sessionManager) Load(r *http.Request) (*session.Session, error) {

	var token string
	cookie, err := r.Cookie(sm.store.Cookie.Name)
	if err == nil {
		token = cookie.Value
	}

	ctx, err := sm.store.Load(r.Context(), token)
	if err != nil {
		return nil, err
	}

	return session.NewSession(ctx, sm.store), nil
}

func (sm *sessionManager) Save(data *session.Session, w http.ResponseWriter) error {
	ctx := data.Context()
	w.Header().Add("Vary", "Cookie")
	sm.store.Put(ctx, _lastLogin, time.Now())

	switch sm.store.Status(data.Context()) {
	case scs.Modified:
		token, expiry, err := sm.store.Commit(ctx)
		if err != nil {
			return err
		}
		sm.store.WriteSessionCookie(ctx, w, token, expiry)
	case scs.Destroyed:
		sm.store.WriteSessionCookie(ctx, w, "", time.Time{})
	}

	return nil
}
