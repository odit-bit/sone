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

const (
	// media storage value
	FS = iota // default
	MINIO
)

// monolith represent infrastructure code that module can use
type Monolith interface {
	Logger() *logrus.Logger
	RPC() *RPC
	RTMP() *rtmp.HandlerRegister
	HTTP() *HTTP

	FS() afero.Fs
	KV() *kvstore.Client
	DB() *sql.DB
	Minio() *minio.Client

	//media storage type
	MediaStorageVendor() int

	Observer() *observer.Observer
}
