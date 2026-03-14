# envstruct

A Go library that populates struct fields from environment variables. Drop-in replacement for [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig).

## Install

```
go get github.com/agentine/envstruct
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/agentine/envstruct"
)

type Config struct {
    Host  string `default:"localhost" desc:"Server hostname"`
    Port  int    `default:"8080"     desc:"Server port"`
    Debug bool   `                   desc:"Enable debug mode"`
}

func main() {
    var c Config
    if err := envstruct.Process("APP", &c); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Listening on %s:%d\n", c.Host, c.Port)
}
```

Set environment variables and run:

```
APP_HOST=0.0.0.0 APP_PORT=9090 go run main.go
# Listening on 0.0.0.0:9090
```

## Supported Types

| Type | Example env value |
|------|-------------------|
| `string` | `hello` |
| `bool` | `true`, `false`, `1`, `0` |
| `int`, `int8`..`int64` | `42`, `-1` |
| `uint`, `uint8`..`uint64` | `42` |
| `float32`, `float64` | `3.14` |
| `time.Duration` | `5s`, `100ms` |
| `url.URL` / `*url.URL` | `https://example.com` |
| `[]T` (any scalar T) | `a,b,c` |
| `map[string]T` | `key1=val1,key2=val2` |
| Custom `Decoder` | User-defined |
| Custom `Setter` | envconfig compat |
| `encoding.TextUnmarshaler` | User-defined |

## Struct Tags

| Tag | Description |
|-----|-------------|
| `env:"VAR_NAME"` | Override env var name |
| `env:"VAR_NAME,required"` | Mark field as required |
| `env:"-"` | Skip field |
| `envconfig:"VAR_NAME"` | envconfig compat tag |
| `default:"value"` | Default if env var unset |
| `desc:"description"` | Description for Usage() |

## Nested Structs

Nested struct fields are flattened with `_` separators:

```go
type DB struct {
    Host string
    Port int
}
type Config struct {
    Database DB
}
// Reads APP_DATABASE_HOST, APP_DATABASE_PORT
```

Embedded structs are flattened without adding a prefix segment.

## Usage Text

```go
envstruct.Usage("APP", &Config{}, os.Stderr)
```

Outputs:

```
  APP_HOST   string  [default: localhost]  Server hostname
  APP_PORT   int     [default: 8080]       Server port
  APP_DEBUG  bool                          Enable debug mode
```

## Migration from envconfig

envstruct is a drop-in replacement. The function signatures are identical:

```go
// Before
envconfig.Process("APP", &config)
envconfig.MustProcess("APP", &config)

// After
envstruct.Process("APP", &config)
envstruct.MustProcess("APP", &config)
```

Both `env` and `envconfig` struct tags are supported for smooth migration.

## License

MIT
