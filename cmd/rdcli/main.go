package main

import "github.com/virat-mankali/reddit-mcp/cmd/rdcli/commands"

var (
	version = "0.1.0"
	commit  = "none"
	date    = "unknown"
)

func main() {
	commands.SetBuildInfo(version, commit, date)
	commands.Execute()
}
