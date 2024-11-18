package streaming

import (
	"github.com/go-chi/chi/v5"
	"github.com/odit-bit/sone/pkg/kvstore"
	"github.com/odit-bit/sone/streaming/internal/application"
	"github.com/odit-bit/sone/streaming/internal/httphandler"
	"github.com/odit-bit/sone/streaming/internal/localfs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func StartUp(logger *logrus.Logger, afs afero.Fs, kv *kvstore.Client, mux *chi.Mux) {
	fs := localfs.New(afs, kv)
	app := application.New(fs, fs, fs)
	httphandler.Register(app, mux)
}
