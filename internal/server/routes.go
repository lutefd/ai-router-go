package server

import (
	"github.com/go-chi/chi/v5"
)

func routes() chi.Router {

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {

	})

	return r

}
