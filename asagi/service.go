package asagi

import (
	"github.com/minio/minio-go/v7"
	"github.com/uptrace/bun"
)

//Service imports data from the asagi schema
type Service struct {
	Pg *bun.DB

	S3MediaClient      *minio.Client
	S3MediaBucket      *string
	S3ThumbnailsClient *minio.Client
	S3ThumbnailsBucket *string
	S3OekakiClient     *minio.Client
	S3OekakiBucket     *string

	AsagiDb           *bun.DB
	AsagiImagesFolder string
	BatchSize         int
}
