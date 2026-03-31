package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/virat-mankali/reddit-mcp/internal/output"
)

var (
	feedSort  string
	feedTime  string
	feedLimit int
)

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Browse your home feed",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		posts, err := client.Feed(cmd.Context(), feedSort, feedTime, feedLimit)
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(posts)
		}

		rows := make([][]string, 0, len(posts))
		for i, post := range posts {
			rows = append(rows, []string{
				fmt.Sprintf("%d", i+1),
				fmt.Sprintf("%d", post.Score),
				"r/" + post.Subreddit,
				trimBody(post.Title, 70),
				"u/" + post.Author,
			})
		}
		output.PrintTable([]string{"#", "Score", "Subreddit", "Title", "Author"}, rows)
		return nil
	},
}

func init() {
	feedCmd.Flags().StringVar(&feedSort, "sort", "best", "Feed sort: best, hot, new, top")
	feedCmd.Flags().StringVar(&feedTime, "time", "day", "Time filter for top sort: hour, day, week, month, year, all")
	feedCmd.Flags().IntVar(&feedLimit, "limit", 10, "Number of posts to fetch")
}
