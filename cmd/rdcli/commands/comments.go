package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/virat-mankali/reddit-mcp/internal/output"
)

var (
	commentsLimit int
	replyParent   string
)

var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Browse and manage comments",
}

var commentsListCmd = &cobra.Command{
	Use:   "list <postId>",
	Short: "List comments on a post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		comments, err := client.ListComments(cmd.Context(), args[0], commentsLimit)
		if err != nil {
			return err
		}
		if jsonOutput {
			return printJSON(comments)
		}
		rows := make([][]string, 0, len(comments))
		for _, comment := range comments {
			rows = append(rows, []string{
				strings.Repeat("  ", comment.Depth) + "u/" + comment.Author,
				fmt.Sprintf("%d", comment.Score),
				trimBody(comment.Body, 110),
			})
		}
		output.PrintTable([]string{"Author", "Score", "Body"}, rows)
		return nil
	},
}

var commentsReplyCmd = &cobra.Command{
	Use:   "reply <thingId> <text>",
	Short: "Reply to a post or comment",
	Args: func(cmd *cobra.Command, args []string) error {
		if replyParent != "" && len(args) < 1 {
			return fmt.Errorf("reply text is required when using --parent")
		}
		if replyParent == "" && len(args) < 2 {
			return fmt.Errorf("post ID and reply text are required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		parent := replyParent
		bodyIndex := 1
		if parent == "" {
			parent = fullname(args[0], "t3")
		} else {
			parent = fullname(parent, "t1")
			bodyIndex = 0
		}
		text := strings.Join(args[bodyIndex:], " ")
		if err := client.Reply(cmd.Context(), parent, text); err != nil {
			return err
		}
		fmt.Println("Reply posted.")
		return nil
	},
}

type commentActions interface {
	DeleteThing(ctx context.Context, thingID string) error
	Vote(ctx context.Context, thingID string, dir int) error
}

var commentsDeleteCmd = commentActionCommand("delete", "Delete a comment", func(client commentActions, cmd *cobra.Command, id string) error {
	return client.DeleteThing(cmd.Context(), fullname(id, "t1"))
})
var commentsUpvoteCmd = commentActionCommand("upvote", "Upvote a comment", func(client commentActions, cmd *cobra.Command, id string) error {
	return client.Vote(cmd.Context(), fullname(id, "t1"), 1)
})
var commentsDownvoteCmd = commentActionCommand("downvote", "Downvote a comment", func(client commentActions, cmd *cobra.Command, id string) error {
	return client.Vote(cmd.Context(), fullname(id, "t1"), -1)
})

func commentActionCommand(use, short string, action func(commentActions, *cobra.Command, string) error) *cobra.Command {
	return &cobra.Command{
		Use:   use + " <commentId>",
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
	commentsCmd.AddCommand(commentsListCmd, commentsReplyCmd, commentsDeleteCmd, commentsUpvoteCmd, commentsDownvoteCmd)
	commentsListCmd.Flags().IntVar(&commentsLimit, "limit", 20, "Number of comments to fetch")
	commentsReplyCmd.Flags().StringVar(&replyParent, "parent", "", "Parent comment ID to reply to instead of a post")
}
