package media

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/odit-bit/sone/internal/monolith"
	"github.com/odit-bit/sone/media/internal/api"
	"github.com/odit-bit/sone/media/internal/application"
	"github.com/odit-bit/sone/media/internal/domain"
	mediaGRPC "github.com/odit-bit/sone/media/internal/grpc"
	"github.com/odit-bit/sone/media/internal/segment"
)

func StartModule(mono monolith.Monolith) {
	var segmentRepo domain.SegmentRepository

	conf := mono.Minio()
	if conf.Addr != "" {
		minioCli, err := initMinio(conf)
		if err != nil {
			mono.Logger().Fatal("media module: initiate minio:", err)
		}
		segmentRepo, err = segment.NewMiniStore(context.Background(), "sone", minioCli)
		if err != nil {
			mono.Logger().Fatal("error initiate segment repository:", err)
		}
		mono.Logger().Debug("media module : initiate segment repository with minio")
	} else {
		mono.Logger().Debug("media module : minio address is not given, try use filesystem")
		segmentRepo = segment.NewLocalFS(mono.FS())
		mono.Logger().Debug("media module : use filesystem for segment repository ")
	}

	app := application.New(segmentRepo)
	if err := mediaGRPC.RegisterServer(app, mono.RPC()); err != nil {
		mono.Logger().Fatal("error register media grpc-server, THIS IS BUG!!")
	}

	h := api.NewHLSHandler(app, mono.Logger())
	mux := mono.HTTP().Mux()
	mux.Mount("/segment", mux)
	h.Register(mux)

}

func initMinio(conf monolith.MinioConfig) (*minio.Client, error) {
	// if conf.MediaStorage != monolith.MINIO {
	// 	return &minio.Client{}, nil
	// }

	minioUrl := conf.Addr
	minioID := conf.AccessKey
	minioSecret := conf.SecretKey
	mioCli, err := minio.New(minioUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(minioID, minioSecret, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	_, err = mioCli.HealthCheck(5 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed connect to minio: %v ", err)
	}
	if ok := mioCli.IsOnline(); !ok {
		return nil, fmt.Errorf("minio health check failed")
	}

	return mioCli, nil
}
