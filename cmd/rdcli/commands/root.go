package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:     "rdcli",
	Short:   "Reddit CLI + MCP",
	Long:    "Reddit CLI + MCP - browse, post, search Reddit from your terminal.",
	Version: "0.1.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as raw JSON")

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(meCmd)
	rootCmd.AddCommand(feedCmd)
	rootCmd.AddCommand(postsCmd)
	rootCmd.AddCommand(commentsCmd)
	rootCmd.AddCommand(subredditsCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(messagesCmd)
	rootCmd.AddCommand(serveCmd)
}
