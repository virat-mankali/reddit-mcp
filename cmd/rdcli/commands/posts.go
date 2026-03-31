package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/virat-mankali/reddit-mcp/internal/output"
)

var (
	postsListSort  string
	postsListTime  string
	postsListLimit int
	postComments   int

	submitTitle string
	submitText  string
	submitURL   string
	submitType  string
)

var postsCmd = &cobra.Command{
	Use:   "posts",
	Short: "Browse and manage posts",
}

var postsListCmd = &cobra.Command{
	Use:   "list <subreddit>",
	Short: "List posts from a subreddit",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		posts, err := client.ListPosts(cmd.Context(), args[0], postsListSort, postsListTime, postsListLimit)
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
				trimBody(post.Title, 72),
				"u/" + post.Author,
				fmt.Sprintf("%d", post.NumComments),
			})
		}
		output.PrintTable([]string{"#", "Score", "Title", "Author", "Comments"}, rows)
		return nil
	},
}

var postsGetCmd = &cobra.Command{
	Use:   "get <postId>",
	Short: "Get a post and its comments",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		details, err := client.GetPost(cmd.Context(), args[0], postComments)
		if err != nil {
			return err
		}
		if jsonOutput {
			return printJSON(details)
		}
		fmt.Printf("Title     : %s\n", details.Post.Title)
		fmt.Printf("Subreddit : r/%s\n", details.Post.Subreddit)
		fmt.Printf("Author    : u/%s\n", details.Post.Author)
		fmt.Printf("Score     : %d\n", details.Post.Score)
		fmt.Printf("Comments  : %d\n", details.Post.NumComments)
		if details.Post.URL != "" {
			fmt.Printf("URL       : %s\n", details.Post.URL)
		}
		if details.Post.SelfText != "" {
			fmt.Printf("\n%s\n", details.Post.SelfText)
		}
		if len(details.Comments) == 0 {
			return nil
		}
		fmt.Println("\nComments")
		rows := make([][]string, 0, len(details.Comments))
		for _, comment := range details.Comments {
			rows = append(rows, []string{
				strings.Repeat("  ", comment.Depth) + "u/" + comment.Author,
				fmt.Sprintf("%d", comment.Score),
				trimBody(comment.Body, 100),
			})
		}
		output.PrintTable([]string{"Author", "Score", "Body"}, rows)
		return nil
	},
}

var postsSubmitCmd = &cobra.Command{
	Use:   "submit <subreddit>",
	Short: "Submit a post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		if submitTitle == "" {
			return fmt.Errorf("--title is required")
		}
		if submitType == "link" && submitURL == "" {
			return fmt.Errorf("--url is required for link posts")
		}
		if submitType != "link" && submitText == "" {
			return fmt.Errorf("--text is required for text posts")
		}
		resp, err := client.SubmitPost(cmd.Context(), args[0], submitTitle, submitText, submitURL, submitType)
		if err != nil {
			return err
		}
		if jsonOutput {
			return printJSON(resp)
		}
		fmt.Printf("Submitted post %s to r/%s\n", resp.Name, resp.Subreddit)
		if resp.URL != "" {
			fmt.Printf("URL: %s\n", resp.URL)
		}
		return nil
	},
}

var postsDeleteCmd = postActionCommand("delete", "Delete a post", func(client APIActions, ctx *cobra.Command, id string) error {
	return client.DeleteThing(ctx.Context(), fullname(id, "t3"))
})
var postsUpvoteCmd = postActionCommand("upvote", "Upvote a post", func(client APIActions, ctx *cobra.Command, id string) error {
	return client.Vote(ctx.Context(), fullname(id, "t3"), 1)
})
var postsDownvoteCmd = postActionCommand("downvote", "Downvote a post", func(client APIActions, ctx *cobra.Command, id string) error {
	return client.Vote(ctx.Context(), fullname(id, "t3"), -1)
})
var postsSaveCmd = postActionCommand("save", "Save a post", func(client APIActions, ctx *cobra.Command, id string) error {
	return client.Save(ctx.Context(), fullname(id, "t3"), true)
})
var postsUnsaveCmd = postActionCommand("unsave", "Unsave a post", func(client APIActions, ctx *cobra.Command, id string) error {
	return client.Save(ctx.Context(), fullname(id, "t3"), false)
})

type APIActions interface {
	DeleteThing(ctx context.Context, thingID string) error
	Vote(ctx context.Context, thingID string, dir int) error
	Save(ctx context.Context, thingID string, save bool) error
}

func postActionCommand(use, short string, action func(APIActions, *cobra.Command, string) error) *cobra.Command {
	return &cobra.Command{
		Use:   use + " <postId>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newAPIClient()
			if err != nil {
				return err
			}
			if err := action(client, cmd, args[0]); err != nil {
				return err
			}
			fmt.Printf("%s %s\n", strings.Title(use), args[0])
			return nil
		},
	}
}

func init() {
	postsCmd.AddCommand(postsListCmd, postsGetCmd, postsSubmitCmd, postsDeleteCmd, postsUpvoteCmd, postsDownvoteCmd, postsSaveCmd, postsUnsaveCmd)

	postsListCmd.Flags().StringVar(&postsListSort, "sort", "hot", "Sort: hot, new, top")
	postsListCmd.Flags().StringVar(&postsListTime, "time", "day", "Time filter for top sort")
	postsListCmd.Flags().IntVar(&postsListLimit, "limit", 10, "Number of posts to fetch")

	postsGetCmd.Flags().IntVar(&postComments, "comments", 10, "Number of comments to fetch")

	postsSubmitCmd.Flags().StringVar(&submitTitle, "title", "", "Post title")
	postsSubmitCmd.Flags().StringVar(&submitText, "text", "", "Text body for self posts")
	postsSubmitCmd.Flags().StringVar(&submitURL, "url", "", "URL for link posts")
	postsSubmitCmd.Flags().StringVar(&submitType, "type", "text", "Post type: text or link")
}
