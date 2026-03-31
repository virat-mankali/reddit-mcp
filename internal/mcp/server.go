package mcpserver

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/virat-mankali/reddit-mcp/internal/api"
	"github.com/virat-mankali/reddit-mcp/internal/config"
)

func NewServer() *server.MCPServer {
	s := server.NewMCPServer(
		"reddit-mcp",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	registerTools(s)
	return s
}

func newClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if cfg.ClientID == "" && cfg.ClientSecret == "" {
		return nil, errors.New("not authenticated. Run: rdcli auth login")
	}
	return api.NewClient(cfg), nil
}

func toolError(err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(err.Error()), nil
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultJSON(v)
}

func requiredString(req mcp.CallToolRequest, key string) (string, *mcp.CallToolResult, error) {
	value, err := req.RequireString(key)
	if err != nil {
		return "", mcp.NewToolResultError(err.Error()), nil
	}
	return value, nil, nil
}

func listLimit(req mcp.CallToolRequest, defaultLimit int) int {
	limit := req.GetInt("limit", defaultLimit)
	if limit <= 0 {
		return defaultLimit
	}
	return limit
}

func withClient(ctx context.Context, fn func(*api.Client) (*mcp.CallToolResult, error)) (*mcp.CallToolResult, error) {
	client, err := newClient()
	if err != nil {
		return toolError(err)
	}
	return fn(client)
}
