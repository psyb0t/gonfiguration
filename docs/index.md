---
title: "Home"
description: "gonfiguration is a Go library that fills a config struct from environment variables through env and default struct tags, with required fields."
nav_order: 1
permalink: /
seo:
  type: WebSite
---

# gonfiguration

A thread-safe Go library that fills a config struct from environment variables, so you stop writing the same `os.Getenv` and `strconv` boilerplate in every service.

Tag the fields you want with `env:"NAME"`, mark the ones that must be set with `,required`, give fallbacks with a `default` tag or from code, and call `Parse`. Errors are wrapped with [ctxerrors](https://github.com/psyb0t/ctxerrors), so each one carries the file, line and function of every hop it passed through.

## Install

```bash
go get github.com/psyb0t/gonfiguration
```

The module declares `go 1.26` in its [`go.mod`](https://github.com/psyb0t/gonfiguration/blob/master/go.mod), so you need Go 1.26 or newer. Besides ctxerrors, the library imports only the standard library.

## Quick look

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/psyb0t/gonfiguration"
)

type Config struct {
	ListenAddress string        `env:"LISTEN_ADDRESS" default:"127.0.0.1:8080"`
	Debug         bool          `env:"DEBUG" default:"false"`
	Timeout       time.Duration `env:"TIMEOUT" default:"30s"`
	AllowedHosts  []string      `env:"ALLOWED_HOSTS" default:"localhost,127.0.0.1"`
	DBName        string        `env:"DB_NAME,required"`
}

func main() {
	var cfg Config
	if err := gonfiguration.Parse(&cfg); err != nil {
		log.Fatalf("config: %v", err)
	}

	fmt.Printf("%+v\n", cfg)
}
```

Run it with `DB_NAME=myapp ALLOWED_HOSTS="api.example.com, cdn.example.com"` and `AllowedHosts` comes out as `[api.example.com cdn.example.com]`, the spaces trimmed. Leave `DB_NAME` unset and `Parse` returns an error that matches `gonfiguration.ErrRequiredFieldNotSet`.

## Where to go next

| You want to | Read |
| --- | --- |
| Know which tags exist and how each field type is parsed | [Struct tags and types](types.md) |
| Set fallbacks and know which value wins | [Defaults and precedence](defaults.md) |
| Match errors with `errors.Is` and read what they print | [Errors](errors.md) |
| Look up a signature, or use the library from several goroutines | [API reference](api.md) |
| Build, lint and test the library itself | [Development](development.md) |

The Go package docs are on [pkg.go.dev](https://pkg.go.dev/github.com/psyb0t/gonfiguration), and the source is on [GitHub](https://github.com/psyb0t/gonfiguration). A longer runnable example is in [`.example/main.go`](https://github.com/psyb0t/gonfiguration/blob/master/.example/main.go).
