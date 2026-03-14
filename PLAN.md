# envstruct

**Replacement for:** [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig)
**Language:** Go
**Package:** `github.com/agentine/envstruct`

## Why

kelseyhightower/envconfig has 8,524 importers and 20.2K GitHub dependents but has been dormant since May 2019 (last release v1.4.0). The sole maintainer is no longer active in OSS. There are 27 open issues and 31 open PRs with zero maintainer engagement. Existing alternatives (caarlos0/env, sethvargo/go-envconfig) have not captured envconfig's market share.

## Scope

A Go library that populates struct fields from environment variables using struct tags. Drop-in replacement for envconfig with modern Go practices.

## Core Features

1. **Struct tag parsing** — `env:"VAR_NAME"` tags to map env vars to struct fields
2. **Type support** — string, int (all sizes), uint, float, bool, time.Duration, url.URL, custom types
3. **Prefix support** — scoped env var lookup via configurable prefix (e.g., `APP_`)
4. **Nested structs** — flatten or use delimiter-separated names (e.g., `APP_DB_HOST`)
5. **Required fields** — `env:"VAR_NAME,required"` tag option
6. **Default values** — `env:"VAR_NAME" default:"value"` tag
7. **Custom decoders** — `Decoder` interface for user-defined types
8. **Slice/map support** — comma-separated values for slices, key=value pairs for maps
9. **Usage generation** — auto-generate usage text from struct tags
10. **Error reporting** — clear error messages with field name, expected type, and env var name

## Architecture

```
envstruct/
├── envstruct.go        # Core Process/MustProcess functions
├── decoder.go          # Type decoders and Decoder interface
├── tags.go             # Struct tag parsing
├── usage.go            # Usage text generation
├── errors.go           # Typed error types
├── envstruct_test.go   # Core tests
├── decoder_test.go     # Decoder tests
├── usage_test.go       # Usage tests
├── example_test.go     # Runnable examples
├── go.mod
├── go.sum
└── README.md
```

## API Surface

```go
// Process populates a struct from environment variables.
func Process(prefix string, spec interface{}) error

// MustProcess is like Process but panics on error.
func MustProcess(prefix string, spec interface{})

// Usage writes a usage message to the given writer.
func Usage(prefix string, spec interface{}, out io.Writer) error

// Decoder is implemented by types that can decode themselves from a string.
type Decoder interface {
    Decode(value string) error
}

// Setter is implemented by types that can set themselves from a string (envconfig compat).
type Setter interface {
    Set(value string) error
}
```

## Migration Path

- Same `Process(prefix, &spec)` function signature as envconfig
- Supports envconfig's struct tag format for easy migration
- Also supports `env:"NAME"` tag format (more conventional)
- Decoder interface is compatible with envconfig's interface

## Deliverables

1. Core library with full type support
2. Comprehensive test suite (>90% coverage)
3. README with migration guide from envconfig
4. Runnable examples
5. CI configuration (GitHub Actions)

## Non-Goals

- File-based config (.env, YAML, TOML) — use dedicated libraries
- Config hot-reloading — environment variables are read once at startup
- CLI flag parsing — separate concern
