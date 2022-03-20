package main

import (
	"chartographer-go/server"
	"log"
	"os"

	"chartographer-go/chart"
	"chartographer-go/imagetile"
	"chartographer-go/kvstore"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	pathToImages := os.Args[1]
	tileRepo, err := imagetile.NewFileSystemTileRepo(pathToImages)
	if err != nil {
		return err
	}
	bmpService := imagetile.NewBmpService(tileRepo)
	tileMaxSize := 1000
	imageRepo := kvstore.NewInMemoryStore()
	// хорошо бы вписался дженерик при инициализации imageRepo, необходимо указать *chart.TiledImage
	chartService := chart.NewChartographerService(imageRepo, bmpService, tileMaxSize)
	// TODO вынести хост и порт в .env
	config := server.NewConfig("8080")
	srv := server.NewServer(config, chartService)

	return srv.Run()
}
