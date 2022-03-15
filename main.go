package main

import (
	"chartographer-go/chart"
	"chartographer-go/store"
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
	chart.TileRepo = tileRepo

	store.TileMaxSize = 1000
	chart.ImageRepo = store.NewInMemoryImageRepo()

	router := chi.NewRouter()

	router.Route("/chartas", func(r chi.Router) {
		r.Post("/", createImage)

		r.Route("/{id}", func(r chi.Router) {
			//r.Post("/", setFragment)
			r.Get("/", fragment)
			r.Delete("/", deleteImage)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
