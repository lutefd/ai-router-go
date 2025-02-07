package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/lutefd/ai-router-go/internal/handler"
)

func routes(handler *handler.AIHandler) chi.Router {

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/ai", func(r chi.Router) {
			r.Route("/generate", func(r chi.Router) {
				r.Post("/", handler.ProxyRequest)
			})
		})
	})

	return r

}
