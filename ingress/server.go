package ingress

import (
	"context"

	"github.com/odit-bit/sone/ingress/internal/app/commands"
	"github.com/odit-bit/sone/ingress/internal/domain"
	"github.com/odit-bit/sone/ingress/internal/localstore"
	"github.com/odit-bit/sone/pkg/rtmp"
	"github.com/spf13/afero"
)

func InitModule(ctx context.Context, rootDir string, afs afero.Fs, kv domain.StreamKeyRepository, hr *rtmp.HandlerRegister) {
	// handlerConf.OnConnect = onConnect(ctx, logger, rootDir, afs, kv)
	hlsRepo := localstore.NewHLSRepository(rootDir, afs, kv)
	fn := func() rtmp.Handler {
		h := commands.NewCreateHLSHandler(ctx, hlsRepo)
		return h
	}
	hr.Set(fn)
}
