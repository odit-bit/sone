package users

import (
	"context"

	"github.com/odit-bit/sone/internal/monolith"
	"github.com/odit-bit/sone/users/internal/application"
	"github.com/odit-bit/sone/users/internal/grpc"
	"github.com/odit-bit/sone/users/internal/sqlite"
	"github.com/odit-bit/sone/users/kv"
)

func StartModule(ctx context.Context, mono monolith.Monolith) {
	users, err := sqlite.New("users", mono.DB())
	if err != nil {
		mono.Logger().Panic(err)
	}
	streams := kv.New(mono.KV())
	glusers, err := sqlite.NewUserGlRepo(ctx, mono.DB(), "glusers")
	if err != nil {
		mono.Logger().Panic(err)
	}

	app := application.New(glusers, users, streams)
	grpc.RegisterUserService(app, mono.RPC())
}
