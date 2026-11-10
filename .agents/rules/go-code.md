# Go code

Prefix: GO.

- **GO1:** The library imports the standard library and `github.com/psyb0t/ctxerrors`, nothing else. `depguard` in `.golangci.yml` allows `$gostd`, `github.com/psyb0t`, testify, logrus and `github.com/pkg/errors`; the last two are not used and are no reason to start. A new third-party dependency needs a reason that outweighs adding it to every importer's module graph.
- **GO3:** Test files are excluded from golangci-lint in `.golangci.yml`.
- **GO5:** Reflection setters write through `reflect.Value.Set*`. A field must be exported for that to work; an unexported field with an `env` tag panics as soon as a value is applied to it. Keep any new setter on the same path so the behavior stays the same for every type.
