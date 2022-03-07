package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"chartographer-go/chart"
)

func main() {
	pathToImages := os.Args[1]
	err := chart.SetImagesDir(pathToImages)
	if err != nil {
		log.Fatal(err)
	}

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
