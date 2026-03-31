package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jsonOutput bool
var (
	buildVersion = "0.1.0"
	buildCommit  = "none"
	buildDate    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "rdcli",
	Short: "Reddit CLI + MCP",
	Long:  "Reddit CLI + MCP - browse, post, search Reddit from your terminal.",
}

func Execute() {
	rootCmd.Version = buildVersion
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func SetBuildInfo(version, commit, date string) {
	if version != "" {
		buildVersion = version
	}
	if commit != "" {
		buildCommit = commit
	}
	if date != "" {
		buildDate = date
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
