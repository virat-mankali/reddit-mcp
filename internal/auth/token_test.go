package auth

import "testing"

func TestLoadTokenFromEnv(t *testing.T) {
	t.Setenv("REDDIT_ACCESS_TOKEN", "access-token")
	t.Setenv("REDDIT_REFRESH_TOKEN", "refresh-token")
	t.Setenv("REDDIT_USERNAME", "virat")

	token := loadTokenFromEnv()
	if token == nil {
		t.Fatal("expected token from env")
	}
	if token.AccessToken != "access-token" {
		t.Fatalf("unexpected access token: %q", token.AccessToken)
	}
	if token.RefreshToken != "refresh-token" {
		t.Fatalf("unexpected refresh token: %q", token.RefreshToken)
	}
	if token.Username != "virat" {
		t.Fatalf("unexpected username: %q", token.Username)
	}
}
