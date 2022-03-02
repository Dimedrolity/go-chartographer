package main

import (
	"chartographer-go/bmp"
	"os"
	"strconv"

	"net/http"
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

	img, err := bmp.NewImage(width, height)
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
