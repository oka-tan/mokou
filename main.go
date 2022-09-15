package main

import (
	"context"
	"database/sql"
	"log"
	"mokou/asagi"
	"mokou/badger"
	"mokou/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	log.Println("Starting Mokou")

	conf := config.LoadConfig()

	time.Sleep(5 * time.Second)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.PostgresConfig.ConnectionString)))
	pg := bun.NewDB(sqldb, pgdialect.New())

	var s3MediaClient *minio.Client
	var s3MediaBucket *string
	var s3ThumbnailsClient *minio.Client
	var s3ThumbnailsBucket *string

	if conf.MediaConfig != nil {
		s3MediaBucket = &conf.MediaConfig.S3BucketName

		var err error
		s3MediaClient, err = minio.New(conf.MediaConfig.S3Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(conf.MediaConfig.S3AccessKeyID, conf.MediaConfig.S3SecretAccessKey, ""),
			Secure: conf.MediaConfig.S3UseSSL,
		})

		if err != nil {
			log.Fatalf("Error creating media connection: %s", err)
		}

		bucketExists, err := s3MediaClient.BucketExists(context.Background(), conf.MediaConfig.S3BucketName)

		if err != nil {
			log.Fatalf("Error checking if media bucket exists: %s", err)
		}

		if !bucketExists {
			err := s3MediaClient.MakeBucket(context.Background(), conf.MediaConfig.S3BucketName, minio.MakeBucketOptions{
				Region:        conf.MediaConfig.S3Region,
				ObjectLocking: false,
			})

			if err != nil {
				log.Fatalf("Error creating S3 media bucket: %s", err)
			}
		}
	}

	if conf.ThumbnailsConfig != nil {
		s3ThumbnailsBucket = &conf.ThumbnailsConfig.S3BucketName

		var err error
		s3ThumbnailsClient, err = minio.New(conf.ThumbnailsConfig.S3Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(conf.ThumbnailsConfig.S3AccessKeyID, conf.ThumbnailsConfig.S3SecretAccessKey, ""),
			Secure: conf.ThumbnailsConfig.S3UseSSL,
		})

		if err != nil {
			log.Fatalf("Error creating thumbnails connection: %s", err)
		}

		bucketExists, err := s3ThumbnailsClient.BucketExists(context.Background(), conf.ThumbnailsConfig.S3BucketName)

		if err != nil {
			log.Fatalf("Error checking if media bucket exists: %s", err)
		}

		if !bucketExists {
			err := s3ThumbnailsClient.MakeBucket(context.Background(), conf.ThumbnailsConfig.S3BucketName, minio.MakeBucketOptions{
				Region:        conf.ThumbnailsConfig.S3Region,
				ObjectLocking: false,
			})

			if err != nil {
				log.Fatalf("Error creating S3 thumbnails bucket: %s", err)
			}
		}
	}

	log.Println("Beginning migration.")

	if conf.AsagiConfig != nil {
		asagiSqldb, err := sql.Open("mysql", conf.AsagiConfig.ConnectionString)

		if err != nil {
			log.Fatalf("Error connecting to Asagi's mysql: %s", err)
		}

		asagiDb := bun.NewDB(asagiSqldb, mysqldialect.New())

		importerService := asagi.Service{
			Pg: pg,

			S3MediaClient:      s3MediaClient,
			S3MediaBucket:      s3MediaBucket,
			S3ThumbnailsClient: s3ThumbnailsClient,
			S3ThumbnailsBucket: s3ThumbnailsBucket,

			AsagiDb:           asagiDb,
			AsagiImagesFolder: *conf.AsagiConfig.ImagesFolder,
			BatchSize:         conf.PostgresConfig.BatchSize,
		}

		for _, boardConfig := range conf.AsagiConfig.Boards {
			if err = importerService.Import(&boardConfig); err != nil {
				log.Fatalf("Error importing board %s: %s", boardConfig.Name, err)
			}
		}
	}

	if conf.BadgerConfig != nil {
		importerService := badger.Service{
			Pg:         pg,
			JsonFolder: conf.BadgerConfig.JsonFolder,
		}

		for _, boardConfig := range conf.BadgerConfig.Boards {
			importerService.Import(boardConfig)
		}
	}

	log.Printf("Done\n")
}
