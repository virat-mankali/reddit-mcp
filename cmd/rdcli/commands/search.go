package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/virat-mankali/reddit-mcp/internal/api"
	"github.com/virat-mankali/reddit-mcp/internal/output"
)

var (
	searchType  string
	searchSub   string
	searchSort  string
	searchTime  string
	searchLimit int
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search Reddit",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		opts := api.SearchOptions{
			Query:     args[0],
			Type:      searchType,
			Subreddit: searchSub,
			Sort:      searchSort,
			Time:      searchTime,
			Limit:     searchLimit,
		}

		switch searchType {
		case "", "posts":
			items, err := client.SearchPosts(cmd.Context(), opts)
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(items)
			}
			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{
					"r/" + item.Subreddit,
					trimBody(item.Title, 70),
					"u/" + item.Author,
					fmt.Sprintf("%d", item.Score),
				})
			}
			output.PrintTable([]string{"Subreddit", "Title", "Author", "Score"}, rows)
		case "comments":
			items, err := client.SearchComments(cmd.Context(), opts)
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(items)
			}
			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{
					"r/" + item.Subreddit,
					"u/" + item.Author,
					fmt.Sprintf("%d", item.Score),
					trimBody(item.Body, 90),
				})
			}
			output.PrintTable([]string{"Subreddit", "Author", "Score", "Body"}, rows)
		case "subreddits":
			items, err := client.SearchSubreddits(cmd.Context(), opts)
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
		case "users":
			items, err := client.SearchUsers(cmd.Context(), opts)
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(items)
			}
			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{
					"u/" + item.Name,
					fmt.Sprintf("%d", item.LinkKarma),
					fmt.Sprintf("%d", item.CommentKarma),
				})
			}
			output.PrintTable([]string{"User", "Link Karma", "Comment Karma"}, rows)
		default:
			return fmt.Errorf("unsupported search type %q", searchType)
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchType, "type", "posts", "Search type: posts, comments, subreddits, users")
	searchCmd.Flags().StringVar(&searchSub, "sub", "", "Restrict search to a subreddit")
	searchCmd.Flags().StringVar(&searchSort, "sort", "relevance", "Search sort: relevance, new, top, comments")
	searchCmd.Flags().StringVar(&searchTime, "time", "all", "Time filter: hour, day, week, month, year, all")
	searchCmd.Flags().IntVar(&searchLimit, "limit", 10, "Number of results to fetch")
}
