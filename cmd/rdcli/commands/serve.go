package commands

import (
	"fmt"

	mcpserver "github.com/virat-mankali/reddit-mcp/internal/mcp"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Reddit MCP server (stdio)",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := mcpserver.NewServer()
		if err := server.ServeStdio(s); err != nil {
			return fmt.Errorf("mcp server error: %w", err)
		}
		return nil
	},
}
