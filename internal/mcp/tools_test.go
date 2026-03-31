package mcpserver

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestHandleGetMeReturnsAuthErrorWhenUnauthenticated(t *testing.T) {
	t.Setenv("REDDIT_CLIENT_ID", "")
	t.Setenv("REDDIT_CLIENT_SECRET", "")
	t.Setenv("REDDIT_ACCESS_TOKEN", "")
	t.Setenv("REDDIT_REFRESH_TOKEN", "")

	result, err := handleGetMe(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("unexpected protocol error: %v", err)
	}
	if !result.IsError {
		t.Fatalf("expected tool error result")
	}
}

func TestHandleListPostsRequiresSubreddit(t *testing.T) {
	result, err := handleListPosts(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("unexpected protocol error: %v", err)
	}
	if !result.IsError {
		t.Fatalf("expected tool error result")
	}
}
