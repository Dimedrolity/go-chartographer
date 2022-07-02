package main

import (
	"log"
	"os"

	"github.com/Dimedrolity/go-chartographer/internal/app"
)

func main() {
	dataDirPath := os.Args[1]
	// TODO вынести хост и порт сервера и макс. размер тайла в os.Args
	const (
		port        = "8080"
		tileMaxSize = 1000
	)

	if err := app.Run(port, dataDirPath, tileMaxSize); err != nil {
		log.Fatal(err)
	}
}
