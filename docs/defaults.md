---
title: "Defaults and precedence"
description: "Set gonfiguration fallbacks with the default struct tag or SetDefault and SetDefaults, and see which value wins when several are present."
nav_order: 3
---

# Defaults and precedence

A field can get its value from three places. `Parse` applies them in this order, and the last one present wins:

1. the `default` struct tag,
2. a programmatic default set with `SetDefault` or `SetDefaults`,
3. the environment variable.

So an environment variable always beats a default, and a default set from code beats the one in the tag. A `required` field fails only when none of the three supplied a value. The order is in `fillFieldValue` in [`gonfiguration.go`](https://github.com/psyb0t/gonfiguration/blob/master/gonfiguration.go).

## The default tag

The simplest way is to put the fallback right next to the field.

```go
type Config struct {
	Port         int           `env:"PORT" default:"8080"`
	Host         string        `env:"HOST" default:"localhost"`
	Timeout      time.Duration `env:"TIMEOUT" default:"30s"`
	AllowedHosts []string      `env:"ALLOWED_HOSTS" default:"localhost,127.0.0.1"`
	Rate         float64       `env:"RATE" default:"0.5"`
	APIKey       string        `env:"API_KEY,required" default:"dev-key"`
}
```

The tag value is a string parsed by the same rules as an environment value, described on [Struct tags and types](types.md). It is parsed on every `Parse`, even when the environment variable is set, so a tag value that does not parse fails `Parse` either way:

```
failed to parse fields: failed to set field value: field T: invalid default tag value "abc": failed to parse int: strconv.ParseInt: parsing "abc": invalid syntax [...]
```

A `default` tag also satisfies `required`, so `APIKey` above never fails on a missing variable.

## Defaults from code

Use `SetDefault` or `SetDefaults` when the fallback lives in a shared package or has to be computed at run time. Call them before `Parse`.

```go
gonfiguration.SetDefault("PORT", 8080)
gonfiguration.SetDefault("TIMEOUT", 30*time.Second)

gonfiguration.SetDefaults(map[string]any{
	"HOST":          "localhost",
	"ALLOWED_HOSTS": []string{"localhost", "127.0.0.1"},
})
```

The key is the environment variable name from the `env` tag, not the Go field name.

These values are Go values, not strings, and their dynamic type must be exactly the field's type. `8080` is an untyped constant that becomes an `int`, so it works for an `int` field and fails for an `int64` or `uint16` field, which need `int64(8080)` and `uint16(8080)`. A `time.Duration` field needs a `time.Duration`, such as `30*time.Second`, not a number or the string `"30s"`. A mismatch fails `Parse` with `ErrDefaultTypeMismatch`, even when the environment variable is set:

```
failed to parse fields: failed to set field value: default value type mismatch [...]
```

That error does not name the key, so check the types of every `SetDefault` call when you see it.

Defaults are process-wide and stay until you overwrite them or call `Reset`. `SetDefaults` adds to the existing defaults; it does not clear the ones it does not mention. Setting a key to `nil` makes `Parse` treat it as having no programmatic default.

## Reading the values back

| Function | Returns |
| --- | --- |
| `GetDefaults()` | A copy of the programmatic defaults. `default` tags are not included, because they belong to a struct, not to the library. |
| `GetEnvVars()` | A copy of the environment cached by the last `Parse` calls. |
| `GetAllValues()` | The programmatic defaults with the cached environment laid over them. |

All three return copies, so changing the map you get back changes nothing inside the library. The cache holds the whole process environment, not only the variables your struct names, so do not log these maps as they are when the environment holds secrets. More on the cache in the [API reference](api.md#global-state).
