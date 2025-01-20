package rtmp

import (
	"io"
	"net"

	"github.com/livekit/go-rtmp"
)

func New(h rtmp.Handler) *rtmp.Server {
	srv := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(c net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			conf := rtmp.ConnConfig{
				Handler: h,
			}
			return c, &conf
		},
	})
	return srv
}
