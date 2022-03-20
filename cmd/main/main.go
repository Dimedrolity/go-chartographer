package main

import (
	"log"
	"os"

	"go-chartographer/internal/app"
)

func main() {
	// TODO вынести хост и порт сервера и макс. размер тайла в os.Args
	const port = "8080"
	dataDirPath := os.Args[1]
	const tileMaxSize = 1000
	if err := app.Run(port, dataDirPath, tileMaxSize); err != nil {
		log.Fatal(err)
	}
}
