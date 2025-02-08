package server

import (
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lutefd/ai-router-go/internal/handler"
	"github.com/lutefd/ai-router-go/internal/middleware"
)

func routes(handler *handler.AIHandler, authHandler *handler.AuthHandler, chatHandler *handler.ChatHandler, userHandler *handler.UserHandler, healthHandler *handler.HealthHandler, authMiddleware *middleware.AuthMiddleware) chi.Router {

	r := chi.NewRouter()

	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Refresh-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/healthz", healthHandler.LivenessCheck)
	r.Get("/readiness", healthHandler.ReadinessCheck)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Get("/google/login", authHandler.GoogleLogin)
			r.Get("/google/callback", authHandler.GoogleCallback)
			r.Post("/google/native/signin", authHandler.HandleNativeSignIn)
			r.Post("/google/refresh", authHandler.RefreshToken)
		})

		r.Route("/ai", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Route("/generate", func(r chi.Router) {
				r.Post("/", handler.ProxyRequest)
			})
		})

		r.Route("/chats", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Post("/", chatHandler.CreateChat)
			r.Get("/{id}", chatHandler.GetChat)
			r.Put("/{id}/title", chatHandler.UpdateChatTitle)
			r.Delete("/{id}", chatHandler.DeleteChat)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Get("/me/chats", userHandler.GetUserChats)
		})
	})

	return r

}
