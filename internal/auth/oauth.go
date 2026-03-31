package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/virat-mankali/reddit-mcp/internal/config"
)

const (
	AuthURL     = "https://www.reddit.com/api/v1/authorize"
	TokenURL    = "https://www.reddit.com/api/v1/access_token"
	RedirectURI = "http://localhost:3141/callback"
	Scopes      = "identity read submit vote history privatemessages mysubreddits save subscribe"
)

type Manager struct {
	cfg        *config.Config
	httpClient *http.Client
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Error        string `json:"error"`
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (m *Manager) Login(ctx context.Context) (*TokenStore, error) {
	state, err := randomString(24)
	if err != nil {
		return nil, err
	}
	verifier, err := randomString(64)
	if err != nil {
		return nil, err
	}
	challenge := pkceChallenge(verifier)

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	server := &http.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Query().Get("error") != "":
			http.Error(w, "Reddit authorization failed. You can close this window.", http.StatusBadRequest)
			errCh <- fmt.Errorf("reddit returned error: %s", r.URL.Query().Get("error"))
		case r.URL.Query().Get("state") != state:
			http.Error(w, "State mismatch. You can close this window.", http.StatusBadRequest)
			errCh <- errors.New("oauth state mismatch")
		case r.URL.Query().Get("code") == "":
			http.Error(w, "Missing code. You can close this window.", http.StatusBadRequest)
			errCh <- errors.New("missing authorization code")
		default:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = io.WriteString(w, "<html><body><h3>Reddit authentication complete.</h3><p>You can close this window and return to the terminal.</p></body></html>")
			codeCh <- r.URL.Query().Get("code")
		}
	})
	server.Handler = mux

	ln, err := net.Listen("tcp", "127.0.0.1:3141")
	if err != nil {
		return nil, fmt.Errorf("listen on callback port: %w", err)
	}
	defer ln.Close()

	go func() {
		if serveErr := server.Serve(ln); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errCh <- serveErr
		}
	}()
	defer server.Shutdown(context.Background())

	authURL, err := buildAuthURL(m.cfg.ClientID, state, challenge)
	if err != nil {
		return nil, err
	}
	if err := openBrowser(authURL); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case code := <-codeCh:
		return m.exchangeCode(ctx, code, verifier)
	}
}

func (m *Manager) Refresh(ctx context.Context, token *TokenStore) (*TokenStore, error) {
	values := url.Values{}
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", token.RefreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, TokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", basicAuth(m.cfg.ClientID, m.cfg.ClientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", m.userAgent())

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("refresh token: %s", strings.TrimSpace(string(body)))
	}

	var parsed tokenResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}
	if parsed.Error != "" {
		return nil, fmt.Errorf("refresh token: %s", parsed.Error)
	}

	refreshed := &TokenStore{
		AccessToken:  parsed.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(parsed.ExpiresIn) * time.Second),
		Username:     token.Username,
	}
	if parsed.RefreshToken != "" {
		refreshed.RefreshToken = parsed.RefreshToken
	}

	if err := SaveToken(refreshed); err != nil {
		return nil, err
	}
	return refreshed, nil
}

func (m *Manager) exchangeCode(ctx context.Context, code, verifier string) (*TokenStore, error) {
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("redirect_uri", RedirectURI)
	values.Set("code_verifier", verifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, TokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", basicAuth(m.cfg.ClientID, m.cfg.ClientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", m.userAgent())

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("exchange code: %s", strings.TrimSpace(string(body)))
	}

	var parsed tokenResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}
	if parsed.Error != "" {
		return nil, fmt.Errorf("exchange code: %s", parsed.Error)
	}

	token := &TokenStore{
		AccessToken:  parsed.AccessToken,
		RefreshToken: parsed.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(parsed.ExpiresIn) * time.Second),
	}
	if err := SaveToken(token); err != nil {
		return nil, err
	}
	return token, nil
}

func buildAuthURL(clientID, state, challenge string) (string, error) {
	u, err := url.Parse(AuthURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("client_id", clientID)
	q.Set("response_type", "code")
	q.Set("state", state)
	q.Set("redirect_uri", RedirectURI)
	q.Set("duration", "permanent")
	q.Set("scope", Scopes)
	q.Set("code_challenge", challenge)
	q.Set("code_challenge_method", "S256")
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func openBrowser(target string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target)
	case "linux":
		cmd = exec.Command("xdg-open", target)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	default:
		return fmt.Errorf("unsupported platform for browser open; visit this URL manually: %s", target)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("open browser: %w", err)
	}
	return nil
}

func randomString(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func basicAuth(clientID, clientSecret string) string {
	raw := clientID + ":" + clientSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}

func (m *Manager) userAgent() string {
	if m.cfg.UserAgent != "" {
		return m.cfg.UserAgent
	}
	return "reddit-mcp/1.0"
}
