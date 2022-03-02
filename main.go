package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()

	router.Post("/chartas/", createImage)

	log.Fatal(http.ListenAndServe(":8080", router))
}
