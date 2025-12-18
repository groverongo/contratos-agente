package common

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMinio(config *ServerConfig) (*minio.Client, error) {
	useSSL := false

	minioClient, err := minio.New(config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Minio.AccessKey, config.Minio.SecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("failed to connect minio: %v", err)
	}

	return minioClient, nil
}
