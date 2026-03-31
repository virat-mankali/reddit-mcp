# reddit-mcp

Reddit CLI and MCP server written in Go.

The shipped command is `rdcli`.

## Current scope

- Phase 1: scaffold and auth flow
- Phase 2: core CLI commands
- Phase 3: stdio MCP server with Reddit tools
- Phase 4: GoReleaser/Homebrew release config started

## Local usage

```bash
go build -o bin/rdcli ./cmd/rdcli
./bin/rdcli --help
./bin/rdcli auth login --client-id YOUR_ID --client-secret YOUR_SECRET
./bin/rdcli serve
```

## Environment variables

```bash
REDDIT_CLIENT_ID=...
REDDIT_CLIENT_SECRET=...
REDDIT_ACCESS_TOKEN=...
REDDIT_REFRESH_TOKEN=...
REDDIT_USER_AGENT=reddit-mcp/1.0 by u/yourusername
```

When access or refresh tokens are provided via environment variables, the CLI and MCP server will prefer them over the keyring.

## Release

The repo includes a `.goreleaser.yaml` configured to build and package the `rdcli` binary and publish a Homebrew formula named `reddit-mcp`.
