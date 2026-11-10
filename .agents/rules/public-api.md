# Public API

Prefix: API. Every exported name is in use by other modules, so the surface changes only on purpose.

- **API1:** The exported surface is eight functions (`Parse`, `MustParse`, `SetDefault`, `SetDefaults`, `GetDefaults`, `GetEnvVars`, `GetAllValues`, `Reset`) and the seven sentinels in `errors.go`. Do not rename, remove or change the signature of any of them in a minor or patch release.
- **API2:** The struct tag contract is `env:"NAME"`, `env:"NAME,required"` and `default:"value"` (`parseTag` and `tagDefaultFromField` in `gonfiguration.go`). A field without an `env` tag is skipped, whatever its type. Unknown options after the name are ignored today; giving one a meaning is a behavior change for every caller.
- **API3:** The precedence is `default` tag, then programmatic default, then environment variable, the last one present winning (`fillFieldValue`). A variable that is set but empty still counts as present. Keep both facts unless a release says otherwise.
- **API4:** The supported field types are `string`, `bool`, every `int` and `uint` width, `float32`, `float64`, `time.Duration` and `[]string`. A new type is added in both `isSupportedType` and `setEnvVarValue`, in the same change, with tests and a line in `docs/types.md`.
- **API5:** Every exported function goes through the `gonfig` value and its mutex. New state goes into the `gonfiguration` struct behind the same lock, never into a new package-level variable. `Reset` is the one gap: it replaces the `gonfig` pointer itself, which no lock covers, so `go test -race` reports a race when it runs alongside any other call. Treat it as test-only and do not build on it.
- **API6:** A change to exported behavior updates `docs/` (the API reference is `docs/api.md`) and the CHANGELOG entry of the release that ships it.
