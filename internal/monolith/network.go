package monolith

import (
	"net/url"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

// HTTP

type HTTP struct {
	addr string
	mux  *chi.Mux
}

func NewHTTP(addr string, mux *chi.Mux) *HTTP {
	u := url.URL{}
	u.Scheme = "http"
	u.Host = addr

	return &HTTP{addr: u.String(), mux: mux}
}

func (h *HTTP) Address() string {
	return h.addr
}

func (h *HTTP) Mux() *chi.Mux {
	return h.mux
}

// RPC

type RPC struct {
	addr string
	*grpc.Server
}

func NewRPC(addr string) *RPC {
	srv := grpc.NewServer()
	return &RPC{addr: addr, Server: srv}
}

func (r *RPC) Address() string {
	return r.addr
}
