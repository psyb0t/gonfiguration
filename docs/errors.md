---
title: "Errors"
description: "The seven gonfiguration sentinel errors, when Parse returns each one, how to match them with errors.Is, and what the wrapped messages print."
nav_order: 4
---

# Errors

`Parse` returns an error wrapped with [ctxerrors](https://github.com/psyb0t/ctxerrors) at every hop, so the message says what failed and the file, line and function of each wrap. The sentinels underneath are plain `errors.New` values from [`errors.go`](https://github.com/psyb0t/gonfiguration/blob/master/errors.go), and `errors.Is` matches them through all the wrapping. `MustParse` panics with the same error.

## Sentinels

| Sentinel | Message | Returned when |
| --- | --- | --- |
| `ErrNilDestination` | `destination is nil` | `Parse(nil)`. |
| `ErrTargetNotPointer` | `destination must be a pointer` | You passed the struct itself instead of a pointer to it. |
| `ErrDestinationNotStruct` | `destination must be a struct` | The pointer does not point at a struct. |
| `ErrUnsupportedFieldType` | `unsupported field type` | A field with an `env` tag has a type from outside the [supported list](types.md#supported-types). |
| `ErrRequiredFieldNotSet` | `required field not set` | A `required` field got no value from the environment, a `default` tag or a programmatic default. |
| `ErrDefaultTypeMismatch` | `default value type mismatch` | A [programmatic default](defaults.md#defaults-from-code) has a different Go type than its field. |
| `ErrInvalidEnvVar` | `invalid environment variable` | An `os.Environ()` entry has no `=`. The Go runtime does not hand out entries like that in normal use, so this is close to unreachable. |

A value that fails to parse (`PORT=abc` for an `int` field) has no sentinel of its own. The error wraps the `strconv` or `time` error, so `errors.Is(err, strconv.ErrSyntax)` and `errors.Is(err, strconv.ErrRange)` work, as does `errors.As` with `*strconv.NumError`.

```go
err := gonfiguration.Parse(&cfg)

switch {
case err == nil:
case errors.Is(err, gonfiguration.ErrRequiredFieldNotSet):
	log.Fatalf("missing config: %v", err)
case errors.Is(err, strconv.ErrSyntax):
	log.Fatalf("bad config value: %v", err)
default:
	log.Fatalf("config: %v", err)
}
```

## What the messages look like

Each wrap adds a phrase to the front and a `[file:line in function]` location to the end. These are real outputs from the current code. The file paths depend on where the module sits on your machine, shortened here to `/path/to`.

Passing the struct instead of a pointer:

```
invalid destination: destination must be a pointer [/path/to/gonfiguration.go:40 in github.com/psyb0t/gonfiguration.Parse]
```

A missing `required` field. This one names the variable:

```
failed to parse fields: failed to set field value: field API_KEY: required field not set [/path/to/gonfiguration.go:169 in github.com/psyb0t/gonfiguration.fillFieldValue] [/path/to/gonfiguration.go:112 in github.com/psyb0t/gonfiguration.parseDstFields] [/path/to/gonfiguration.go:44 in github.com/psyb0t/gonfiguration.Parse]
```

A value that does not parse. This one names the setter and quotes the value, but not the variable name:

```
failed to parse fields: failed to set field value: failed to parse int: strconv.ParseInt: parsing "not-a-number": invalid syntax [/path/to/gonfiguration.go:238 in github.com/psyb0t/gonfiguration.setInt] [/path/to/gonfiguration.go:112 in github.com/psyb0t/gonfiguration.parseDstFields] [/path/to/gonfiguration.go:44 in github.com/psyb0t/gonfiguration.Parse]
```

An unsupported field type stops at the first such field and names neither the field nor the type:

```
failed to parse fields: unsupported field type [/path/to/gonfiguration.go:44 in github.com/psyb0t/gonfiguration.Parse]
```

Two things follow from those outputs:

- A parse failure quotes the bad value. If the variable holds a secret, the secret is in the error, so be careful where you log `Parse` errors.
- Only `required` and `default` tag errors name the variable. For the others, the quoted value and the setter (`setInt`, `setBool`, `setDuration`, ...) are what you have to go on.

How ctxerrors formats the chain is on its [error output page](https://psyb0t.github.io/ctxerrors/error-output/).
