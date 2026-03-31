package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	userAgent := cfg.UserAgent
	if userAgent == "" {
		userAgent = "reddit-mcp/1.0"
	}

	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		userAgent:  userAgent,
		auth:       auth.NewManager(cfg),
	}
}

func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, body io.Reader) (*http.Request, error) {
	if query == nil {
		query = url.Values{}
	}
	if _, ok := query["raw_json"]; !ok {
		query.Set("raw_json", "1")
	}

	target := c.baseURL + path
	if encoded := query.Encode(); encoded != "" {
		target += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, method, target, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	token, err := auth.LoadToken()
	if err != nil {
		return nil, fmt.Errorf("load auth token: %w", err)
	}
	if token.NeedsRefresh() && token.RefreshToken != "" {
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

	if err := applyRateLimit(resp); err != nil {
		resp.Body.Close()
		return nil, err
	}

	if resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("reddit api %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	return resp, nil
}

func (c *Client) getJSON(ctx context.Context, path string, query url.Values, out any) error {
	req, err := c.newRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return err
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) postForm(ctx context.Context, path string, form url.Values, out any) error {
	body := bytes.NewBufferString(form.Encode())
	req, err := c.newRequest(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if out == nil {
		_, err = io.Copy(io.Discard, resp.Body)
		return err
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func applyRateLimit(resp *http.Response) error {
	remainingHeader := resp.Header.Get("X-Ratelimit-Remaining")
	resetHeader := resp.Header.Get("X-Ratelimit-Reset")
	if remainingHeader == "" || resetHeader == "" {
		return nil
	}

	remaining, err := strconv.ParseFloat(remainingHeader, 64)
	if err != nil {
		return nil
	}
	if remaining > 0 {
		return nil
	}

	resetSeconds, err := strconv.ParseFloat(resetHeader, 64)
	if err != nil {
		return nil
	}
	sleepFor := time.Duration(math.Ceil(resetSeconds)) * time.Second
	if sleepFor <= 0 {
		return nil
	}
	if sleepFor > 10*time.Second {
		sleepFor = 10 * time.Second
	}
	time.Sleep(sleepFor)
	return nil
}
