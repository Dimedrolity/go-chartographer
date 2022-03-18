package main

import (
	"chartographer-go/chart"
	"chartographer-go/store"
	"chartographer-go/tiledimage"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	pathToImages := os.Args[1]
	tileRepo, err := store.NewFileSystemTileRepo(pathToImages)
	if err != nil {
		log.Fatal(err)
	}

	imageRepo := tiledimage.NewInMemoryImageRepo()
	tileMaxSize := 1000

	chartService := chart.NewChartographerService(imageRepo, tileRepo, tileMaxSize)

	ChartService = chartService

	router := chi.NewRouter()

	router.Route("/chartas", func(r chi.Router) {
		r.Post("/", createImage)

		r.Route("/{id}", func(r chi.Router) {
			r.Post("/", setFragment)
			r.Get("/", fragment)
			r.Delete("/", deleteImage)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
