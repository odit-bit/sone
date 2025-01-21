package monolith

import (
	"database/sql"

	"github.com/odit-bit/sone/pkg/kvstore"
	"github.com/odit-bit/sone/pkg/observer"
	"github.com/odit-bit/sone/pkg/rtmp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// monolith represent infrastructure code that module can use
type Monolith interface {
	Logger() *logrus.Logger

	// network
	RPC() *RPC
	RTMP() *rtmp.HandlerRegister
	HTTP() *HTTP

	// storage infrastructure
	FS() afero.Fs
	KV() *kvstore.Client
	DB() *sql.DB
	Minio() MinioConfig

	Observer() *observer.Observer
}
