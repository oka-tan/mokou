package main

import (
	"context"
	"log"
	"mokou/asagi"
	"mokou/config"
	"mokou/importers"
	"mokou/koiwai"
	"time"
)

func main() {
	log.Println("Starting mokou. Kill this right now if you haven't manually created the bucket or the table partitions.")

	conf, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	initialNap, err := time.ParseDuration(conf.InitialNap)
	if err == nil {
		time.Sleep(initialNap)
	} else {
		time.Sleep(5 * time.Second)
	}

	asagiDb := asagi.Connect(conf.AsagiConfig.ConnectionString)
	koiwaiDb, koiwaiS3Client := koiwai.Connect(conf.KoiwaiConfig)
	koiwaiS3Service := koiwai.S3Service{
		KoiwaiDb:   koiwaiDb,
		S3Client:   koiwaiS3Client,
		BucketName: conf.KoiwaiConfig.S3BucketName,
	}

	importerService := importers.Service{
		AsagiDb:           asagiDb,
		AsagiImagesFolder: *conf.AsagiConfig.ImagesFolder,
		KoiwaiDb:          koiwaiDb,
		KoiwaiS3Service:   &koiwaiS3Service,
		BatchSize:         conf.BatchSize,
	}

	log.Println("Beginning migration. Hopefully you've disabled autovacuum and created.")

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
