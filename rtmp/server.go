package rtmp

import (
	"context"
	"io"
	"net"

	"github.com/livekit/go-rtmp"
	"github.com/odit-bit/sone/internal/kvstore"
	"github.com/odit-bit/sone/rtmp/internal/app/commands"
	"github.com/odit-bit/sone/rtmp/internal/localstore"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func NewRTMPServer(ctx context.Context, logger *logrus.Logger, rootDir string, afs afero.Fs) *rtmp.Server {
	// ingest := fs.NewHLSRepository(rootDir, afs)

	rs := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			// handler will create every time received connection (dispatch)
			// h := handler.New(ctx, ingest, handler.WithBitRateLimit(handler.DefaultMaxBitrate))
			ingest := commands.NewCreateHLSHandler(ctx, localstore.NewHLSRepository(rootDir, afs, kvstore.Open()))
			conf := &rtmp.ConnConfig{
				Handler: ingest,
				ControlState: rtmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},
				Logger: logger,
			}

			return conn, conf
		},
	})
	return rs

}
