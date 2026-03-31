# reddit-mcp

Reddit CLI and MCP server built in Go.

`reddit-mcp` gives you one binary, `rdcli`, that can:

- browse Reddit from the terminal
- post, comment, vote, save, and message
- expose Reddit actions as MCP tools for AI agents

No Node.js runtime. No Python setup. Just a Go binary.

## Install

### Homebrew

```bash
brew install virat-mankali/tap/reddit-mcp
```

### Build locally

```bash
git clone https://github.com/virat-mankali/reddit-mcp.git
cd reddit-mcp
go build -o bin/rdcli ./cmd/rdcli
./bin/rdcli --help
```

## Reddit app setup

Create a Reddit OAuth app first:

1. Go to `https://www.reddit.com/prefs/apps`
2. Click `create another app`
3. Choose `web app`
4. Set redirect URI to `http://localhost:3141/callback`
5. Copy the `client_id` and `client_secret`

## Authentication

Login once and tokens will be stored securely in your system keyring.

```bash
rdcli auth login --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET
rdcli auth status
rdcli auth logout
```

Config is stored in:

```bash
~/.rdcli/config.json
```

## CLI usage

### Profile

```bash
rdcli me
rdcli me --json
```

### Feed

```bash
rdcli feed
rdcli feed --sort hot
rdcli feed --sort top --time week
rdcli feed --limit 25
```

### Posts

```bash
rdcli posts list r/golang
rdcli posts get abc123

rdcli posts submit r/golang --title "Hello" --text "Post body" --type text
rdcli posts submit r/golang --title "Project" --url "https://github.com/..." --type link

rdcli posts upvote abc123
rdcli posts downvote abc123
rdcli posts save abc123
rdcli posts unsave abc123
rdcli posts delete abc123
```

### Comments

```bash
rdcli comments list abc123
rdcli comments reply abc123 "Nice post"
rdcli comments reply --parent def456 "Nice reply"
rdcli comments upvote def456
rdcli comments downvote def456
rdcli comments delete def456
```

### Subreddits

```bash
rdcli subreddits list
rdcli subreddits info r/golang
rdcli subreddits search golang
rdcli subreddits subscribe r/golang
rdcli subreddits unsubscribe r/golang
```

### Search

```bash
rdcli search "reddit cli"
rdcli search "mcp server" --type posts
rdcli search "open source" --type comments
rdcli search "golang" --type subreddits
rdcli search "virat" --type users
```

### Messages

```bash
rdcli messages list
rdcli messages list --unread
rdcli messages list --sent
rdcli messages send --to u/username --subject "Hello" --body "Message text"
rdcli messages read t4_abc123
```

## JSON output

Every main command supports `--json`.

```bash
rdcli me --json
rdcli feed --json
rdcli posts list r/golang --json
```

## MCP server

Start the MCP server over stdio:

```bash
rdcli serve
```

This exposes Reddit actions as MCP tools so clients like Claude Desktop, Cursor, or other MCP-compatible apps can use them.

## Environment variables

You can run the CLI or MCP server with environment variables instead of saved credentials:

```bash
REDDIT_CLIENT_ID=...
REDDIT_CLIENT_SECRET=...
REDDIT_ACCESS_TOKEN=...
REDDIT_REFRESH_TOKEN=...
REDDIT_USER_AGENT="reddit-mcp/1.0 by u/yourusername"
```

If token env vars are present, they take priority over keyring-stored tokens.

## Release

This repo includes:

- GoReleaser config
- GitHub Actions CI
- GitHub Actions release workflow
- Homebrew formula publishing to `virat-mankali/homebrew-tap`

Release flow:

1. push code to `main`
2. create and push a version tag like `v1.0`
3. GitHub Actions runs GoReleaser
4. binaries are published to GitHub Releases
5. Homebrew formula is updated in your tap repo

## License

MIT
