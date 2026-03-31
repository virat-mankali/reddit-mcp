package commands

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/virat-mankali/reddit-mcp/internal/api"
	"github.com/virat-mankali/reddit-mcp/internal/config"
	"github.com/virat-mankali/reddit-mcp/internal/output"
)

func newAPIClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if cfg.ClientID == "" && cfg.ClientSecret == "" {
		return nil, errors.New("missing Reddit app configuration; run `rdcli auth login` or set REDDIT_CLIENT_ID and REDDIT_CLIENT_SECRET")
	}
	return api.NewClient(cfg), nil
}

func printJSON(v any) error {
	return output.PrintJSON(v)
}

func trimBody(text string, limit int) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	if len(text) <= limit {
		return text
	}
	if limit <= 3 {
		return text[:limit]
	}
	return text[:limit-3] + "..."
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format("2006-01-02 15:04")
}

func boolVote(v *bool) string {
	if v == nil {
		return ""
	}
	if *v {
		return "up"
	}
	return "down"
}

func fullname(id, kind string) string {
	id = strings.TrimSpace(id)
	if strings.HasPrefix(id, "t1_") || strings.HasPrefix(id, "t3_") || strings.HasPrefix(id, "t4_") || strings.HasPrefix(id, "t5_") {
		return id
	}
	return kind + "_" + id
}

func requireArgs(n int, got []string) error {
	if len(got) < n {
		return fmt.Errorf("expected at least %d args", n)
	}
	return nil
}
