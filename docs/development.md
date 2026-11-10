---
title: "Development"
description: "Build, lint and test gonfiguration with its Make targets, the 90 percent coverage gate and the vendored golangci-lint, and how CI and releases run."
nav_order: 6
---

# Development

Working on gonfiguration itself. You need Go 1.26 or newer and Make. golangci-lint comes from the module's own `tool` dependency in [`go.mod`](https://github.com/psyb0t/gonfiguration/blob/master/go.mod), and every dependency is vendored in [`vendor/`](https://github.com/psyb0t/gonfiguration/tree/master/vendor), so there is nothing to install globally.

## Make targets

Run `make help` for the list. The [`Makefile`](https://github.com/psyb0t/gonfiguration/blob/master/Makefile) has:

| Target | What it does |
| --- | --- |
| `make all` | Runs `dep`, `lint-fix` and `test-coverage`, in that order. |
| `make dep` | `go mod tidy` and `go mod vendor`. |
| `make lint` | Fails if `go fix -diff ./...` has suggestions, then runs `go tool golangci-lint run ./...` with the rules in [`.golangci.yml`](https://github.com/psyb0t/gonfiguration/blob/master/.golangci.yml). |
| `make lint-fix` | Applies `go fix ./...` and runs golangci-lint with `--fix`. |
| `make test` | `go test -race ./...`. |
| `make test-coverage` | Runs the tests with the race detector and a coverage profile, writes the total to `coverage-percent.txt`, and fails below 90 percent. |
| `make docs` | Builds this site into `.testing/docs-site` with GitHub's Pages build image. Needs Docker. |

[`pre-commit.sh`](https://github.com/psyb0t/gonfiguration/blob/master/pre-commit.sh) runs `make lint && make test-coverage`, the same pair CI runs.

No Go 1.26 on the machine? Run the same targets in the official image:

```bash
docker run --rm \
  --user "$(id -u):$(id -g)" \
  -e HOME=/tmp \
  -e GOCACHE=/repo/.testing/gocache \
  -v "$PWD:/repo" \
  -w /repo \
  golang:1.26 \
  make lint test-coverage
```

## Layout

| Path | Holds |
| --- | --- |
| [`gonfiguration.go`](https://github.com/psyb0t/gonfiguration/blob/master/gonfiguration.go) | `Parse`, `MustParse`, `GetAllValues`, `Reset`, the global state, tag parsing and the setter for each type. |
| [`defaults.go`](https://github.com/psyb0t/gonfiguration/blob/master/defaults.go) | `SetDefault`, `SetDefaults` and `GetDefaults`. |
| [`env_vars.go`](https://github.com/psyb0t/gonfiguration/blob/master/env_vars.go) | `GetEnvVars` and the environment cache. |
| [`errors.go`](https://github.com/psyb0t/gonfiguration/blob/master/errors.go) | The sentinel errors. |
| [`.example/main.go`](https://github.com/psyb0t/gonfiguration/blob/master/.example/main.go) | A runnable example. The leading dot keeps it out of `./...`, so lint and tests skip it. |
| [`docs/`](https://github.com/psyb0t/gonfiguration/tree/master/docs) | This site, built by GitHub Pages from the `master` branch. |

The black-box tests are in `gonfiguration_test.go` (`package gonfiguration_test`), and the tests of unexported helpers in `gonfiguration_internal_test.go`. The library has global state, so each test resets it with `defer gonfiguration.Reset()` and sets variables with `t.Setenv`.

## CI

[`.github/workflows/pipeline.yml`](https://github.com/psyb0t/gonfiguration/blob/master/.github/workflows/pipeline.yml) runs on every push to any branch or tag, once a week on a schedule, and by hand.

- On a branch push, the scheduled run and a manual run, the shared Go workflow runs `make lint`, `make test-coverage` and govulncheck.
- When that run is on `master`, the badges job then publishes the coverage, version, license and imported-by badges to the `badges` branch.
- On a tag, the pipeline does not rerun the checks. It waits for the `master` run of the same commit to pass, cuts the GitHub release with the annotated tag message as its notes, and refreshes the badges with the coverage from that `master` run. The `master` run keeps its coverage artifact for one day, so a tag pushed later than that still gets its release, but the badges job fails.

The weekly run exists for the imported-by badge, which pkg.go.dev updates on its own schedule.

## Releases

Versions follow [semver](https://semver.org), as the [CHANGELOG](https://github.com/psyb0t/gonfiguration/blob/master/CHANGELOG.md) says at the top. The Git tag is the version; there is no version file. A release is a commit on `master` that has passed CI, a CHANGELOG section for the version, and an annotated `vX.Y.Z` tag on that commit.
