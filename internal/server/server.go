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

	"github.com/lutefd/ai-router/internal/config"
	"google.golang.org/genai"
)

func run() error {
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	geminiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.GEMINI_SK,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      routes(),
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
