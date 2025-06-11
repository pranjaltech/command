cask "cmd" do
  arch arm: "arm64", intel: "amd64"

  version :latest
  sha256 :no_check

  url "https://github.com/pranjaltech/command/releases/latest/download/cmd_darwin_#{arch}.tar.gz"

  name "cmd"
  desc "AI-assisted CLI that turns natural language into shell commands"
  homepage "https://github.com/pranjaltech/command"

  binary "cmd"

  postflight do
    if system_command("/usr/bin/xattr", args: ["-h"]).exit_status == 0
      system_command "/usr/bin/xattr", args: ["-dr", "com.apple.quarantine", "#{staged_path}/cmd"]
    end
  end
end
