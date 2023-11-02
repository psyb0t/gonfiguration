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

Copyright 2023 Ciprian Mandache ([ciprian.51k.eu](https://ciprian.51k.eu))

Listen up! Permission is straight-up given, no strings attached, to any badass out there snagging a copy of this masterpiece (let's call it the "Software"). You can rock out with the Software any damn way you please. Want to use it? Go for it. Modify it? Be my guest. Merge, publish, distribute, sublicense, or even make a quick buck selling it? Hell yeah, you can. Just if you're handing this gem to someone else, don't be a douche – include this copyright notice and my cool permission ramble in all copies or major parts of the Software.

Now, here's the kicker: the Software is provided "as is". I ain't making any pinky promises on how it'll perform or if it might royally screw things up. So, if some shit hits the fan, don't come crying to me or any other folks holding the copyright. We're just chilling and ain't responsible for whatever chaos you or this code might stir up.
