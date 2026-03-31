package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/virat-mankali/reddit-mcp/internal/output"
)

var subredditLimit int

var subredditsCmd = &cobra.Command{
	Use:   "subreddits",
	Short: "Browse and manage subreddits",
}

var subredditsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your subscribed subreddits",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		items, err := client.ListMySubreddits(cmd.Context(), subredditLimit)
		if err != nil {
			return err
		}
		if jsonOutput {
			return printJSON(items)
		}
		rows := make([][]string, 0, len(items))
		for _, item := range items {
			rows = append(rows, []string{
				"r/" + item.DisplayName,
				fmt.Sprintf("%d", item.Subscribers),
				trimBody(item.PublicDescription, 80),
			})
		}
		output.PrintTable([]string{"Subreddit", "Subscribers", "Description"}, rows)
		return nil
	},
}

var subredditsInfoCmd = &cobra.Command{
	Use:   "info <subreddit>",
	Short: "Get info about a subreddit",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		item, err := client.GetSubreddit(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		if jsonOutput {
			return printJSON(item)
		}
		fmt.Printf("Name        : r/%s\n", item.DisplayName)
		fmt.Printf("Title       : %s\n", item.Title)
		fmt.Printf("Subscribers : %d\n", item.Subscribers)
		fmt.Printf("NSFW        : %t\n", item.Over18)
		fmt.Printf("Subscribed  : %t\n", item.UserIsSubscriber)
		fmt.Printf("Created     : %s\n", item.CreatedAt().Format("January 2, 2006"))
		if item.PublicDescription != "" {
			fmt.Printf("Description : %s\n", item.PublicDescription)
		}
		return nil
	},
}

var subredditsSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search subreddits",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		items, err := client.SearchSubredditsByName(cmd.Context(), args[0], subredditLimit)
		if err != nil {
			return err
		}
		if jsonOutput {
			return printJSON(items)
		}
		rows := make([][]string, 0, len(items))
		for _, item := range items {
			rows = append(rows, []string{
				"r/" + item.DisplayName,
				fmt.Sprintf("%d", item.Subscribers),
				trimBody(item.PublicDescription, 90),
			})
		}
		output.PrintTable([]string{"Subreddit", "Subscribers", "Description"}, rows)
		return nil
	},
}

func subredditActionCommand(use, short string, subscribe bool) *cobra.Command {
	return &cobra.Command{
		Use:   use + " <subreddit>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newAPIClient()
			if err != nil {
				return err
			}
			if err := client.Subscribe(cmd.Context(), args[0], subscribe); err != nil {
				return err
			}
			if subscribe {
				fmt.Printf("Subscribed to %s\n", args[0])
			} else {
				fmt.Printf("Unsubscribed from %s\n", args[0])
			}
			return nil
		},
	}
}

func init() {
	subredditsCmd.AddCommand(subredditsListCmd, subredditsInfoCmd, subredditsSearchCmd, subredditActionCommand("subscribe", "Subscribe to a subreddit", true), subredditActionCommand("unsubscribe", "Unsubscribe from a subreddit", false))
	subredditsListCmd.Flags().IntVar(&subredditLimit, "limit", 25, "Number of subreddits to fetch")
	subredditsSearchCmd.Flags().IntVar(&subredditLimit, "limit", 25, "Number of subreddits to fetch")
}
