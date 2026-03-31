package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show your Reddit profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		me, err := client.Me(cmd.Context())
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(me)
		}

		fmt.Printf("Username      : u/%s\n", me.Name)
		fmt.Printf("Link karma    : %d\n", me.LinkKarma)
		fmt.Printf("Comment karma : %d\n", me.CommentKarma)
		fmt.Printf("Total karma   : %d\n", me.TotalKarma)
		fmt.Printf("Cake day      : %s\n", me.CreatedAt().Format("January 2, 2006"))
		fmt.Printf("Gold          : %t\n", me.IsGold)
		fmt.Printf("Verified email: %t\n", me.HasVerifiedEM)
		return nil
	},
}
