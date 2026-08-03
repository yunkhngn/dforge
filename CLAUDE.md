# CLAUDE.md

Guidance for working in the **dforge** repository — an offline-first CLI that bootstraps,
inspects, and manages Docker Compose projects.

- Design spec: `docs/superpowers/specs/2026-08-03-dforge-design.md`
- Implementation plan: `docs/superpowers/plans/2026-08-03-dforge.md`
- Module path: `github.com/yunkhngn/dforge` · Go floor: `1.22`

## Everyday commands

```bash
go test ./...          # full suite (unit tests, no Docker daemon needed)
go build ./...         # build all packages
go vet ./...           # static checks
go run . <command>     # run the CLI locally, e.g. `go run . doctor`
```

Integration tests that need a real Docker daemon are gated behind a build tag and skipped
by default. Run them explicitly with `go test -tags=integration ./...`.

## Architecture (one-line)

Pure logic lives in `internal/` (unit-testable, no Docker/TTY); `cmd/` (Cobra) and
`internal/ui/` (Bubble Tea + Lipgloss) are thin shells; `internal/docker` is the only
impure package and sits behind a `Client` interface. See the design spec for details.

---

# CI/CD: automated build & release

Releases are **fully automated and tag-driven**. You do not build or upload binaries by
hand — pushing a semver tag does everything.

## How it works

1. A git tag matching `v*` (e.g. `v1.2.0`) is pushed to GitHub.
2. `.github/workflows/release.yml` triggers on that tag.
3. The workflow runs **GoReleaser** (`.goreleaser.yaml`), which:
   - builds cross-platform binaries (linux/darwin/windows × amd64/arm64),
   - creates archives + `checksums.txt`,
   - publishes a GitHub Release with those assets,
   - updates the Homebrew formula in the `yunkhngn/homebrew-tap` repo.

Nothing is published on pushes to branches or on pull requests — **only tags** cut a release.

## Prerequisites (one-time setup)

- A `yunkhngn/homebrew-tap` GitHub repository must exist.
- A repository/organization secret **`HOMEBREW_TAP_GITHUB_TOKEN`** must be set — a token
  with write access to `homebrew-tap`. Without it the build succeeds but the Homebrew
  formula update fails.
- `GITHUB_TOKEN` is provided automatically by Actions for the GitHub Release itself.

## Cutting a release (the process)

Do these in order. **Never tag until CI on `main` is green.**

```bash
# 1. Make sure everything passes locally and on main
go test ./... && go vet ./... && go build ./...

# 2. (Optional but recommended) dry-run the release locally
goreleaser check                                   # config is valid
goreleaser release --snapshot --clean --skip=publish   # builds into dist/, no publish

# 3. Tag with the chosen version and push the tag
git tag v1.2.0
git push origin v1.2.0
```

Then watch the `release` workflow in the Actions tab. When it finishes:

- verify the GitHub Release has all archives + `checksums.txt`,
- verify `brew install yunkhngn/tap/dforge` installs the new version,
- run `dforge --version` to confirm the version string matches the tag.

If the workflow fails **before publishing**, delete the tag, fix the issue, and re-tag:

```bash
git tag -d v1.2.0 && git push origin :refs/tags/v1.2.0
```

Do **not** reuse or move a tag that has already produced a published release — cut a new
patch version instead.

## Versioning (semver)

`vMAJOR.MINOR.PATCH`. Map changes to the version like this:

- **PATCH** (`v1.2.0 → v1.2.1`): bug fixes, doctor rule tweaks, template fixes, dependency
  bumps — no change to command behavior or flags.
- **MINOR** (`v1.2.0 → v1.3.0`): new command, new supported framework/service, new flag,
  or additive behavior that stays backward compatible.
- **MAJOR** (`v1.x → v2.0.0`): a breaking change — removed/renamed command or flag, changed
  output format that scripts may depend on, or a changed generated-file layout.

Pre-1.0 (`v0.x.y`): treat MINOR as the breaking-change lever and PATCH for everything else.

## When to release — decision guide

Release when **all** of these hold:

1. `main` is green: `go test ./...`, `go vet ./...`, `go build ./...` all pass in CI.
2. The change is user-facing and complete — a whole command, framework, service, or fix,
   not a half-finished slice.
3. The safety invariants still hold (these are release blockers if broken):
   - compose mutations preserve user comments and write a `.bak` backup,
   - `env` never prints variable **values**,
   - `clean` and any destructive action still require confirmation (or explicit `--force`),
   - generated Dockerfiles are multi-stage, non-root, pinned (no `:latest`), with a
     `HEALTHCHECK` — and `dforge doctor` on dforge's own generated output still scores high.
4. Docs reflect the change (README/spec/help text), so the release notes make sense.

**Do not release** when: CI is red or flaky; a feature is behind a feature-work branch not
yet merged; the only changes are internal refactors with no user impact (batch these into
the next feature release); or you are unsure whether a change is breaking — in that case,
size the version conservatively (round *up*: treat "maybe breaking" as MAJOR).

**Cadence:** release per meaningful unit of user value rather than on a fixed clock. A
single new command or framework is enough to justify a MINOR release; accumulate trivial
patches until there is something worth installing.

## Claude's role in releases

Claude may prepare and recommend releases (run the verification commands, draft release
notes, suggest the version bump) but must **not push a tag or trigger a release without the
user's explicit go-ahead** — publishing is outward-facing and irreversible.
