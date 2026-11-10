# Testing

Prefix: TEST.

- **TEST1:** `make test-coverage` is the CI gate: `go test -race` with a coverage profile, failing below `MIN_TEST_COVERAGE` (90) in the `Makefile`. Cover new code instead of lowering the floor.
- **TEST2:** Black-box tests go in `gonfiguration_test.go` (`package gonfiguration_test`) and use only the exported API. White-box tests of unexported helpers go in `gonfiguration_internal_test.go`.
- **TEST3:** The library keeps global state, so every test that sets defaults or calls `Parse` starts with `defer gonfiguration.Reset()`, and environment variables are set with `t.Setenv` so they are restored afterwards. A test that has to clear the whole environment restores it in a `defer`, as `TestInvalidEnvVariable` does. Tests do not call `t.Parallel()`: `t.Setenv` panics in a parallel test, and the state is shared.
- **TEST4:** A new supported type gets a success case, a parse error case and a `default` tag case.
- **TEST5:** Tests use testify (`require`).
