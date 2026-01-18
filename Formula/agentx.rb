# typed: false
# frozen_string_literal: true

class Agentx < Formula
  desc "CLI tool for managing MCP servers and skills across AI coding agents"
  homepage "https://github.com/agentsdance/agentx"
  url "https://github.com/agentsdance/agentx/archive/refs/tags/v0.0.6.tar.gz"
  sha256 "60231952b4d070f2256a3e0dec39a4093d33d36920d7bec1342ff5f4508e2142"
  license "Apache-2.0"
  head "https://github.com/agentsdance/agentx.git", branch: "master"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X github.com/agentsdance/agentx/internal/version.Version=#{version}
    ]
    system "go", "build", *std_go_args(ldflags: ldflags)
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/agentx version")
  end
end
