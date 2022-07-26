package koiwai

import (
	"database/sql"
	"log"
	"mokou/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

//Connect creates a connection to postgres given
//a connection string.
func Connect(conf *config.KoiwaiConfig) (*bun.DB, *minio.Client) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.ConnectionString)))
	db := bun.NewDB(sqldb, pgdialect.New())

	s3Client, err := minio.New(conf.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.S3AccessKeyID, conf.S3SecretAccessKey, ""),
		Secure: conf.S3UseSSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

	return db, s3Client
}
