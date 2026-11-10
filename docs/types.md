---
title: "Struct tags and types"
description: "The env, required and default struct tags in gonfiguration, every supported field type, and exactly how each one is parsed from a string."
nav_order: 2
---

# Struct tags and types

`Parse` walks the fields of the struct you hand it and fills each one that has an `env` tag. This page covers the tags it reads, the field types it accepts and how each type is parsed. Which value wins when several are present is on [Defaults and precedence](defaults.md).

## Tags

| Tag | Meaning |
| --- | --- |
| `env:"NAME"` | Fill this field from the environment variable `NAME`. |
| `env:"NAME,required"` | Same, and fail with `ErrRequiredFieldNotSet` when no value comes from anywhere. |
| `default:"value"` | Fallback used when `NAME` is not set. Parsed exactly like an environment value. |

```go
type Config struct {
	Port    int    `env:"PORT" default:"8080"`
	APIKey  string `env:"API_KEY,required"`
	Verbose bool   `env:"VERBOSE"`
	Ignored string
}
```

The rules, straight from [`gonfiguration.go`](https://github.com/psyb0t/gonfiguration/blob/master/gonfiguration.go):

- A field without an `env` tag is skipped, whatever its type. `Ignored` above is never touched.
- The name and options are split on commas and trimmed, so `env:"PORT, required"` works. Any option other than `required` is ignored without an error.
- `required` is satisfied by an environment variable that is set, even to an empty string, by a `default` tag, or by a programmatic default.
- A tagged field must be exported. An unexported field with an `env` tag makes `Parse` panic (`reflect: reflect.Value.SetString using value obtained using unexported field`) as soon as a value is applied to it.
- Only top-level fields are read. A nested struct, a pointer, a map or any other unsupported type with an `env` tag fails the whole `Parse` with `ErrUnsupportedFieldType`.

## Supported types

| Field type | Parsed with | Accepts |
| --- | --- | --- |
| `string` | nothing | The value as it is, including an empty string. |
| `bool` | `strconv.ParseBool` | `1`, `t`, `T`, `true`, `TRUE`, `True`, `0`, `f`, `F`, `false`, `FALSE`, `False`. |
| `int`, `int8`, `int16`, `int32`, `int64` | `strconv.ParseInt(value, 10, 64)` | Base 10 integers such as `-42`. No `0x` prefix. |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `strconv.ParseUint(value, 10, 64)` | Base 10 integers without a sign, such as `42`. |
| `float32`, `float64` | `strconv.ParseFloat` with the field's bit size | `3.14`, `1e-3`, and the other forms `ParseFloat` takes. |
| `time.Duration` | `time.ParseDuration` | `300ms`, `30s`, `5m`, `1h30m`. A bare number such as `30` is rejected. |
| `[]string` | split on `,`, each item trimmed | `a, b ,c` becomes `["a", "b", "c"]`. |

A named type over one of the scalar kinds, such as `type Mode string` or `type Port uint16`, is accepted and parsed by its kind. `time.Duration` is matched by type before its `int64` kind is looked at, so it gets `time.ParseDuration`. Slices are matched by exact type, so only `[]string` works; a named slice type such as `type Hosts []string` fails with `ErrUnsupportedFieldType`.

### Parsing details worth knowing

- **Narrow integers wrap without an error.** Every integer is parsed as 64 bits and then stored in the field. A value that does not fit the field's width is truncated: `300` in an `int8` becomes `44`, and `256` in a `uint8` becomes `0`. Only values outside the 64-bit range fail. Validate the range yourself when it matters.
- **An empty value is not "unset".** `PORT=` sets `PORT` to an empty string. That overrides every default, so an `int` field fails because `strconv.ParseInt` rejects `""`. A `string` field becomes empty, and a `[]string` field becomes an empty slice, `[]string{}`.
- **`[]string` has no escaping.** A comma always splits, and empty items stay: `a,,b` becomes `["a", "", "b"]`.
- **Errors quote the value.** The `strconv` and `time.ParseDuration` errors include the text that failed to parse, so a typo in a secret ends up in the error message. See [Errors](errors.md).
