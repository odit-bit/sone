package ingress

import (
	"context"

	"github.com/odit-bit/sone/ingress/internal/app/commands"
	"github.com/odit-bit/sone/ingress/internal/localstore"
	"github.com/odit-bit/sone/ingress/internal/rpc"
	"github.com/odit-bit/sone/internal/monolith"
	"github.com/odit-bit/sone/pkg/rtmp"
	"github.com/odit-bit/sone/streaming/streamingpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitModule(ctx context.Context, mono monolith.Monolith) {

	conn, err := dial(ctx, mono.RPC().Address())
	if err != nil {
		mono.Logger().Panic(err)
	}
	streamsClient := streamingpb.NewLiveStreamClient(conn)
	app := localstore.NewHLSRepository(
		rpc.NewLiveStreams(streamsClient),
	)
	fn := func() rtmp.Handler {
		h := commands.NewCreateHLSHandler(ctx, app)
		return h
	}
	mono.RTMP().Set(fn)
}

func dial(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if err = conn.Close(); err != nil {
				// for logging
				return
			}
		}

		go func() {
			<-ctx.Done()
			if err = conn.Close(); err != nil {
				// for logging
			}
		}()
	}()

	return conn, nil
}
