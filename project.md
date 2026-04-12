# reddit-mcp · Project Blueprint

> Reddit CLI + MCP tool — browse, post, search, and manage Reddit from your terminal or let AI agents do it for you.

Command name: **`rd`** — short, fast, memorable.
Language: **Go** — single binary, tiny file size, no runtime.

---

## Table of Contents

1. [Overview](#overview)
2. [Tech Stack](#tech-stack)
3. [Reddit API — What Actually Works](#reddit-api--what-actually-works)
4. [Project Structure](#project-structure)
5. [Phase 1 — Scaffold + Auth](#phase-1--scaffold--auth)
6. [Phase 2 — Core CLI Commands](#phase-2--core-cli-commands)
7. [Phase 3 — MCP Server](#phase-3--mcp-server)
8. [Phase 4 — GoReleaser + Homebrew](#phase-4--goreleaser--homebrew)
9. [Phase 5 — GitHub Actions CI/CD](#phase-5--github-actions-cicd)
10. [All CLI Commands Reference](#all-cli-commands-reference)
11. [All MCP Tools Reference](#all-mcp-tools-reference)
12. [Environment Variables](#environment-variables)
13. [Homebrew Formula](#homebrew-formula)
14. [MCP Config](#mcp-config)

---

## Overview

`reddit-mcp` is a developer-first CLI and MCP server for Reddit built in Go.
Single binary. No Node.js. No runtime. Just `brew install` and go.

It lets you:
- Browse your feed, subreddits, hot/new/top posts from terminal
- Submit posts, comments, replies
- Search Reddit across posts, comments, subreddits
- Read and send DMs
- Vote on posts and comments
- Let AI agents (Claude, Cursor, Kiro) manage Reddit via MCP

Unlike LinkedIn — **everything works** with a normal self-serve OAuth app.

---

## Tech Stack

| Tool | Purpose |
|------|---------|
| **Go 1.22+** | Primary language |
| **cobra** (`github.com/spf13/cobra`) | CLI framework — same as tgcli, discocli |
| **viper** (`github.com/spf13/viper`) | Config management |
| **golang.org/x/oauth2** | OAuth 2.0 flow + token refresh |
| **net/http** | Reddit REST API calls (stdlib, no extra dep) |
| **encoding/json** | JSON parsing (stdlib) |
| **mcp-go** (`github.com/mark3labs/mcp-go`) | MCP server in Go |
| **tablewriter** (`github.com/olekukonko/tablewriter`) | Pretty table output |
| **color** (`github.com/fatih/color`) | Colored terminal output |
| **keyring** (`github.com/zalando/go-keyring`) | Secure OS keychain token storage |
| **GoReleaser** | Cross-platform binary builds + Homebrew tap update |
| **golangci-lint** | Linting |

> **Why `mcp-go` and not the official SDK?**
> The official MCP SDK is TypeScript-first. `mcp-go` by mark3labs is the best Go MCP library — actively maintained, used by many Go MCP servers, implements the full MCP spec over stdio.

---

## Reddit API — What Actually Works

### Getting credentials (2 minutes, no approval)

1. Go to **https://www.reddit.com/prefs/apps**
2. Click **"create another app"**
3. Pick type: **"web app"** (for browser OAuth flow)
4. Set redirect URI: `http://localhost:3141/callback`
5. Hit create → instantly get `client_id` and `client_secret`

No waiting. No approval. No company page. Just works.

### OAuth 2.0 Scopes — all self-serve, all free

| Scope | What it enables |
|-------|----------------|
| `identity` | Read username, karma, account info |
| `read` | Read posts, comments, subreddits, feeds |
| `submit` | Submit posts and comments |
| `vote` | Upvote/downvote posts and comments |
| `history` | Read your post/comment history |
| `privatemessages` | Read and send DMs |
| `mysubreddits` | List subreddits you're subscribed to |
| `save` | Save/unsave posts |
| `subscribe` | Subscribe/unsubscribe to subreddits |
| `flair` | Set your flair in subreddits |

**All of these work with a basic self-serve app. No gating. No partner program.**

### Rate limits

- **100 requests/minute** with OAuth (more than enough)
- Use `User-Agent: reddit-mcp/1.0.0 by u/yourusername` header (Reddit requires this)
- Respect `X-Ratelimit-Remaining` header in responses

### Base API URL

```
https://oauth.reddit.com
```
(Always use this, not `api.reddit.com`, when using OAuth tokens)

### Token details

- Access token: **1 hour** expiry
- Refresh token: **permanent** (until user revokes)
- Auto-refresh before every request when near expiry

---

## Project Structure

```
reddit-mcp/
├── cmd/
│   └── rdcli/
│       ├── main.go           # Entry point
│       ├── root.go           # Cobra root command, --json flag, global flags
│       ├── auth.go           # rd auth login/logout/status
│       ├── me.go             # rd me (profile, karma, trophies)
│       ├── feed.go           # rd feed (home feed)
│       ├── posts.go          # rd posts hot/new/top/submit/delete/vote
│       ├── comments.go       # rd comments list/reply/delete/vote
│       ├── subreddits.go     # rd subreddits list/search/subscribe
│       ├── search.go         # rd search posts/comments/subreddits
│       ├── messages.go       # rd messages list/send/read
│       ├── save.go           # rd save / rd unsave
│       └── serve.go          # rd serve (starts MCP server)
├── internal/
│   ├── auth/
│   │   ├── oauth.go          # OAuth 2.0 PKCE flow, local callback server
│   │   └── token.go          # Token storage (keyring), refresh logic
│   ├── api/
│   │   ├── client.go         # HTTP client with auth, rate limit, User-Agent
│   │   ├── me.go             # /api/v1/me endpoints
│   │   ├── feed.go           # /best, /hot, /new, /top endpoints
│   │   ├── posts.go          # /submit, /del, /vote endpoints
│   │   ├── comments.go       # /comment, /del, vote endpoints
│   │   ├── subreddits.go     # /subreddits/mine, /r/{sub}/about
│   │   ├── search.go         # /search endpoint
│   │   └── messages.go       # /message/inbox, /message/compose
│   ├── config/
│   │   └── config.go         # ~/.rd/config.json via viper
│   ├── mcp/
│   │   ├── server.go         # MCP server setup, tool registration
│   │   └── tools/
│   │       ├── me.go         # MCP tools for profile
│   │       ├── feed.go       # MCP tools for feed
│   │       ├── posts.go      # MCP tools for posts
│   │       ├── comments.go   # MCP tools for comments
│   │       ├── search.go     # MCP tools for search
│   │       ├── subreddits.go # MCP tools for subreddits
│   │       └── messages.go   # MCP tools for messages
│   └── output/
│       └── output.go         # Table, JSON, pretty formatters
├── .github/
│   └── workflows/
│       ├── ci.yml            # PR → vet + lint + test
│       └── release.yml       # Tag → GoReleaser → Homebrew update
├── .goreleaser.yaml          # Cross-platform build + Homebrew tap
├── .gitignore
├── .golangci.yml             # Linter config
├── LICENSE
├── README.md
├── project.md                # This file
├── go.mod
└── go.sum
```

---

## Phase 1 — Scaffold + Auth

**Goal:** Repo initialized, `rd auth login` works end-to-end, tokens saved securely.

### 1.1 Init the repo

```bash
mkdir reddit-mcp && cd reddit-mcp
go mod init github.com/virat-mankali/reddit-mcp

# Core deps
go get github.com/spf13/cobra
go get github.com/spf13/viper
go get golang.org/x/oauth2
go get github.com/mark3labs/mcp-go
go get github.com/olekukonko/tablewriter
go get github.com/fatih/color
go get github.com/zalando/go-keyring
```

### 1.2 `go.mod` (key deps)

```
module github.com/virat-mankali/reddit-mcp

go 1.22

require (
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.0
    golang.org/x/oauth2 v0.18.0
    github.com/mark3labs/mcp-go v0.8.0
    github.com/olekukonko/tablewriter v0.0.5
    github.com/fatih/color v1.16.0
    github.com/zalando/go-keyring v0.2.3
)
```

### 1.3 `cmd/rdcli/main.go`

```go
package main

import (
    "github.com/virat-mankali/reddit-mcp/cmd/rdcli/commands"
)

func main() {
    commands.Execute()
}
```

### 1.4 `cmd/rdcli/root.go`

```go
package commands

import (
    "github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
    Use:     "rd",
    Short:   "Reddit CLI + MCP — browse, post, search Reddit from your terminal",
    Version: "1.0.0",
}

func Execute() {
    rootCmd.Execute()
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
```

### 1.5 OAuth 2.0 Flow (`internal/auth/oauth.go`)

Reddit uses standard OAuth 2.0 with PKCE. Flow:

1. Generate `state` param (random 32-char string for CSRF protection)
2. Build authorize URL:
   ```
   https://www.reddit.com/api/v1/authorize
     ?client_id=CLIENT_ID
     &response_type=code
     &state=STATE
     &redirect_uri=http://localhost:3141/callback
     &duration=permanent
     &scope=identity read submit vote history privatemessages mysubreddits save subscribe
   ```
3. Open browser with `open` (macOS) / `xdg-open` (Linux)
4. Spin up `net/http` server on `localhost:3141` to catch the callback
5. Exchange `code` for tokens at:
   ```
   POST https://www.reddit.com/api/v1/access_token
   Authorization: Basic base64(client_id:client_secret)
   Body: grant_type=authorization_code&code=CODE&redirect_uri=...
   ```
6. Store `access_token` + `refresh_token` in OS keyring via `go-keyring`
7. Store `client_id` + `client_secret` in `~/.rd/config.json` via viper

```go
// Key constants in internal/auth/oauth.go
const (
    AuthURL     = "https://www.reddit.com/api/v1/authorize"
    TokenURL    = "https://www.reddit.com/api/v1/access_token"
    RedirectURI = "http://localhost:3141/callback"
    Scopes      = "identity read submit vote history privatemessages mysubreddits save subscribe"
    KeyringService = "reddit-mcp"
)
```

### 1.6 Token storage (`internal/auth/token.go`)

```go
type TokenStore struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
    Username     string    `json:"username"`
}

// Store in OS keychain:
// keyring.Set("reddit-mcp", "access_token", token)
// keyring.Set("reddit-mcp", "refresh_token", token)

// Refresh logic: if time.Now().After(expiresAt.Add(-5 * time.Minute))
// POST https://www.reddit.com/api/v1/access_token
// grant_type=refresh_token&refresh_token=REFRESH_TOKEN
```

### 1.7 `cmd/rdcli/auth.go`

```go
var authCmd = &cobra.Command{
    Use:   "auth",
    Short: "Manage Reddit authentication",
}

// rd auth login
var authLoginCmd = &cobra.Command{
    Use:   "login",
    Short: "Authenticate with Reddit via OAuth",
    RunE:  runAuthLogin,
}

// rd auth logout
var authLogoutCmd = &cobra.Command{
    Use:   "logout",
    Short: "Clear stored credentials",
    RunE:  runAuthLogout,
}

// rd auth status
var authStatusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show current authentication status",
    RunE:  runAuthStatus,
}
```

Output of `rd auth login`:
```
Opening browser for Reddit authorization...
Waiting for callback on http://localhost:3141/callback

✓ Authenticated as u/virat-mankali
  Token expires: in 59 minutes (auto-refreshes)
  Refresh token: stored permanently
```

Output of `rd auth status`:
```
✓ Logged in as u/virat-mankali
  Link karma  : 1,204
  Comment karma: 847
  Token expires: in 42 minutes
```

### 1.8 `internal/api/client.go`

```go
type Client struct {
    httpClient *http.Client
    baseURL    string
    userAgent  string
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{Timeout: 30 * time.Second},
        baseURL:    "https://oauth.reddit.com",
        userAgent:  "reddit-mcp/1.0.0 by u/virat-mankali",
    }
}

// Before every request:
// 1. Check token expiry, refresh if needed
// 2. Set Authorization: bearer ACCESS_TOKEN
// 3. Set User-Agent header (Reddit bans requests without it)
// 4. Check X-Ratelimit-Remaining in response, sleep if 0
```

### 1.9 Deliverables for Phase 1

- [ ] `go mod init` + all deps fetched
- [ ] `rd --help` and `rd <command> --help` work
- [ ] `rd auth login` opens browser, captures code, exchanges for tokens
- [ ] Tokens stored in OS keyring (not plaintext files)
- [ ] `rd auth logout` clears keyring
- [ ] `rd auth status` shows username + karma + token expiry
- [ ] Auto-refresh works silently before every API call
- [ ] `User-Agent` header set correctly on all requests

---

## Phase 2 — Core CLI Commands

**Goal:** All main CLI commands working against the Reddit REST API.

### 2.1 Profile — `rd me`

**Endpoint:** `GET /api/v1/me`

```
rd me                          # Your profile: username, karma, cake day, etc.
rd me --json                   # Raw JSON
```

Output:
```
Username    : u/virat-mankali
Link karma  : 1,204
Comment karma: 847
Cake day    : March 15, 2022
Gold        : No
Verified    : Yes
```

### 2.2 Feed — `rd feed`

**Endpoints:**
- `GET /best` — best posts on your home feed
- `GET /hot` — hot posts home feed
- `GET /new` — new posts home feed
- `GET /top?t=day|week|month|year|all`

```
rd feed                        # Home feed (best, default)
rd feed --sort hot
rd feed --sort new
rd feed --sort top --time week
rd feed --limit 25             # Default 10
rd feed --json
```

Output:
```
#   Score  Sub              Title                                        Author
1   4,821  r/programming    "Go 1.23 released with new iterator..."      u/gopherine
2   2,103  r/golang         "I built a Reddit CLI in Go"                 u/virat-mankali
3    987   r/opensource     "Why open source matters in 2026"            u/devhacker
```

### 2.3 Posts — `rd posts`

**Endpoints:**
- `GET /r/{subreddit}/hot|new|top` — browse subreddit
- `POST /api/submit` — submit post
- `POST /api/del` — delete post
- `POST /api/vote` — vote

```
rd posts list r/golang                         # Hot posts in r/golang
rd posts list r/golang --sort new
rd posts list r/golang --sort top --time month
rd posts list r/golang --limit 20

rd posts get <postId>                          # Get full post + top comments
rd posts get <postId> --comments 25

rd posts submit r/golang \
  --title "I built a Reddit CLI in Go" \
  --text "Check it out at github.com/..." \
  --type text                                  # text post

rd posts submit r/golang \
  --title "Cool project" \
  --url "https://github.com/virat-mankali/reddit-mcp" \
  --type link                                  # link post

rd posts delete <postId>
rd posts upvote <postId>
rd posts downvote <postId>
rd posts save <postId>
rd posts unsave <postId>
```

### 2.4 Comments — `rd comments`

**Endpoints:**
- `GET /comments/{postId}` — list comments on a post
- `POST /api/comment` — reply to post or comment
- `POST /api/del` — delete comment
- `POST /api/vote` — vote on comment

```
rd comments list <postId>                      # List comments on a post
rd comments list <postId> --limit 50
rd comments reply <postId> "Great post!"       # Comment on a post
rd comments reply --parent <commentId> "Nice!" # Reply to a comment
rd comments delete <commentId>
rd comments upvote <commentId>
rd comments downvote <commentId>
```

### 2.5 Subreddits — `rd subreddits`

**Endpoints:**
- `GET /subreddits/mine/subscriber` — subreddits you're subscribed to
- `GET /r/{sub}/about` — subreddit info
- `POST /api/subscribe` — subscribe/unsubscribe

```
rd subreddits list                             # Your subscribed subreddits
rd subreddits list --limit 50
rd subreddits info r/golang                    # Info about a subreddit
rd subreddits search "typescript"              # Search subreddits by name/topic
rd subreddits subscribe r/golang
rd subreddits unsubscribe r/golang
```

Output of `rd subreddits list`:
```
Subreddit          Subscribers   Description
r/golang           150,234       The Go programming language
r/programming      6,200,100     Computer programming
r/opensource       89,400        Open source software
```

### 2.6 Search — `rd search`

**Endpoint:** `GET /search?q=QUERY&type=link|comment|sr`

```
rd search "reddit cli golang"                  # Search posts (default)
rd search "mcp server" --type posts
rd search "virat-mankali" --type users
rd search "golang cli" --type subreddits
rd search "open source tools" --sub r/programming    # Search within subreddit
rd search "hello" --sort relevance|new|top
rd search "hello" --time hour|day|week|month|year|all
rd search "hello" --limit 20
rd search "hello" --json
```

### 2.7 Messages (DMs) — `rd messages`

**Endpoints:**
- `GET /message/inbox` — inbox
- `GET /message/unread` — unread messages
- `GET /message/sent` — sent messages
- `POST /api/compose` — send a message
- `POST /api/read_message` — mark as read

```
rd messages list                               # Inbox
rd messages list --unread                      # Unread only
rd messages list --sent                        # Sent messages
rd messages send --to u/username --subject "Hey" --body "What's up!"
rd messages read <messageId>                   # Mark as read
```

### 2.8 Output formatting

Every command supports:
- `--json` — raw Reddit API JSON, great for piping to `jq`
- `--limit N` — control how many results to fetch
- Default pretty table output via `tablewriter`

```bash
rd posts list r/golang --json | jq '.[0].title'
rd me --json | jq '.total_karma'
rd feed --json | jq '[.[] | {title, score, subreddit}]'
```

### 2.9 Deliverables for Phase 2

- [ ] `rd me` works
- [ ] `rd feed` with all sort options works
- [ ] `rd posts list/get/submit/delete/upvote/downvote/save` work
- [ ] `rd comments list/reply/delete/upvote/downvote` work
- [ ] `rd subreddits list/info/search/subscribe/unsubscribe` work
- [ ] `rd search` with all type/sort/time options works
- [ ] `rd messages list/send/read` work
- [ ] `--json` flag works on every command
- [ ] Rate limit headers respected (`X-Ratelimit-Remaining`)
- [ ] All commands have `--help`

---

## Phase 3 — MCP Server

**Goal:** `rd serve` starts a stdio MCP server exposing all Reddit actions as tools for AI agents.

### 3.1 MCP library

Using **`github.com/mark3labs/mcp-go`** — the standard Go MCP library.

```go
// internal/mcp/server.go
import "github.com/mark3labs/mcp-go/mcp"
import "github.com/mark3labs/mcp-go/server"

func NewServer() *server.MCPServer {
    s := server.NewMCPServer(
        "reddit-mcp",
        "1.0.0",
        server.WithToolCapabilities(true),
    )

    registerMeTools(s)
    registerFeedTools(s)
    registerPostTools(s)
    registerCommentTools(s)
    registerSearchTools(s)
    registerSubredditTools(s)
    registerMessageTools(s)

    return s
}
```

### 3.2 MCP tool pattern (example: `submit_post`)

```go
// internal/mcp/tools/posts.go
func registerPostTools(s *server.MCPServer) {
    s.AddTool(mcp.NewTool(
        "submit_post",
        mcp.WithDescription("Submit a new post to a subreddit"),
        mcp.WithString("subreddit",
            mcp.Required(),
            mcp.Description("Subreddit name without r/ prefix, e.g. golang"),
        ),
        mcp.WithString("title",
            mcp.Required(),
            mcp.Description("Post title"),
        ),
        mcp.WithString("type",
            mcp.Description("Post type: text or link"),
            mcp.Enum("text", "link"),
            mcp.DefaultValue("text"),
        ),
        mcp.WithString("text",
            mcp.Description("Post body (for text posts)"),
        ),
        mcp.WithString("url",
            mcp.Description("URL to share (for link posts)"),
        ),
    ), handleSubmitPost)
}

func handleSubmitPost(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Extract params, call api.SubmitPost(), return result
}
```

### 3.3 `cmd/rdcli/serve.go`

```go
var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the Reddit MCP server (stdio)",
    Long:  "Starts an MCP server over stdio for use with Claude, Cursor, Kiro, etc.",
    RunE:  runServe,
}

func runServe(cmd *cobra.Command, args []string) error {
    s := mcp.NewServer()
    if err := server.ServeStdio(s); err != nil {
        return fmt.Errorf("MCP server error: %w", err)
    }
    return nil
}
```

### 3.4 Auth in MCP mode

MCP mode reads tokens from keyring exactly like CLI mode. If no tokens found, tool responses return a clean error:

```json
{
  "error": "Not authenticated. Run: rd auth login"
}
```

Env var override (headless/CI):
```
REDDIT_CLIENT_ID=xxx
REDDIT_CLIENT_SECRET=xxx
REDDIT_ACCESS_TOKEN=xxx
REDDIT_REFRESH_TOKEN=xxx
```

If env vars are set, they take priority over keyring. This is the pattern for MCP config with explicit tokens.

### 3.5 Deliverables for Phase 3

- [ ] `rd serve` starts over stdio with no errors
- [ ] All MCP tools listed in [All MCP Tools Reference](#all-mcp-tools-reference) registered
- [ ] Input schema validated for each tool (required fields, enums)
- [ ] Auth errors surface cleanly in tool result (don't crash server)
- [ ] Env var token override works
- [ ] Tested with Claude Desktop — tools appear and execute correctly
- [ ] Tested with Cursor/Kiro MCP config

---

## Phase 4 — GoReleaser + Homebrew

**Goal:** Single binary installable via `brew install virat-mankali/tap/reddit-mcp`.

### 4.1 `.goreleaser.yaml`

```yaml
version: 2

project_name: reddit-mcp

before:
  hooks:
    - go mod tidy
    - go generate ./...

builds:
  - id: rd
    main: ./cmd/rdcli/
    binary: rd
    env:
      - CGO_ENABLED=0
    goos:
      - darwin
      - linux
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}

archives:
  - id: rd
    builds: [rd]
    name_template: >-
      rd_
      {{- .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "arm64" }}arm64
      {{- else }}{{ .Arch }}{{ end }}
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: "checksums.txt"
  algorithm: sha256

brews:
  - name: reddit-mcp
    ids: [rd]
    repository:
      owner: virat-mankali
      name: homebrew-tap
      token: "{{ .Env.TAP_GITHUB_TOKEN }}"
    directory: Formula
    homepage: "https://github.com/virat-mankali/reddit-mcp"
    description: "Reddit CLI + MCP tool — browse, post, search Reddit from your terminal or let AI agents do it"
    license: "MIT"
    test: |
      system "#{bin}/rd --version"
    install: |
      bin.install "rd"

release:
  github:
    owner: virat-mankali
    name: reddit-mcp
  draft: false
  prerelease: auto
  name_template: "{{.ProjectName}} v{{.Version}}"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
      - Merge pull request
```

### 4.2 What GoReleaser does on tag push

1. Builds `rd` binary for: `darwin/amd64`, `darwin/arm64`, `linux/amd64`, `linux/arm64`, `windows/amd64`
2. Strips debug symbols (`-s -w`) — smallest possible binary
3. Tars each into `rd_darwin_arm64.tar.gz`, etc.
4. Computes SHA256 checksums
5. Creates GitHub Release with all archives + `checksums.txt`
6. **Auto-updates `Formula/reddit-mcp.rb`** in your existing `homebrew-tap` repo with new version + SHAs

Zero manual work after tagging.

### 4.3 Binary size estimate

Go with `-s -w` + `CGO_ENABLED=0`:
- macOS arm64: ~8-12MB
- Linux x64: ~8-12MB

vs TypeScript binary via pkg: ~50-80MB. That's the Go advantage.

### 4.4 Homebrew formula (auto-generated by GoReleaser, lives in tap repo)

```ruby
class RedditMcp < Formula
  desc "Reddit CLI + MCP tool — browse, post, search Reddit from your terminal"
  homepage "https://github.com/virat-mankali/reddit-mcp"
  version "1.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/virat-mankali/reddit-mcp/releases/download/v#{version}/rd_darwin_arm64.tar.gz"
      sha256 "AUTO_FILLED_BY_GORELEASER"
    else
      url "https://github.com/virat-mankali/reddit-mcp/releases/download/v#{version}/rd_darwin_x86_64.tar.gz"
      sha256 "AUTO_FILLED_BY_GORELEASER"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/virat-mankali/reddit-mcp/releases/download/v#{version}/rd_linux_arm64.tar.gz"
      sha256 "AUTO_FILLED_BY_GORELEASER"
    else
      url "https://github.com/virat-mankali/reddit-mcp/releases/download/v#{version}/rd_linux_x86_64.tar.gz"
      sha256 "AUTO_FILLED_BY_GORELEASER"
    end
  end

  def install
    bin.install "rd"
  end

  test do
    assert_match "reddit-mcp", shell_output("#{bin}/rd --version")
  end
end
```

### 4.5 Install command (end user)

```bash
brew install virat-mankali/tap/reddit-mcp
```

### 4.6 Deliverables for Phase 4

- [ ] `.goreleaser.yaml` committed and tested with `goreleaser build --snapshot --clean`
- [ ] `TAP_GITHUB_TOKEN` secret added to repo settings
- [ ] `Formula/reddit-mcp.rb` auto-created in homebrew-tap on first release
- [ ] Tested: `brew install virat-mankali/tap/reddit-mcp && rd --help` works
- [ ] Binary is truly self-contained — no Go runtime needed on user machine

---

## Phase 5 — GitHub Actions CI/CD

**Goal:** PR checks + tag → full automated release.

### 5.1 `ci.yml` — runs on every PR

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
          cache: true

      - name: Vet
        run: go vet ./...

      - name: Lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest

      - name: Test
        run: go test ./... -v -race

      - name: Build (smoke test)
        run: go build -o rd ./cmd/rdcli/
```

### 5.2 `release.yml` — runs on `v*` tag push

```yaml
name: Release

on:
  push:
    tags:
      - "v*"

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0       # GoReleaser needs full git history for changelog

      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
          cache: true

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          TAP_GITHUB_TOKEN: ${{ secrets.TAP_GITHUB_TOKEN }}
```

### 5.3 Required GitHub Secrets

| Secret | Where to set | What it is |
|--------|-------------|-----------|
| `TAP_GITHUB_TOKEN` | reddit-mcp repo → Settings → Secrets | GitHub PAT with `repo` scope, for writing to homebrew-tap |
| `GITHUB_TOKEN` | auto-provided by GitHub Actions | For creating GitHub Release |

### 5.4 Release workflow (one command)

```bash
git tag v1.0.0 && git push origin v1.0.0

# GitHub Actions will automatically:
# 1. go vet + lint + test
# 2. GoReleaser builds 5 platform binaries (darwin arm64/x64, linux arm64/x64, windows x64)
# 3. Strips debug symbols → tiny binaries
# 4. Creates GitHub Release with all .tar.gz + checksums.txt
# 5. Updates Formula/reddit-mcp.rb in homebrew-tap with new SHAs
# 6. Users run: brew upgrade reddit-mcp  ✓
```

### 5.5 Deliverables for Phase 5

- [ ] `ci.yml` passes on main branch
- [ ] `release.yml` triggers on `v*` tags
- [ ] All platform binaries created in GitHub Release
- [ ] `checksums.txt` published alongside binaries
- [ ] `Formula/reddit-mcp.rb` auto-updated in homebrew-tap
- [ ] Full cycle tested: `git tag v0.1.0 && git push origin v0.1.0` → brew install works

---

## All CLI Commands Reference

```bash
# Auth
rd auth login                                  # OAuth login via browser
rd auth logout                                 # Clear stored credentials
rd auth status                                 # Show current user + token info

# Profile
rd me                                          # Your profile, karma, cake day
rd me --json

# Feed
rd feed                                        # Home feed (best)
rd feed --sort hot|new|top
rd feed --sort top --time day|week|month|year|all
rd feed --limit 25
rd feed --json

# Posts
rd posts list r/golang                         # Hot posts in subreddit
rd posts list r/golang --sort new|top|hot|rising
rd posts list r/golang --sort top --time week
rd posts list r/golang --limit 20
rd posts get <postId>                          # Full post + comments
rd posts get <postId> --comments 25
rd posts submit r/golang --title "Title" --text "Body" --type text
rd posts submit r/golang --title "Title" --url "https://..." --type link
rd posts delete <postId>
rd posts upvote <postId>
rd posts downvote <postId>
rd posts save <postId>
rd posts unsave <postId>

# Comments
rd comments list <postId>
rd comments list <postId> --limit 50
rd comments reply <postId> "Comment text"
rd comments reply --parent <commentId> "Reply text"
rd comments delete <commentId>
rd comments upvote <commentId>
rd comments downvote <commentId>

# Subreddits
rd subreddits list
rd subreddits list --limit 100
rd subreddits info r/golang
rd subreddits search "typescript"
rd subreddits subscribe r/golang
rd subreddits unsubscribe r/golang

# Search
rd search "golang cli"
rd search "mcp server" --type posts|comments|subreddits|users
rd search "hello" --sub r/programming
rd search "tools" --sort relevance|new|top
rd search "tools" --time hour|day|week|month|year|all
rd search "tools" --limit 20
rd search "tools" --json

# Messages
rd messages list
rd messages list --unread
rd messages list --sent
rd messages send --to u/username --subject "Subject" --body "Message body"
rd messages read <messageId>

# MCP Server
rd serve

# Meta
rd --version
rd --help
rd <command> --help
```

---

## All MCP Tools Reference

| Tool | Description |
|------|-------------|
| `get_my_profile` | Get your Reddit profile, karma, account info |
| `get_home_feed` | Get your personalized home feed with sort/time options |
| `list_subreddit_posts` | List posts from any subreddit |
| `get_post` | Get full post details with top comments |
| `submit_post` | Submit a text or link post to a subreddit |
| `delete_post` | Delete one of your posts |
| `vote_post` | Upvote or downvote a post |
| `save_post` | Save a post to your saved list |
| `list_comments` | List comments on a post |
| `reply_to_post` | Comment on a post |
| `reply_to_comment` | Reply to an existing comment |
| `delete_comment` | Delete your comment |
| `vote_comment` | Upvote or downvote a comment |
| `search_reddit` | Search posts, comments, subreddits, or users |
| `search_subreddit` | Search within a specific subreddit |
| `list_my_subreddits` | List subreddits you're subscribed to |
| `get_subreddit_info` | Get info and stats about a subreddit |
| `subscribe_subreddit` | Subscribe to a subreddit |
| `unsubscribe_subreddit` | Unsubscribe from a subreddit |
| `list_messages` | List inbox, unread, or sent messages |
| `send_message` | Send a DM to a Reddit user |
| `read_message` | Mark a message as read |

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `REDDIT_CLIENT_ID` | OAuth app Client ID |
| `REDDIT_CLIENT_SECRET` | OAuth app Client Secret |
| `REDDIT_ACCESS_TOKEN` | Access token override (headless/CI/MCP) |
| `REDDIT_REFRESH_TOKEN` | Refresh token override (headless/CI/MCP) |

Default token storage: **OS keyring** (macOS Keychain, Linux Secret Service, Windows Credential Manager).
Config file: `~/.rd/config.json` (client_id, client_secret, username).

---

## MCP Config

### Claude Desktop / Cursor / Kiro

```json
{
  "mcpServers": {
    "reddit": {
      "command": "rd",
      "args": ["serve"]
    }
  }
}
```

### Headless with env var tokens

```json
{
  "mcpServers": {
    "reddit": {
      "command": "rd",
      "args": ["serve"],
      "env": {
        "REDDIT_CLIENT_ID": "your_client_id",
        "REDDIT_CLIENT_SECRET": "your_client_secret",
        "REDDIT_ACCESS_TOKEN": "your_access_token",
        "REDDIT_REFRESH_TOKEN": "your_refresh_token"
      }
    }
  }
}
```

---

## Milestones Summary

| Phase | Deliverable | Effort |
|-------|-------------|--------|
| Phase 1 | Scaffold + Auth (`rd auth login` works) | ~1 day |
| Phase 2 | All CLI commands working | ~2-3 days |
| Phase 3 | MCP server + all 22 tools | ~1 day |
| Phase 4 | GoReleaser + Homebrew formula | ~2 hours |
| Phase 5 | GitHub Actions CI/CD | ~2 hours |
| **Total** | **Ship-ready `reddit-mcp`** | **~5-6 days** |

---

## Order of Operations (what to build first)

1. `rd auth login` — nothing else works without this
2. `rd me` — validates auth + API client are working
3. `rd feed` + `rd posts list` — core browse experience
4. `rd posts submit` + `rd comments reply` — write access
5. `rd search` — discovery
6. `rd messages` — DMs
7. `rd subreddits` — subscription management
8. `rd serve` (MCP) — AI agent access
9. GoReleaser config + smoke test
10. GitHub Actions → tag v0.1.0 → ship

---

## Getting Reddit Credentials.

1. Go to **https://www.reddit.com/prefs/apps**
2. Click **"create another app"**
3. Select type: **"web app"**
4. Name: `reddit-mcp` (or anything)
5. Redirect URI: `http://localhost:3141/callback`
6. Click create
7. Copy `client_id` (under the app name) and `client_secret`
8. Run `rd auth login --client-id X --client-secret Y`

---

*Built by Virat Mankali · MIT License*
