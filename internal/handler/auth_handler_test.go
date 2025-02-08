package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/lutefd/ai-router-go/internal/mocks"
	"github.com/lutefd/ai-router-go/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
)

func TestAuthHandler_GoogleLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)
	handler := NewAuthHandler(mockAuthService, "client-id", "client-secret", "http://localhost:8080/callback", "http://localhost:3000", "client-id")

	tests := []struct {
		name             string
		redirectURI      string
		expectedLocation string
	}{
		{
			name:             "default redirect",
			redirectURI:      "",
			expectedLocation: "https://accounts.google.com/o/oauth2/auth",
		},
		{
			name:             "custom redirect",
			redirectURI:      "http://custom-client.com",
			expectedLocation: "https://accounts.google.com/o/oauth2/auth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/auth/google/login", nil)
			if tt.redirectURI != "" {
				q := req.URL.Query()
				q.Add("redirect_uri", tt.redirectURI)
				req.URL.RawQuery = q.Encode()
			}

			rr := httptest.NewRecorder()
			handler.GoogleLogin(rr, req)

			require.Equal(t, http.StatusTemporaryRedirect, rr.Code)
			location := rr.Header().Get("Location")
			assert.Contains(t, location, tt.expectedLocation)
		})
	}
}

func TestAuthHandler_GoogleCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/token":
			json.NewEncoder(w).Encode(&oauth2.Token{
				AccessToken: "test-access-token",
				TokenType:   "Bearer",
			})
		case "/userinfo":
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Bearer test-access-token" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    "123",
				"email": "test@example.com",
				"name":  "Test User",
			})
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))
	defer ts.Close()

	handler := NewAuthHandler(
		mockAuthService,
		"client-id",
		"client-secret",
		"http://localhost:8080/callback",
		"http://localhost:3000",
		"client-id",
	)

	handler.oauthConfig = &oauth2.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  ts.URL + "/auth",
			TokenURL: ts.URL + "/token",
		},
	}

	transport := &http.Transport{}
	client := &http.Client{Transport: transport}

	oldClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: &mockTransport{
			base: transport,
			fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.String() == "https://www.googleapis.com/oauth2/v2/userinfo" {
					req.URL, _ = url.Parse(ts.URL + "/userinfo")
				}
				return client.Transport.RoundTrip(req)
			},
		},
	}
	defer func() { http.DefaultClient = oldClient }()

	tests := []struct {
		name           string
		setupMocks     func()
		state          string
		code           string
		expectedStatus int
	}{
		{
			name: "successful callback",
			setupMocks: func() {
				mockAuthService.EXPECT().
					AuthenticateUser(
						gomock.Any(),
						"test@example.com",
						"Test User",
						"123",
					).
					Return(&models.User{
						ID:    "123",
						Email: "test@example.com",
						Name:  "Test User",
						Role:  "user",
					}, &TokenPair{
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						ExpiresIn:    900,
					}, nil)
			},
			state:          url.QueryEscape("http://localhost:3000"),
			code:           "valid-code",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "missing code",
			setupMocks:     func() {},
			state:          url.QueryEscape("http://localhost:3000"),
			code:           "",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "missing state",
			setupMocks:     func() {},
			state:          "",
			code:           "valid-code",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			reqURL := fmt.Sprintf("/auth/google/callback?code=%s&state=%s", tt.code, tt.state)
			req := httptest.NewRequest(http.MethodGet, reqURL, nil)
			rr := httptest.NewRecorder()

			handler.GoogleCallback(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusTemporaryRedirect {
				location := rr.Header().Get("Location")
				expectedURL := "http://localhost:3000?access_token=access-token&refresh_token=refresh-token&expires_in=900"
				assert.Equal(t, expectedURL, location)
			}
		})
	}
}

type mockTransport struct {
	base http.RoundTripper
	fn   func(*http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.fn != nil {
		return t.fn(req)
	}
	if t.base != nil {
		return t.base.RoundTrip(req)
	}
	return http.DefaultTransport.RoundTrip(req)
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)
	handler := NewAuthHandler(mockAuthService, "client-id", "client-secret", "http://localhost:8080/callback", "http://localhost:3000", "client-id")

	tests := []struct {
		name           string
		setupMocks     func()
		refreshToken   string
		expectedStatus int
		expectedBody   *TokenPair
	}{
		{
			name: "successful refresh",
			setupMocks: func() {
				mockAuthService.EXPECT().
					RefreshAccessToken("valid-refresh-token").
					Return(&TokenPair{
						AccessToken:  "new-access-token",
						RefreshToken: "new-refresh-token",
						ExpiresIn:    900,
					}, nil)
			},
			refreshToken:   "valid-refresh-token",
			expectedStatus: http.StatusOK,
			expectedBody: &TokenPair{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
				ExpiresIn:    900,
			},
		},
		{
			name:           "missing refresh token",
			setupMocks:     func() {},
			refreshToken:   "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "invalid refresh token",
			setupMocks: func() {
				mockAuthService.EXPECT().
					RefreshAccessToken("invalid-token").
					Return(nil, fmt.Errorf("invalid token"))
			},
			refreshToken:   "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
			if tt.refreshToken != "" {
				req.Header.Set("X-Refresh-Token", tt.refreshToken)
			}
			rr := httptest.NewRecorder()

			handler.RefreshToken(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != nil {
				var response TokenPair
				err := json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, &response)
			}
		})
	}
}
