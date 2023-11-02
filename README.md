# gonfiguration 🔧

A kickass configuration package built on top of `viper` for Golang. Because, why the hell not? Simplify setting defaults and parsing configurations without the unnecessary bullshit.

## Installation

```bash
go get github.com/psyb0t/gonfiguration
```

## Usage

### Structuring your Config

Define your configuration struct. Ensure you use `mapstructure` tags because `gonfiguration` takes the good stuff from environment variables, and they must match these tags.

```go
type config struct {
	LogLevel  level  `mapstructure:"LOG_LEVEL"`
	LogFormat format `mapstructure:"LOG_FORMAT"`
}
```

### Set Defaults

```go
gonfiguration.SetDefaults(
	configuration.Default{
		Key:   "LOG_LEVEL",
		Value: "INFO", // Replace with your default level
	},
	configuration.Default{
		Key:   "LOG_FORMAT",
		Value: "text", // Replace with your default format
	},
)
```

### Parse Configurations

Parsing is a piece of cake. But for god's sake, handle those errors:

```go
c := config{}
if err := gonfiguration.Parse(&c); err != nil {
	panic("Damn! Couldn't parse the log config.")
}
```

## Development

### Makefile

Run these bad boys:

- `make dep` for dependency management.
- `make lint` to lint all your Golang files because nobody likes messy code.
- `make test` for the usual tests.
- `make test-coverage` to see if you're covering your ass enough with tests.

For a full list of commands:

```bash
make help
```

## Contributing

Got some wicked improvements or just found a dumb bug? Open a PR or shoot an issue. Let's get chaotic together.

## License

DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE. Seriously, just do what you want.
