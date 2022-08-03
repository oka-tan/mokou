package main

import (
	"context"
	"log"
	"mokou/asagi"
	"mokou/config"
	"mokou/importers"
	"mokou/koiwai"
	"time"

	"github.com/minio/minio-go/v7"
)

func main() {
	log.Println("Starting mokou. Kill this right now if you haven't manually created the bucket or the table partitions.")

	conf, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	asagiDb := asagi.Connect(conf.AsagiConfig.ConnectionString)
	koiwaiDb, koiwaiS3Client := koiwai.Connect(conf.PostgresConfig, conf.S3Config)

	bucketExists, err := koiwaiS3Client.BucketExists(context.Background(), conf.S3Config.S3BucketName)

	if err != nil {
		log.Fatal(err)
	}

	if !bucketExists {
		err = koiwaiS3Client.MakeBucket(context.Background(), conf.S3Config.S3BucketName, minio.MakeBucketOptions{
			Region:        conf.S3Config.S3Region,
			ObjectLocking: false,
		})

		if err != nil {
			log.Fatal("Error creating S3 media bucket")
		}
	}

	koiwaiS3Service := koiwai.S3Service{
		KoiwaiDb:   koiwaiDb,
		S3Client:   koiwaiS3Client,
		BucketName: conf.S3Config.S3BucketName,
	}

	importerService := importers.Service{
		AsagiDb:           asagiDb,
		AsagiImagesFolder: *conf.AsagiConfig.ImagesFolder,
		KoiwaiDb:          koiwaiDb,
		KoiwaiS3Service:   &koiwaiS3Service,
		BatchSize:         conf.PostgresConfig.BatchSize,
	}

	log.Println("Beginning migration.")

	for _, boardConfig := range conf.Boards {
		if err = importerService.AsagiToKoiwai(&boardConfig); err != nil {
			panic(err)
		}
	}

	log.Println("Fully vacuuming all tables manually and rebuilding indexes.")

	if _, err := koiwaiDb.QueryContext(context.Background(), "VACUUM (FULL, ANALYZE) post"); err != nil {
		log.Println("Error vacuuming all post tables. Consider executing 'VACUUM (FULL, ANALYZE) post' manually in psql.")
		log.Println("Migration went alright, however.")
	}

	if _, err := koiwaiDb.QueryContext(context.Background(), "VACUUM (FULL, ANALYZE) media"); err != nil {
		log.Println("Error vacuuming the media table. Consider executing 'VACUUM (FULL, ANALYZE) media' manually in psql.")
		log.Println("Migration went alright, however.")
	}

	log.Println("Done.")
}
