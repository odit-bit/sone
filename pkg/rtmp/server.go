package rtmp

import (
	"io"
	"net"

	gortmp "github.com/livekit/go-rtmp"
)

type Server struct {
	rtmpSrv *gortmp.Server
	Handler *HandlerRegister
}

func NewServer() *Server {
	return &Server{
		rtmpSrv: &gortmp.Server{},
	}
}

func (s *Server) Serve(l net.Listener) error {
	srv := gortmp.NewServer(&gortmp.ServerConfig{
		OnConnect: func(c net.Conn) (io.ReadWriteCloser, *gortmp.ConnConfig) {
			conf := gortmp.ConnConfig{
				Handler: s.Handler.handlerFunc(),
			}
			return c, &conf
		},
	})
	s.rtmpSrv = srv
	return s.rtmpSrv.Serve(l)
}

func (s *Server) Close() error {
	return s.rtmpSrv.Close()
}

type HandlerFunc func() Handler
type Handler gortmp.Handler

type HandlerRegister struct {
	handlerFunc HandlerFunc
}

func NewHandler() *HandlerRegister {
	return &HandlerRegister{}
}

func (h *HandlerRegister) Set(fn HandlerFunc) {
	h.handlerFunc = fn
}
