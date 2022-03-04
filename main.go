package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()

	router.Post("/chartas/", createImage)
	router.Post("/chartas/{id}/", setFragment)
	router.Get("/chartas/{id}/", fragment)
	router.Delete("/chartas/{id}/", deleteImage)

	log.Fatal(http.ListenAndServe(":8080", router))
}
