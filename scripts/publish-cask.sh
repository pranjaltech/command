#!/usr/bin/env bash
set -euo pipefail

# Publish the GoReleaser-generated Homebrew cask to the tap.
#
# GoReleaser writes the cask to dist/ (homebrew_casks.skip_upload is "true"),
# but its template does not always match Homebrew's rubocop (e.g. the trailing
# blank line before `end` that Homebrew 6.0.3+ flags as
# Layout/EmptyLinesAroundBlockBody). The tap's CI runs `brew test-bot
# --only-tap-syntax` (i.e. `brew style`) and goes red on any offense, so we run
# `brew style --fix` here — inside a real tap checkout, where the Cask/* cops
# apply — before committing. This way the file that lands in the tap is already
# style-clean regardless of how GoReleaser's template evolves.
#
# Requires (provided by the release workflow):
#   - brew on PATH (Homebrew/actions/setup-homebrew)
#   - HOMEBREW_TAP_GITHUB_TOKEN with write access to the tap
#   - GITHUB_REF_NAME set to the release tag (falls back to `git describe`)

# `TAP` is the Homebrew short name (maps to the GitHub repo pranjaltech/
# homebrew-tools); `TAP_REPO` is that GitHub repo, used for the push URL.
TAP="pranjaltech/tools"
TAP_REPO="pranjaltech/homebrew-tools"
CASK_NAME="cmd"
GENERATED="dist/homebrew/Casks/${CASK_NAME}.rb"
TAG="${GITHUB_REF_NAME:-$(git describe --tags --abbrev=0)}"

if [[ ! -f "$GENERATED" ]]; then
    echo "::error::generated cask not found at ${GENERATED}" >&2
    exit 1
fi
if [[ -z "${HOMEBREW_TAP_GITHUB_TOKEN:-}" ]]; then
    echo "::error::HOMEBREW_TAP_GITHUB_TOKEN is not set" >&2
    exit 1
fi

# Clone the tap into Homebrew's taps directory. `brew style` only applies the
# Cask/* cops when the file lives inside a registered tap, so styling the loose
# dist/ file would lint it as plain Ruby and miss/misreport offenses.
brew tap "$TAP"
TAP_DIR="$(brew --repository "$TAP")"
CASK_PATH="${TAP_DIR}/Casks/${CASK_NAME}.rb"

cp "$GENERATED" "$CASK_PATH"

# Normalise to the current Homebrew rubocop, then hard-fail if anything is left
# uncorrected so a malformed cask never reaches the tap.
brew style --fix "$CASK_PATH"
brew style "$CASK_PATH"

cd "$TAP_DIR"
git config user.name "GitHub Actions"
git config user.email "actions@github.com"
git add "Casks/${CASK_NAME}.rb"

if git diff --cached --quiet; then
    echo "Cask ${CASK_NAME} is unchanged; nothing to publish."
    exit 0
fi

git commit -m "Brew cask update for ${CASK_NAME} version ${TAG}"
git remote set-url origin \
    "https://x-access-token:${HOMEBREW_TAP_GITHUB_TOKEN}@github.com/${TAP_REPO}.git"
git push origin HEAD:main

echo "Published style-clean cask for ${CASK_NAME} ${TAG} to ${TAP_REPO}."
