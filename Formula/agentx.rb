# typed: false
# frozen_string_literal: true

class Agentx < Formula
  desc "CLI tool for managing MCP servers and skills across AI coding agents"
  homepage "https://github.com/agentsdance/agentx"
  url "https://github.com/agentsdance/agentx/archive/refs/tags/v0.1.1.tar.gz"
  sha256 "d96a3bca1ee5ecf15c7311ee92f91eb15e3f16f6ed7552f620e8bf6190342b9e"
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
