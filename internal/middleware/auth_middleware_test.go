package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lutefd/ai-router-go/internal/middleware"
	"github.com/lutefd/ai-router-go/internal/mocks"
	"github.com/lutefd/ai-router-go/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthMiddleware_RequireAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)
	authMiddleware := middleware.NewAuthMiddleware(mockAuthService)

	tests := []struct {
		name           string
		setupAuth      func(r *http.Request)
		setupMocks     func()
		expectedStatus int
		checkContext   bool
	}{
		{
			name: "valid token",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer valid-token")
			},
			setupMocks: func() {
				mockAuthService.EXPECT().
					ValidateToken("valid-token").
					Return(&service.Claims{
						UserID: "123",
						Email:  "test@example.com",
						Role:   "user",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkContext:   true,
		},
		{
			name:      "missing auth header",
			setupAuth: func(r *http.Request) {},
			setupMocks: func() {
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid auth format",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "InvalidFormat token")
			},
			setupMocks: func() {
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid token",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid-token")
			},
			setupMocks: func() {
				mockAuthService.EXPECT().
					ValidateToken("invalid-token").
					Return(nil, fmt.Errorf("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			handler := authMiddleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkContext {
					claims := r.Context().Value(middleware.UserContextKey)
					assert.True(t, claims != nil)
					assert.NotNil(t, claims.(*service.Claims))
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			tt.setupAuth(req)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
