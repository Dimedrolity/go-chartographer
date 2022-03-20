package server

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go-chartographer/internal/chart"
)

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

	img, err := s.chartService.AddImage(width, height)

	var errSize *chart.SizeError
	if err != nil {
		if errors.As(err, &errSize) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(img.Id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) setFragment(w http.ResponseWriter, req *http.Request) {
	x, err := getQueryParamInt(req, "x")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	y, err := getQueryParamInt(req, "y")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// пусть width и height будут обязательными параметрами, несмотря на то, что
	// размеры можно получить при декодировании изображения в теле запроса
	_, err = getQueryParamInt(req, "width")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = getQueryParamInt(req, "height")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(req, "id")

	img, err := s.chartService.GetImage(id)
	if err != nil {
		if errors.Is(err, chart.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO не декодировать сразу, сначала проверить, что есть пересечение img и width height
	b, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fragment, err := s.chartService.Decode(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.chartService.SetFragment(img, x, y, fragment)
	if err != nil {
		if errors.Is(err, chart.ErrNotOverlaps) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getFragment(w http.ResponseWriter, req *http.Request) {
	x, err := getQueryParamInt(req, "x")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	y, err := getQueryParamInt(req, "y")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	id := chi.URLParam(req, "id")

	img, err := s.chartService.GetImage(id)
	if err != nil {
		if errors.Is(err, chart.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fragment, err := s.chartService.GetFragment(img, x, y, width, height)

	var errSize *chart.SizeError
	if err != nil {
		if errors.As(err, &errSize) || errors.Is(err, chart.ErrNotOverlaps) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := s.chartService.Encode(fragment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "image/bmp")
}

func (s *Server) deleteImage(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	err := s.chartService.DeleteImage(id)
	if err != nil {
		if errors.Is(err, chart.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
