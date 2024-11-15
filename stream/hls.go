package stream

import (
	"net/http"

	"github.com/odit-bit/sone/internal/kvstore"
	"github.com/odit-bit/sone/stream/internal/fs"
	"github.com/odit-bit/sone/stream/internal/handler"
	"github.com/spf13/afero"

	"github.com/sirupsen/logrus"
)

func NewHLSServer(logger *logrus.Logger, afs afero.Fs) *http.Server {
	kv := kvstore.Open()
	fs, err := fs.New(afs, kv)
	if err != nil {
		panic(err)
	}

	fs.SetLogger(logger)
	ui := handler.UI()
	fh := handler.NewWebHandler(fs, ui, fs)

	srv := http.Server{Handler: fh.Handler()}
	return &srv
}
