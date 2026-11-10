---
title: "API reference"
description: "Every exported gonfiguration function with its signature and behavior, plus how the process-wide state and its locking work across goroutines."
nav_order: 5
---

# API reference

Everything `github.com/psyb0t/gonfiguration` exports, apart from the sentinel errors, which are on [Errors](errors.md). The same list, generated from the source, is on [pkg.go.dev](https://pkg.go.dev/github.com/psyb0t/gonfiguration).

## Parsing

| Signature | Behavior |
| --- | --- |
| `func Parse(dst any) error` | Reads `os.Environ()` into the cache, checks that `dst` is a non-nil pointer to a struct, then fills every field with an `env` tag. Fields are filled in declaration order and the first error stops it, so fields before the failing one are already set. |
| `func MustParse(dst any)` | `Parse`, then `panic(err)` when it fails. For `main` or `init` code that should not start with bad config. |

```go
var cfg Config
if err := gonfiguration.Parse(&cfg); err != nil {
	log.Fatalf("config: %v", err)
}

var other Config
gonfiguration.MustParse(&other)
```

`Parse` only writes fields that get a value. A field with no environment variable and no default keeps whatever it held before the call, so you can prefill a struct and let the environment override parts of it.

## Defaults

| Signature | Behavior |
| --- | --- |
| `func SetDefault(key string, val any)` | Stores `val` as the programmatic default for the environment variable `key`, replacing any earlier one. |
| `func SetDefaults(defaults map[string]any)` | `SetDefault` for each entry. Keys it does not mention keep their defaults. |
| `func GetDefaults() map[string]any` | A copy of the programmatic defaults. |

How defaults combine with `default` tags and the environment, and the type rule they follow, is on [Defaults and precedence](defaults.md).

## Reading state

| Signature | Behavior |
| --- | --- |
| `func GetEnvVars() map[string]string` | A copy of the cached environment. Empty until the first `Parse`. |
| `func GetAllValues() map[string]any` | A copy of the programmatic defaults with the cached environment written over them. The environment values stay strings. |
| `func Reset()` | Drops every default and the cached environment. |

## Global state

The library keeps one process-wide value, created when the package loads, holding two maps: the programmatic defaults and the cached environment. Every function above works on it, so two packages in the same program that both call `SetDefault` share the same defaults.

The environment cache is filled by `Parse`, from all of `os.Environ()`, not only the variables your struct names. Each `Parse` adds and overwrites entries but never removes one, so a variable unset after an earlier `Parse` stays in `GetEnvVars()` until `Reset`. `Parse` itself always reads the live environment, so this only affects what the getters return.

## Concurrency

The maps are guarded by a `sync.RWMutex`, and the getters return copies. `Parse`, `MustParse`, `SetDefault`, `SetDefaults`, `GetDefaults`, `GetEnvVars` and `GetAllValues` are safe to call from several goroutines at once.

```go
go func() {
	gonfiguration.SetDefault("KEY1", "value1")
}()

go func() {
	var cfg Config
	_ = gonfiguration.Parse(&cfg)
}()
```

Two limits:

- `Reset` is not safe alongside other calls. It swaps out the whole state value, and that swap is outside the lock, so `go test -race` flags it when another goroutine is using the library at the same moment. Use it between tests, not while the program runs.
- Each lock covers one map operation, not a whole call. `SetDefaults` locks once per key, and `Parse` reads each default under its own lock, so a `Parse` running at the same time as a `SetDefaults` can see some of the new defaults and not others. Set defaults before you start parsing.
