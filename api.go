package main

import (
	"chartographer-go/chart"
	"chartographer-go/store"
	"chartographer-go/tiledimage"
	"errors"
	"golang.org/x/image/bmp"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

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

	img, err := chart.NewRgbaBmp(width, height)
	var errSize *chart.SizeError
	if err != nil {
		if errors.As(err, &errSize) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(img.Id))
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
	//width, err := strconv.Atoi(queryValues.Get("width"))
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//height, err := strconv.Atoi(queryValues.Get("height"))
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	id := chi.URLParam(req, "id")

	// TODO не декодировать сразу, сначала проверить, что есть пересечение img и width height
	fragment, err := bmp.Decode(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	store.ShiftRect(fragment, x, y)

	err = chart.SetFragment2(id, fragment)
	if err != nil {
		if errors.Is(err, tiledimage.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	img, err := chart.ImageRepo.Get(id)
	if err != nil {
		if errors.Is(err, tiledimage.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	fragment, err := chart.GetFragment(img, x, y, width, height)
	var errSize *chart.SizeError
	if err != nil {
		if errors.As(err, &errSize) || errors.Is(err, chart.ErrNotOverlaps) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	err := chart.DeleteImage(id)
	if err != nil {
		if errors.Is(err, tiledimage.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
