package api

import (
	"context"

	"github.com/odit-bit/sone/streaming/internal/scheme"
	"github.com/odit-bit/sone/users/userpb"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	Users interface {
		AuthToken(ctx context.Context, token string) (*scheme.UserInfo, error)
	}
}

type UserRepo interface {
	AuthToken(ctx context.Context, token string) (*scheme.UserInfo, error)
}

func NewGrpcClient(conn *grpc.ClientConn) GrpcClient {
	return GrpcClient{
		Users: &UserGrpc{
			cli: userpb.NewUserServiceClient(conn),
		},
	}
}

var _ UserRepo = (*UserGrpc)(nil)

type UserGrpc struct {
	cli userpb.UserServiceClient
}

// AuthToken implements streaming.UserRepo.
func (u *UserGrpc) AuthToken(ctx context.Context, token string) (*scheme.UserInfo, error) {
	info, err := u.cli.AuthToken(ctx, &userpb.AuthTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}
	return &scheme.UserInfo{ID: info.ID, Name: info.Username}, nil
}

/////

func NewMockGRPC() GrpcClient {
	return GrpcClient{
		&mockUserGrpc{},
	}
}

type mockUserGrpc struct {
}

// AuthToken implements streaming.UserRepo.
func (u *mockUserGrpc) AuthToken(ctx context.Context, token string) (*scheme.UserInfo, error) {
	return &scheme.UserInfo{
		ID:   "123x123",
		Name: "testing",
	}, nil
}
