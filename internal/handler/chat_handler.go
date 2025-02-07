package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lutefd/ai-router-go/internal/middleware"
	"github.com/lutefd/ai-router-go/internal/models"
	"github.com/lutefd/ai-router-go/internal/service"
)

type ChatHandler struct {
	chatService service.ChatServiceInterface
}

func NewChatHandler(chatService service.ChatServiceInterface) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var chat models.Chat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*service.Claims)
	chat.User = claims.UserID

	if err := h.chatService.CreateChat(r.Context(), &chat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}

func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	claims := r.Context().Value(middleware.UserContextKey).(*service.Claims)

	chat, err := h.chatService.GetChat(r.Context(), chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if chat.User != claims.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(chat)
}

func (h *ChatHandler) UpdateChatTitle(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	var update struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*service.Claims)

	chat, err := h.chatService.GetChat(r.Context(), chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if chat.User != claims.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chat.Title = update.Title
	if err := h.chatService.UpdateChat(r.Context(), chat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chat)
}

func (h *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	claims := r.Context().Value(middleware.UserContextKey).(*service.Claims)

	chat, err := h.chatService.GetChat(r.Context(), chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if chat.User != claims.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.chatService.DeleteChat(r.Context(), chatID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
