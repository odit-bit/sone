package app

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/odit-bit/sone/internal/monolith"
)

func initMinio(conf *Config) (*minio.Client, error) {
	if conf.MediaStorage != monolith.MINIO {
		return &minio.Client{}, nil
	}
	
	minioUrl := conf.Minio.Address
	minioID := conf.Minio.AccessKey
	minioSecret := conf.Minio.SecretKey
	mioCli, err := minio.New(minioUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(minioID, minioSecret, ""),
		Secure: false,
	})
	return mioCli, err
}
