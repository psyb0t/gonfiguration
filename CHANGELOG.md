# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v1.5.3 — 2026-07-27

Self-hosted README badges.

- **Coverage / version / license badges are self-rendered SVGs** served from
  `raw.githubusercontent.com/psyb0t/gonfiguration/badges/*.svg` — no third-party
  render service. `make test-coverage` writes the percentage to
  `coverage-percent.txt`, the pipeline uploads it, and a `badges` job bakes it
  into the SVG. The CI badge is switched to GitHub's native `badge.svg`. No
  library code changed.

## v1.5.2 — 2026-07-27

- Bump `github.com/stretchr/testify` 1.10.0 → 1.11.1 (test dependency).

## v1.5.1 — 2026-07-26

README badges + repo housekeeping.

- pkg.go.dev reference + GitHub Actions CI status badges.
- Added a GitHub Sponsors funding config; CI restricted to collaborators only;
  README tweaks. No library code changed.

## v1.5.0 — 2026-03-13

- Added default-value support via a struct tag.

## v1.4.1 — 2026-01-16

- Modernized tooling and updated the Go version.

## v1.4.0 — 2026-01-16

- Added required-field support via `env:"VAR,required"`.
- Added `MustParse()` that panics on error.
- Added `ErrNilDestination`, `ErrRequiredFieldNotSet`, `ErrDefaultTypeMismatch`.
- Removed `pkg/errors`; use stdlib `fmt.Errorf` with `%w`. Cleaned up error
  messages.

## v1.3.1 — 2025-09-11

- Maintenance release.

## v1.3.0 — 2025-09-11

- Added `[]string` support.

## v1.2.0 — 2025-09-07

- Added `time.Duration` support.

## v1.1.0 — 2025-09-07

- Maintenance release.

## v1.0.0 — 2023-11-04

- Initial release.
