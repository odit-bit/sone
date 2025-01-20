package grpc

import (
	"github.com/odit-bit/sone/users/gluserpb"
	"github.com/odit-bit/sone/users/internal/application"
	"github.com/odit-bit/sone/users/userpb"
	"google.golang.org/grpc"
)

var _ userpb.UserServiceServer = (*Server)(nil)
var _ gluserpb.GoogleUserServiceServer = (*Server)(nil)

type Server struct {
	app *application.Application
	userpb.UnimplementedUserServiceServer
	gluserpb.UnimplementedGoogleUserServiceServer
}

func New(app *application.Application) *Server {
	s := Server{
		app:                                  app,
		UnimplementedUserServiceServer:       userpb.UnimplementedUserServiceServer{},
		UnimplementedGoogleUserServiceServer: gluserpb.UnimplementedGoogleUserServiceServer{},
	}
	return &s
}

func RegisterUserService(app *application.Application, registrar grpc.ServiceRegistrar) {
	s := Server{
		app:                                  app,
		UnimplementedUserServiceServer:       userpb.UnimplementedUserServiceServer{},
		UnimplementedGoogleUserServiceServer: gluserpb.UnimplementedGoogleUserServiceServer{},
	}
	userpb.RegisterUserServiceServer(registrar, &s)
	gluserpb.RegisterGoogleUserServiceServer(registrar, &s)
}
