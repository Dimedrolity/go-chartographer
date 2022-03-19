package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"golang.org/x/image/bmp"

	"chartographer-go/chart"
	"chartographer-go/store"
	"chartographer-go/tiledimage"
)

func paramError(name string, err error) error {
	return fmt.Errorf(
		"некорректный параметр запроса - %v: %w", name, err)
}
func getQueryParam(req *http.Request, name string) (string, error) {
	q := req.URL.Query()
	if !q.Has(name) || q.Get(name) == "" {
		return "", paramError(name, errors.New("отсутствует или пустая строка"))
	}
	return q.Get(name), nil
}
func getQueryParamInt(req *http.Request, name string) (int, error) {
	s, err := getQueryParam(req, name)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, paramError(name, err)
	}
	return i, nil
}

func (s *Server) createImage(w http.ResponseWriter, req *http.Request) {
	width, err := getQueryParamInt(req, "width")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	height, err := getQueryParamInt(req, "height")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	img, err := s.chartService.NewRgbaBmp(width, height)
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

func (s *Server) setFragment(w http.ResponseWriter, req *http.Request) {
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

	err = s.chartService.SetFragment(id, fragment)
	if err != nil {
		if errors.Is(err, tiledimage.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func (s *Server) getFragment(w http.ResponseWriter, req *http.Request) {
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

	img, err := s.chartService.GetTiledImage(id)
	if err != nil {
		if errors.Is(err, tiledimage.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	fragment, err := s.chartService.GetFragment(img, x, y, width, height)
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

func (s *Server) deleteImage(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	err := s.chartService.DeleteImage(id)
	if err != nil {
		if errors.Is(err, tiledimage.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
