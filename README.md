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

## API Reference

### Functions

#### `Process(prefix string, spec interface{}) error`

Populates the struct pointed to by `spec` with values from environment variables. The `prefix` is prepended to each field's env var name (uppercased, joined with `_`).

Returns a joined error of all field errors encountered. Returns an error immediately if `spec` is not a non-nil pointer to a struct.

```go
var cfg Config
if err := envstruct.Process("MYAPP", &cfg); err != nil {
    log.Fatal(err)
}
```

#### `MustProcess(prefix string, spec interface{})`

Like `Process` but panics on error. Useful in `main()` when a misconfigured environment is always fatal.

```go
var cfg Config
envstruct.MustProcess("MYAPP", &cfg)
```

#### `Usage(prefix string, spec interface{}, out io.Writer) error`

Writes a human-readable table of expected environment variables to `out`, including their type, default value, and description.

```go
envstruct.Usage("APP", &Config{}, os.Stderr)
```

Output:

```
  APP_HOST   string  [default: localhost]  Server hostname
  APP_PORT   int     [default: 8080]       Server port
  APP_DEBUG  bool                          Enable debug mode
```

### Interfaces

Types can implement any of the following interfaces to control their own decoding. Interface methods are checked via pointer receivers.

#### `Decoder`

```go
type Decoder interface {
    Decode(value string) error
}
```

The primary custom decoding interface. If a field's type implements `Decoder`, `Decode` is called with the raw environment variable string.

#### `Setter`

```go
type Setter interface {
    Set(value string) error
}
```

Compatibility interface for types that implement `kelseyhightower/envconfig`'s `Setter` interface. Checked after `Decoder`.

#### `encoding.TextUnmarshaler`

Standard library interface. If a type implements `encoding.TextUnmarshaler`, `UnmarshalText` is called. Checked after `Setter`.

**Custom type example:**

```go
type LogLevel int

const (
    LevelInfo LogLevel = iota
    LevelWarn
    LevelError
)

func (l *LogLevel) Decode(s string) error {
    switch strings.ToLower(s) {
    case "info":  *l = LevelInfo
    case "warn":  *l = LevelWarn
    case "error": *l = LevelError
    default:      return fmt.Errorf("unknown level: %s", s)
    }
    return nil
}

type Config struct {
    Level LogLevel `default:"info"`
}
// LOG_LEVEL=warn envstruct.Process("LOG", &cfg)
```

### Error Types

#### `ParseError`

Returned when an environment variable value cannot be parsed into the target type.

```go
type ParseError struct {
    FieldName string  // struct field name
    EnvVar    string  // full env var key (e.g. APP_PORT)
    Value     string  // the raw string value that failed to parse
    TypeName  string  // target Go type name
    Err       error   // underlying parse error
}
```

`ParseError` implements `Unwrap()`, so it works with `errors.As` and `errors.Is`.

#### `RequiredError`

Returned when a field marked `required` has no environment variable set and no default.

```go
type RequiredError struct {
    FieldName string  // struct field name
    EnvVar    string  // full env var key
}
```

## Supported Types

| Type | Example env value |
|------|-------------------|
| `string` | `hello` |
| `bool` | `true`, `false`, `1`, `0` |
| `int`, `int8`..`int64` | `42`, `-1` |
| `uint`, `uint8`..`uint64` | `42` |
| `float32`, `float64` | `3.14` |
| `time.Duration` | `5s`, `100ms`, `1h30m` |
| `url.URL` / `*url.URL` | `https://example.com` |
| `[]T` (any scalar T) | `a,b,c` |
| `map[string]T` | `key1=val1,key2=val2` |
| Custom `Decoder` | User-defined |
| Custom `Setter` | envconfig compat |
| `encoding.TextUnmarshaler` | User-defined |

### Slices

Slice fields are split on `,` by default. Use the `envSeparator` tag to override the delimiter.

```go
type Config struct {
    Tags    []string `desc:"Comma-separated tags"`         // TAGS=a,b,c
    Paths   []string `envSeparator:":" desc:"Colon paths"` // PATHS=/a:/b:/c
}
```

### Maps

Map fields are split on `,`, with each element parsed as `key=value`:

```go
type Config struct {
    Labels map[string]string // LABELS=env=prod,region=us-east-1
}
```

## Struct Tags

| Tag | Description |
|-----|-------------|
| `env:"VAR_NAME"` | Override env var name for this field |
| `env:"VAR_NAME,required"` | Mark field as required (error if unset and no default) |
| `env:"-"` | Skip field entirely |
| `envconfig:"VAR_NAME"` | Alias for `env` — envconfig compatibility |
| `default:"value"` | Default value if env var is unset |
| `envDefault:"value"` | Alias for `default` — envconfig compatibility |
| `envSeparator:":"` | Custom separator for slice fields (default: `,`) |
| `envExpand:"true"` | Expand `$VAR` references in the value via `os.ExpandEnv` |
| `desc:"description"` | Field description shown in `Usage()` output |

### `envExpand` example

```go
type Config struct {
    DataDir string `default:"$HOME/.myapp" envExpand:"true"`
}
// DataDir will expand to /home/user/.myapp (or whatever $HOME is set to)
```

## Nested Structs

Named nested struct fields are flattened with `_` separators:

```go
type DB struct {
    Host string
    Port int
}
type Config struct {
    Database DB
}
// envstruct.Process("APP", &cfg) reads:
//   APP_DATABASE_HOST
//   APP_DATABASE_PORT
```

Embedded (anonymous) structs are flattened without adding a prefix segment:

```go
type Base struct {
    LogLevel string
}
type Config struct {
    Base                // embedded — no extra segment added
    Port int
}
// envstruct.Process("APP", &cfg) reads:
//   APP_LOG_LEVEL
//   APP_PORT
```

### Pointer-to-struct

Pointer-to-struct fields are automatically allocated only when at least one nested env var is actually set. If no env var matches, the pointer remains `nil`.

```go
type TLSConfig struct {
    Cert string
    Key  string
}
type Config struct {
    TLS *TLSConfig
}
// If neither APP_TLS_CERT nor APP_TLS_KEY is set, cfg.TLS == nil.
// If either is set, cfg.TLS is allocated and populated.
```

## Usage Text

```go
envstruct.Usage("APP", &Config{}, os.Stderr)
```

Outputs a table of all expected environment variables with their type, default value, and description. Useful to print in response to a `--help` flag or on startup error.

## Migration from envconfig

envstruct is a drop-in replacement. Function signatures are identical:

```go
// Before
import "github.com/kelseyhightower/envconfig"
envconfig.Process("APP", &config)
envconfig.MustProcess("APP", &config)

// After
import "github.com/agentine/envstruct"
envstruct.Process("APP", &config)
envstruct.MustProcess("APP", &config)
```

Both `env` and `envconfig` struct tags are supported for smooth migration. Both `default` and `envDefault` tags are also recognized.

If your types implement `envconfig.Setter` (a `Set(string) error` method), they will work without modification.

## License

MIT
