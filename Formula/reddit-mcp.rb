class RedditMcp < Formula
  desc "Reddit CLI + MCP tool for browsing, posting, search, and automation"
  homepage "https://github.com/virat-mankali/reddit-mcp"
  version "1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/virat-mankali/reddit-mcp/releases/download/v#{version}/reddit-mcp_Darwin_arm64.tar.gz"
      sha256 "REPLACED_BY_GORELEASER"
    else
      url "https://github.com/virat-mankali/reddit-mcp/releases/download/v#{version}/reddit-mcp_Darwin_x86_64.tar.gz"
      sha256 "REPLACED_BY_GORELEASER"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/virat-mankali/reddit-mcp/releases/download/v#{version}/reddit-mcp_Linux_arm64.tar.gz"
      sha256 "REPLACED_BY_GORELEASER"
    else
      url "https://github.com/virat-mankali/reddit-mcp/releases/download/v#{version}/reddit-mcp_Linux_x86_64.tar.gz"
      sha256 "REPLACED_BY_GORELEASER"
    end
  end

  def install
    bin.install "rdcli"
  end

  test do
    assert_match "rdcli version", shell_output("#{bin}/rdcli --version")
  end
end
