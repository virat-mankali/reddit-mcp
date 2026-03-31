package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/virat-mankali/reddit-mcp/internal/output"
)

var (
	messageLimit int
	unreadOnly   bool
	sentOnly     bool
	messageTo    string
	messageSubj  string
	messageBody  string
)

var messagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "Read and send Reddit messages",
}

var messagesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List inbox, unread, or sent messages",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		folder := "inbox"
		if unreadOnly {
			folder = "unread"
		}
		if sentOnly {
			folder = "sent"
		}
		items, err := client.ListMessages(cmd.Context(), folder, messageLimit)
		if err != nil {
			return err
		}
		if jsonOutput {
			return printJSON(items)
		}
		rows := make([][]string, 0, len(items))
		for _, item := range items {
			counterparty := item.Author
			if folder == "sent" {
				counterparty = item.Dest
			}
			rows = append(rows, []string{
				counterparty,
				item.Subject,
				formatTime(item.CreatedAt()),
				trimBody(item.Body, 90),
			})
		}
		output.PrintTable([]string{"User", "Subject", "Sent", "Body"}, rows)
		return nil
	},
}

var messagesSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a Reddit message",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		if messageTo == "" || messageSubj == "" || messageBody == "" {
			return fmt.Errorf("--to, --subject, and --body are required")
		}
		if err := client.SendMessage(cmd.Context(), messageTo, messageSubj, messageBody); err != nil {
			return err
		}
		fmt.Printf("Sent message to %s\n", messageTo)
		return nil
	},
}

var messagesReadCmd = &cobra.Command{
	Use:   "read <messageId>",
	Short: "Mark a message as read",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		if err := client.MarkMessageRead(cmd.Context(), args[0]); err != nil {
			return err
		}
		fmt.Printf("Marked %s as read\n", args[0])
		return nil
	},
}

func init() {
	messagesCmd.AddCommand(messagesListCmd, messagesSendCmd, messagesReadCmd)
	messagesListCmd.Flags().IntVar(&messageLimit, "limit", 10, "Number of messages to fetch")
	messagesListCmd.Flags().BoolVar(&unreadOnly, "unread", false, "Show unread messages")
	messagesListCmd.Flags().BoolVar(&sentOnly, "sent", false, "Show sent messages")

	messagesSendCmd.Flags().StringVar(&messageTo, "to", "", "Recipient username")
	messagesSendCmd.Flags().StringVar(&messageSubj, "subject", "", "Message subject")
	messagesSendCmd.Flags().StringVar(&messageBody, "body", "", "Message body")
}
