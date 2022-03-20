package server

import "github.com/go-chi/chi/v5"

func (s *Server) setRoutes() {
	s.router.Route("/chartas", func(r chi.Router) {
		r.Post("/", s.createImage)

		r.Route("/{id}", func(r chi.Router) {
			r.Post("/", s.setFragment)
			r.Get("/", s.getFragment)
			r.Delete("/", s.deleteImage)
		})
	})
}
