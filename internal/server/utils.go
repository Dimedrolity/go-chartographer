package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
