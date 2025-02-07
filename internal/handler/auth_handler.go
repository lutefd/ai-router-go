package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/lutefd/ai-router-go/internal/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	oauthConfig *oauth2.Config
	authService service.AuthServiceInterface
	clientURL   string
}

func NewAuthHandler(authService service.AuthServiceInterface, clientID, clientSecret, redirectURL, clientURL string) *AuthHandler {
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthHandler{
		oauthConfig: oauthConfig,
		authService: authService,
		clientURL:   clientURL,
	}
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	originURL := r.URL.Query().Get("redirect_uri")
	if originURL == "" {
		originURL = h.clientURL
	}

	state := url.QueryEscape(originURL)
	url := h.oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(w, "State parameter missing", http.StatusInternalServerError)
		return
	}

	redirectURL, err := url.QueryUnescape(state)
	if err != nil {
		redirectURL = h.clientURL
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code parameter missing", http.StatusInternalServerError)
		return
	}

	token, err := h.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)
		return
	}

	userInfo, err := h.getUserInfo(r.Context(), token.AccessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}

	_, jwtToken, err := h.authService.AuthenticateUser(r.Context(), userInfo.Email, userInfo.Name, userInfo.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusInternalServerError)
		return
	}

	redirectWithToken := fmt.Sprintf("%s?token=%s", redirectURL, jwtToken)
	http.Redirect(w, r, redirectWithToken, http.StatusTemporaryRedirect)
}

type userInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *AuthHandler) getUserInfo(ctx context.Context, accessToken string) (*userInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var info userInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &info, nil
}
