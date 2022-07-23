package main

import (
	"mokou/asagi"
	"mokou/config"
	"mokou/eientei"
	"mokou/importers"
	"time"
)

func main() {
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
	eienteiDb := eientei.Connect(conf.EienteiConfig.ConnectionString)

	importerService := importers.Service{
		AsagiDb:   asagiDb,
		EienteiDb: eienteiDb,
		BatchSize: conf.BatchSize,
	}

	for _, boardConfig := range conf.Boards {
		if err = importerService.AsagiToEientei(&boardConfig); err != nil {
			panic(err)
		}
	}
}
