package grpc

import (
	"context"
	"fmt"

	"github.com/odit-bit/sone/users/gluserpb"
	"github.com/odit-bit/sone/users/internal/application"
	"github.com/odit-bit/sone/users/internal/domain"
	"github.com/odit-bit/sone/users/userpb"
)

// RegisterUser implements userpb.UserServiceServer.
func (s *Server) RegisterUser(ctx context.Context, in *userpb.RegisterUserRequest) (*userpb.RegisterUserResponse, error) {

	resp, err := s.app.RegisterUser(ctx, application.RegisterUserArgs{
		Username: in.Username,
		Password: in.Password,
	})
	return &userpb.RegisterUserResponse{Id: resp.Id}, err
}

func (s *Server) AuthenticateUser(ctx context.Context, in *userpb.AuthUserRequest) (*userpb.AuthUserResponse, error) {
	token, err := s.app.AuthenticateUser(ctx, application.AuthenticateUserArgs{
		Username: in.Username,
		Password: in.Password,
	})
	return &userpb.AuthUserResponse{Token: token}, err
}

func (s *Server) AuthToken(ctx context.Context, in *userpb.AuthTokenRequest) (*userpb.UserInfo, error) {
	resp, err := s.app.AuthToken(ctx, in.Token)
	if err != nil {
		return nil, err
	}
	return &userpb.UserInfo{ID: resp.ID, Username: resp.Username, Email: resp.Email}, nil
}

func (s *Server) Get(ctx context.Context, in *gluserpb.GetRequest) (*gluserpb.GetResponse, error) {
	var resp gluserpb.GetResponse
	user, ok, err := s.app.GetUserGoogle(ctx, in.Id)
	if err != nil {
		return nil, fmt.Errorf("google user service error %v", err)
	}
	if !ok {
		return &resp, nil
	}
	resp.Id = user.Id
	resp.Email = user.Email
	resp.Name = user.Name
	resp.IsEmailVerified = user.IsEmailVerified
	resp.IsFound = true

	return &resp, nil
}

func (s *Server) Save(ctx context.Context, in *gluserpb.SaveRequest) (*gluserpb.SaveResponse, error) {
	err := s.app.SaveUserGoogle(ctx, domain.UserGl{
		Id:              in.Id,
		Name:            in.Name,
		Email:           in.Email,
		IsEmailVerified: in.IsEmailVerified,
	})

	if err != nil {
		return nil, err
	}
	return &gluserpb.SaveResponse{}, nil
}
