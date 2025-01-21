package media

import (
	"context"

	"github.com/odit-bit/sone/internal/monolith"
	"github.com/odit-bit/sone/media/internal/api"
	"github.com/odit-bit/sone/media/internal/application"
	"github.com/odit-bit/sone/media/internal/domain"
	mediaGRPC "github.com/odit-bit/sone/media/internal/grpc"
	"github.com/odit-bit/sone/media/internal/segment"
)

func StartModule(mono monolith.Monolith) {
	var segmentRepo domain.SegmentRepository
	var err error

	switch mono.MediaStorageVendor() {
	case monolith.FS:
		segmentRepo = segment.NewLocalFS(mono.FS())
		mono.Logger().Debug("media module : initiate segment repository in memory")

	case monolith.MINIO:
		segmentRepo, err = segment.NewMiniStore(context.Background(), "sone", mono.Minio())
		if err != nil {
			mono.Logger().Fatal("error initiate segment repository:", err)
		}
		mono.Logger().Debug("media module : initiate segment repository with minio")
	}

	app := application.New(segmentRepo)
	if err := mediaGRPC.RegisterServer(app, mono.RPC()); err != nil {
		mono.Logger().Fatal("error register media grpc-server, THIS IS BUG!!")
	}

	h := api.NewHLSHandler(app, mono.Logger())
	mux := mono.HTTP().Mux()
	mono.HTTP().Mux().Mount("/segment", mux)
	h.Register(mux)

}
