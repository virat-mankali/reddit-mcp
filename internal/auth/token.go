package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/zalando/go-keyring"
)

const (
	KeyringService = "reddit-mcp"
	tokenKey       = "token"
)

type TokenStore struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	Username     string    `json:"username"`
}

func SaveToken(token *TokenStore) error {
	payload, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}

	if err := keyring.Set(KeyringService, tokenKey, string(payload)); err != nil {
		return fmt.Errorf("save token in keyring: %w", err)
	}

	return nil
}

func LoadToken() (*TokenStore, error) {
	if token := loadTokenFromEnv(); token != nil {
		return token, nil
	}

	payload, err := keyring.Get(KeyringService, tokenKey)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("load token from keyring: %w", err)
	}

	var token TokenStore
	if err := json.Unmarshal([]byte(payload), &token); err != nil {
		return nil, fmt.Errorf("decode token: %w", err)
	}

	return &token, nil
}

func HasEnvToken() bool {
	return loadTokenFromEnv() != nil
}

func ClearToken() error {
	err := keyring.Delete(KeyringService, tokenKey)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("delete token from keyring: %w", err)
	}
	return nil
}

func (t *TokenStore) NeedsRefresh() bool {
	if t == nil {
		return false
	}
	return time.Now().After(t.ExpiresAt.Add(-5 * time.Minute))
}

func loadTokenFromEnv() *TokenStore {
	accessToken := os.Getenv("REDDIT_ACCESS_TOKEN")
	refreshToken := os.Getenv("REDDIT_REFRESH_TOKEN")
	if accessToken == "" && refreshToken == "" {
		return nil
	}

	token := &TokenStore{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(55 * time.Minute),
		Username:     os.Getenv("REDDIT_USERNAME"),
	}

	if accessToken == "" {
		token.ExpiresAt = time.Now().Add(-time.Minute)
	}

	return token
}
