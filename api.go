package main

import (
	"golang.org/x/image/bmp"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"

	"chartographer-go/chart"
)

// TODO директория должна указываться при инициализации приложения
// TODO должна создаваться, если она не существует.
var pathToFolder = "data/"

func createImage(w http.ResponseWriter, req *http.Request) {
	queryValues := req.URL.Query()
	width, err := strconv.Atoi(queryValues.Get("width"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	height, err := strconv.Atoi(queryValues.Get("height"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	img, err := chart.NewRGBA(width, height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgBytes, err := chart.Encode(img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := chart.Guid()
	name := chart.AppendExtension(id)

	err = os.WriteFile(pathToFolder+name, imgBytes, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func setFragment(w http.ResponseWriter, req *http.Request) {
	queryValues := req.URL.Query()
	x, err := strconv.Atoi(queryValues.Get("x"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(queryValues.Get("y"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	width, err := strconv.Atoi(queryValues.Get("width"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	height, err := strconv.Atoi(queryValues.Get("height"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fragment, err := bmp.Decode(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(req, "id")
	name := chart.AppendExtension(id)

	file, err := os.Open(pathToFolder + name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, err := bmp.Decode(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = file.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = chart.SetFragment(img, fragment, x, y, width, height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imgBytes, err := chart.Encode(img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(pathToFolder+name, imgBytes, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fragment(w http.ResponseWriter, req *http.Request) {
	queryValues := req.URL.Query()
	x, err := strconv.Atoi(queryValues.Get("x"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(queryValues.Get("y"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	width, err := strconv.Atoi(queryValues.Get("width"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	height, err := strconv.Atoi(queryValues.Get("height"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(req, "id")
	name := chart.AppendExtension(id)

	file, err := os.Open(pathToFolder + name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, err := bmp.Decode(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = file.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fragment, err := chart.Fragment(img, x, y, width, height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bmp.Encode(w, fragment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteImage(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	name := chart.AppendExtension(id)

	err := os.Remove(pathToFolder + name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
