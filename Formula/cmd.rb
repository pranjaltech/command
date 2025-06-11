class Cmd < Formula
  desc "AI-assisted CLI that turns natural language into shell commands"
  homepage "https://github.com/pranjaltech/command"
  head "https://github.com/pranjaltech/command.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./"
  end

  test do
    assert_match "A prompt is required", shell_output("#{bin}/cmd 2>&1", 1)
  end
end
