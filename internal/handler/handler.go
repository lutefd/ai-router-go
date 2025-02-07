package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/lutefd/ai-router-go/internal/strategy"
)

type AIHandler struct {
	aiStrategy strategy.AIStrategyInterface
}

func NewAIHandler(aiStrategy strategy.AIStrategyInterface) *AIHandler {
	return &AIHandler{aiStrategy: aiStrategy}
}

func (h *AIHandler) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	platform := r.Header.Get("Platform")
	if platform == "" {
		http.Error(w, "Platform header is required", http.StatusBadRequest)
		return
	}

	model := r.Header.Get("Model")
	if model == "" {
		http.Error(w, "Model header is required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	err = h.aiStrategy.GenerateResponse(r.Context(), platform, model,
		string(body), func(chunk string) {
			fmt.Fprintf(w, "data: %s\n\n", chunk)
			flusher.Flush()
		})

	if err != nil {
		log.Printf("Error generating response: %v", err)
		fmt.Fprintf(w, "data: ERROR: %s\n\n", err.Error())
		flusher.Flush()
		return
	}

	fmt.Fprint(w, "data: [DONE]\n\n")
	flusher.Flush()
	log.Println("Stream completed")
}
