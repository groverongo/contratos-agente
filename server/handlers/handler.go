package handlers

import (
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type Handler struct {
	DB    *gorm.DB
	MinIO *minio.Client
}

func New(db *gorm.DB, minioClient *minio.Client) *Handler {
	return &Handler{
		DB:    db,
		MinIO: minioClient,
	}
}
