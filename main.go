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
	log.Println("Starting mokou.")

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
	koiwaiDb := koiwai.Connect(conf.KoiwaiConfig.ConnectionString)

	importerService := importers.Service{
		AsagiDb:   asagiDb,
		KoiwaiDb:  koiwaiDb,
		BatchSize: conf.BatchSize,
	}

	log.Println("Beginning migration. Hopefully you've disabled autovacuum.")

	for _, boardConfig := range conf.Boards {
		if err = importerService.AsagiToKoiwai(&boardConfig); err != nil {
			panic(err)
		}
	}

	log.Println("Fully vacuuming all tables manually and rebuilding indexes.")

	if _, err = koiwaiDb.QueryContext(context.Background(), "VACUUM (FULL, ANALYZE) post"); err != nil {
		log.Println("Error vacuuming all tables. Consider executing 'VACUUM (FULL, ANALYZE) post' manually in psql.")
		log.Println("Migration went alright, however.")
		panic(err)
	}

	log.Println("Done.")
}
