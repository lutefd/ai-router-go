package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lutefd/ai-router-go/internal/config"
	database "github.com/lutefd/ai-router-go/internal/database/mongodb"
	"github.com/lutefd/ai-router-go/internal/handler"
	"github.com/lutefd/ai-router-go/internal/middleware"
	"github.com/lutefd/ai-router-go/internal/repository"
	"github.com/lutefd/ai-router-go/internal/repository/mongodb"
	"github.com/lutefd/ai-router-go/internal/service"
	"github.com/lutefd/ai-router-go/internal/strategy"
	"github.com/lutefd/ai-router-go/pkg/idgen"
)

func Run() error {
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	conn, err := database.NewMongoDBConnection(ctx, cfg.MongoDBURI, cfg.MongoDBDatabase)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer conn.Close(ctx)

	if err := idgen.Init(cfg.WorkerID); err != nil {
		return fmt.Errorf("failed to initialize ID generator: %w", err)
	}

	geminiRepo := repository.NewGeminiRepository(ctx, cfg.GEMINI_SK)
	openaiRepo := repository.NewOpenAIRepository(cfg.OPENAI_SK)
	deepseekRepo := repository.NewDeepSeekRepository(cfg.DEEPSEEK_SK)
	userRepo := mongodb.NewUserRepository(conn.DB)
	aiService := service.NewAIService(geminiRepo, openaiRepo, deepseekRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	aiStrategy := strategy.NewAIStrategy(aiService)
	aiHandler := handler.NewAIHandler(aiStrategy)
	authHandler := handler.NewAuthHandler(authService, cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.AuthRedirectURL, cfg.ClientURL)
	chatRepo := mongodb.NewChatRepository(conn.DB)
	chatService := service.NewChatService(chatRepo)
	chatHandler := handler.NewChatHandler(chatService)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      routes(aiHandler, authHandler, chatHandler, userHandler, authMiddleware),
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(),
			30*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		log.Println("Server stopped")
	}()

	log.Printf("Server is listening on port %d\n", cfg.ServerPort)
	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		return fmt.Errorf("could not listen on %d: %w", cfg.ServerPort, err)
	}

	return nil
}
