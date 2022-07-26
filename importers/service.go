package importers

import (
	"mokou/koiwai"

	"github.com/uptrace/bun"
)

//Service imports data from schema to schema.
type Service struct {
	AsagiDb           *bun.DB
	AsagiImagesFolder string
	KoiwaiDb          *bun.DB
	KoiwaiS3Service   *koiwai.S3Service
	BatchSize         int
}
