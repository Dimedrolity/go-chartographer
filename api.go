package main

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"

	"chartographer-go/bmp"
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

	img, err := bmp.NewRGBA(width, height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgBytes, err := bmp.Encode(img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := bmp.Guid()
	name := bmp.AppendExtension(id)

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

	fragmentBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fragment, err := bmp.Decode(fragmentBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(req, "id")
	name := bmp.AppendExtension(id)

	imgBytes, err := os.ReadFile(pathToFolder + name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, err := bmp.Decode(imgBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bmp.SetFragment(img, fragment, x, y, width, height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imgBytes, err = bmp.Encode(img)
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
