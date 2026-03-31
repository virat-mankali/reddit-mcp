package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/virat-mankali/reddit-mcp/internal/auth"
	"github.com/virat-mankali/reddit-mcp/internal/config"
)

const baseURL = "https://oauth.reddit.com"

type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
	auth       *auth.Manager
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		userAgent:  cfg.UserAgent,
		auth:       auth.NewManager(cfg),
	}
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	token, err := auth.LoadToken()
	if err != nil {
		return nil, fmt.Errorf("load auth token: %w", err)
	}
	if token.NeedsRefresh() {
		token, err = c.auth.Refresh(ctx, token)
		if err != nil {
			return nil, err
		}
	}

	req = req.Clone(ctx)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("reddit api %s: %s", resp.Status, string(body))
	}

	return resp, nil
}
