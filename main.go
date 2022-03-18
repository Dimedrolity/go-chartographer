package main

import (
	"log"
	"os"

	"chartographer-go/chart"
	"chartographer-go/store"
	"chartographer-go/tiledimage"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	pathToImages := os.Args[1]
	tileRepo, err := store.NewFileSystemTileRepo(pathToImages)
	if err != nil {
		return err
	}
	tileMaxSize := 1000
	imageRepo := tiledimage.NewInMemoryImageRepo()
	chartService := chart.NewChartographerService(imageRepo, tileRepo, tileMaxSize)
	// TODO вынести хост и порт в .env
	config := NewConfig("8080")
	server := NewServer(config, chartService)

	return server.Run()
}
