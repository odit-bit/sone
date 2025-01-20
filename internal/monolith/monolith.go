package monolith

import (
	"database/sql"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"github.com/odit-bit/sone/pkg/kvstore"
	"github.com/odit-bit/sone/pkg/observer"
	"github.com/odit-bit/sone/pkg/rtmp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
)

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

type HTTP struct {
	addr string
}

func NewHTTP(addr string) *HTTP {
	u := url.URL{}
	u.Scheme = "http"
	u.Host = addr

	return &HTTP{addr: u.String()}
}

func (h *HTTP) Address() string {
	return h.addr
}

// monolith represent infrastructure code that module can use
type Monolith interface {
	FS() afero.Fs
	KV() *kvstore.Client
	Logger() *logrus.Logger
	RPC() *RPC
	Mux() *chi.Mux
	RTMP() *rtmp.HandlerRegister
	DB() *sql.DB
	Minio() *minio.Client
	Test() bool
	HTTP() *HTTP
	Observer() *observer.Observer
}
