package commands

import "github.com/spf13/cobra"

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show your Reddit profile",
}

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Browse your home feed",
}

var postsCmd = &cobra.Command{
	Use:   "posts",
	Short: "Browse and manage posts",
}

var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Browse and manage comments",
}

var subredditsCmd = &cobra.Command{
	Use:   "subreddits",
	Short: "Browse and manage subreddits",
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Reddit",
}

var messagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "Read and send Reddit messages",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server",
}
