package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
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
