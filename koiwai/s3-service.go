package koiwai

import (
	"github.com/minio/minio-go/v7"
	"github.com/uptrace/bun"
)

type S3Service struct {
	S3Client   *minio.Client
	KoiwaiDb   *bun.DB
	BucketName string
}
