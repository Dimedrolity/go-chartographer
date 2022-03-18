package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"chartographer-go/chart"
)

type Config struct {
	Port string
}

func NewConfig(port string) *Config {
	return &Config{Port: port}
}

type Server struct {
	config *Config
	router *chi.Mux

	chartService *chart.ChartographerService
}

func NewServer(config *Config, chartService *chart.ChartographerService) *Server {
	s := &Server{
		config:       config,
		router:       chi.NewRouter(),
		chartService: chartService,
	}
	s.setRoutes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Run() error {
	p := ":" + s.config.Port
	return http.ListenAndServe(p, s)
}
