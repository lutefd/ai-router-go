package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/lutefd/ai-router-go/internal/models"
	"github.com/lutefd/ai-router-go/internal/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	oauthConfig *oauth2.Config
	authService service.AuthServiceInterface
	clientURL   string
}

type TokenPair = service.TokenPair

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

	_, tokenPair, err := h.authService.AuthenticateUser(r.Context(), userInfo.Email, userInfo.Name, userInfo.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusInternalServerError)
		return
	}

	isMobile := r.URL.Query().Get("platform") == "mobile"
	if isMobile {
		scheme := r.URL.Query().Get("app_scheme")
		redirectURL = fmt.Sprintf("%s://oauth/callback?access_token=%s&refresh_token=%s&expires_in=%d",
			scheme,
			tokenPair.AccessToken,
			tokenPair.RefreshToken,
			tokenPair.ExpiresIn)
	} else {
		redirectURL = fmt.Sprintf("%s?access_token=%s&refresh_token=%s&expires_in=%d",
			redirectURL,
			tokenPair.AccessToken,
			tokenPair.RefreshToken,
			tokenPair.ExpiresIn)
	}

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("X-Refresh-Token")
	if refreshToken == "" {
		http.Error(w, "Refresh token required", http.StatusBadRequest)
		return
	}

	tokenPair, err := h.authService.RefreshAccessToken(refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
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

type GoogleIDTokenClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Sub           string `json:"sub"`
}

func (h *AuthHandler) HandleNativeSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		IDToken string `json:"id_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenInfo, err := h.verifyGoogleIDToken(r.Context(), req.IDToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid ID token: %v", err), http.StatusUnauthorized)
		return
	}

	user, tokenPair, err := h.authService.AuthenticateUser(
		r.Context(),
		tokenInfo.Email,
		tokenInfo.Name,
		tokenInfo.Sub,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		User      *models.User `json:"user"`
		TokenPair *TokenPair   `json:"tokens"`
	}{
		User:      user,
		TokenPair: tokenPair,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) verifyGoogleIDToken(ctx context.Context, idToken string) (*GoogleIDTokenClaims, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: h.oauthConfig.ClientID,
	})

	token, err := verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	var claims GoogleIDTokenClaims
	if err := token.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return &claims, nil
}
