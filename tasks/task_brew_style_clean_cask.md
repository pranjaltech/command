Objective: Make the GoReleaser-generated Homebrew cask (`Casks/cmd.rb` in the
pranjaltech/homebrew-tools tap) pass `brew style` so the tap's
`brew test-bot --only-tap-syntax` CI is green. Homebrew 6.0.3 tightened rubocop
and the generated cask started failing.

Offenses to eliminate (Homebrew 6.0.3, `brew style pranjaltech/tools`):
- `Cask/StanzaOrder` — desc/homepage/version emitted out of canonical order.
- `Cask/StanzaGrouping` — stanza groups not separated by a single blank line.
- `Style/NumericPredicate` — `.exit_status == 0` in the postflight block.
- `Layout/EmptyLinesAroundBlockBody` — blank line before the final `end`
  (surfaced by current GoReleaser; introduced after the tap's last release).

Acceptance Criteria:
- The cask GoReleaser publishes to the tap is `brew style`-clean (0 offenses).
- The fix lives in this repo (not hand-edits to the generated, DO-NOT-EDIT cask).
- `goreleaser check` passes and the release workflow stays valid.

Findings:
- Current GoReleaser (v2.16.0, what `~> v2` resolves to) already emits canonical
  StanzaOrder/StanzaGrouping — those two offenses are gone with the new template.
- `Style/NumericPredicate` comes from this repo's `homebrew_casks.hooks.post.install`
  snippet; fixed at the source by using `.exit_status.zero?`.
- `Layout/EmptyLinesAroundBlockBody` (the trailing blank before `end`) is hardcoded
  in GoReleaser's template and cannot be removed via config (adding `zap`/stanzas to
  fill the space only introduces new offenses). It must be post-processed.

Implementation Checklist:
- [x] `.goreleaser.yaml`: use `.exit_status.zero?` in the postflight hook.
- [x] `.goreleaser.yaml`: set `homebrew_casks.skip_upload: "true"` so GoReleaser
      writes the cask to dist/ but does not push the un-styled file.
- [x] Add `scripts/publish-cask.sh`: style the generated cask with `brew style --fix`
      inside a real tap checkout (where the Cask/* cops apply), hard-fail on any
      remaining offense, then commit/push it to the tap.
- [x] `.github/workflows/release.yml`: add `Homebrew/actions/setup-homebrew` (same as
      the tap CI) and run the publish script with `HOMEBREW_TAP_GITHUB_TOKEN`.
- [x] Verify: `goreleaser release --snapshot` then `brew style --fix` in tap context
      yields `no offenses detected`.
- [x] Verify: full clone→copy→style→commit→push simulation against a local bare
      remote; the pushed cask re-styles clean.

Status: Completed
