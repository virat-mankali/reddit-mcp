package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/virat-mankali/reddit-mcp/internal/api"
	"github.com/virat-mankali/reddit-mcp/internal/auth"
	"github.com/virat-mankali/reddit-mcp/internal/config"
)

var (
	loginClientID     string
	loginClientSecret string
	loginUserAgent    string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage Reddit authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Reddit via OAuth",
	RunE:  runAuthLogin,
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored credentials",
	RunE:  runAuthLogout,
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	RunE:  runAuthStatus,
}

func init() {
	authLoginCmd.Flags().StringVar(&loginClientID, "client-id", "", "Reddit app client ID")
	authLoginCmd.Flags().StringVar(&loginClientSecret, "client-secret", "", "Reddit app client secret")
	authLoginCmd.Flags().StringVar(&loginUserAgent, "user-agent", "", "Custom Reddit user-agent")

	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
}

func runAuthLogin(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if loginClientID != "" {
		cfg.ClientID = loginClientID
	}
	if loginClientSecret != "" {
		cfg.ClientSecret = loginClientSecret
	}
	if loginUserAgent != "" {
		cfg.UserAgent = loginUserAgent
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return errors.New("missing Reddit app credentials; pass --client-id and --client-secret or set RD_CLIENT_ID and RD_CLIENT_SECRET")
	}

	if cfg.UserAgent == "" {
		cfg.UserAgent = "reddit-mcp/0.1.0"
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	fmt.Println("Opening browser for Reddit authorization...")
	fmt.Printf("Waiting for callback on %s\n\n", auth.RedirectURI)

	manager := auth.NewManager(cfg)
	token, err := manager.Login(cmd.Context())
	if err != nil {
		return err
	}

	client := api.NewClient(cfg)
	me, err := client.Me(cmd.Context())
	if err != nil {
		return fmt.Errorf("login succeeded but fetching profile failed: %w", err)
	}

	token.Username = me.Name
	if err := auth.SaveToken(token); err != nil {
		return err
	}

	expiresIn := time.Until(token.ExpiresAt).Round(time.Minute)
	if expiresIn < 0 {
		expiresIn = 0
	}

	fmt.Printf("Authenticated as u/%s\n", me.Name)
	fmt.Printf("Token expires: %s\n", humanizeDuration(expiresIn))
	fmt.Println("Refresh token: stored permanently")
	return nil
}

func runAuthLogout(cmd *cobra.Command, args []string) error {
	if err := auth.ClearToken(); err != nil {
		return err
	}

	fmt.Println("Stored Reddit tokens cleared.")
	return nil
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	token, err := auth.LoadToken()
	if err != nil {
		return errors.New("not logged in")
	}

	client := api.NewClient(cfg)
	me, err := client.Me(cmd.Context())
	if err != nil {
		return err
	}

	token, err = auth.LoadToken()
	if err != nil {
		return err
	}

	if jsonOutput {
		payload := struct {
			Username      string    `json:"username"`
			LinkKarma     int       `json:"link_karma"`
			CommentKarma  int       `json:"comment_karma"`
			TokenExpires  time.Time `json:"token_expires_at"`
			RefreshStored bool      `json:"refresh_token_stored"`
		}{
			Username:      me.Name,
			LinkKarma:     me.LinkKarma,
			CommentKarma:  me.CommentKarma,
			TokenExpires:  token.ExpiresAt,
			RefreshStored: token.RefreshToken != "",
		}

		out, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	}

	fmt.Printf("Logged in as u/%s\n", me.Name)
	fmt.Printf("Link karma   : %d\n", me.LinkKarma)
	fmt.Printf("Comment karma: %d\n", me.CommentKarma)
	fmt.Printf("Token expires: %s\n", humanizeDuration(time.Until(token.ExpiresAt).Round(time.Minute)))
	return nil
}

func humanizeDuration(d time.Duration) string {
	if d <= 0 {
		return "expired"
	}
	return fmt.Sprintf("in %s", d)
}
