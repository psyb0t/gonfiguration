# Errors

Prefix: ERR.

- **ERR1:** Sentinels live in `errors.go` as `errors.New` values named `Err...`. Callers match them with `errors.Is`, so a sentinel is never renamed, removed or given a different meaning.
- **ERR2:** Every error leaving a function is wrapped with `ctxerrors.Wrap` or `ctxerrors.Wrapf` and a short phrase saying what failed, as `Parse` and the `set*` helpers do. The wrap adds the file, line and function, and `errors.Is` still matches the sentinel underneath.
- **ERR3:** Errors that name a field use the env var key (`field %s` in `fillFieldValue`), because that is what the caller has to fix.
- **ERR4:** Check `ctxerrors/commerr` before adding a sentinel. Add a new one to `errors.go` only when no existing sentinel or `commerr` class fits, and list it in `docs/errors.md`.
- **ERR5:** A value can be a secret. The `strconv` and `time.ParseDuration` errors the setters wrap already quote the value that failed to parse; do not add any other place where an environment variable's value goes into an error message. Name the key instead. The `default` tag value quoted in `fillFieldValue` comes from source code, not the environment, so it is fine.
