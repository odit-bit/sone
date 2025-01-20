package streaming

import (
	"context"

	"github.com/odit-bit/sone/internal/monolith"
	streamGRPC "github.com/odit-bit/sone/streaming/grpc"
	"github.com/odit-bit/sone/streaming/internal/api"
	"github.com/odit-bit/sone/streaming/internal/app"
	"github.com/odit-bit/sone/streaming/internal/database"
	"github.com/odit-bit/sone/streaming/streamingpb"
)

func StartModule(ctx context.Context, mono monolith.Monolith) error {
	conn, err := api.Dial(ctx, mono.RPC().Address())
	if err != nil {
		return err
	}

	//construct application neccessary component.
	repo, err := database.NewSqlite(mono.DB())
	if err != nil {
		return err
	}
	rpc := api.NewGrpcClient(conn)
	streaming := app.New(&repo, &rpc)

	// attach app into network interface (grpc, http, etc..)
	ss := streamGRPC.NewServer(streaming)
	streamingpb.RegisterLiveStreamServer(mono.RPC(), ss)
	//
	return nil
}
