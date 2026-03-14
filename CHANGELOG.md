# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2026-03-14

### Added

- `Process(prefix string, v any) error` — parse environment variables into a struct
- `MustProcess(prefix string, v any)` — panics on error
- `Usage(prefix string, v any) error` — print aligned documentation of expected env vars
- Full type support: `string`, `bool`, `int`/`uint` (all sizes), `float32`/`float64`, `time.Duration`, `*url.URL`, slices, maps, pointer fields
- Custom decoder interfaces: `Decoder`, `Setter`, `encoding.TextUnmarshaler`
- Nested struct support with prefix propagation and custom per-field prefix via `env` tag
- Embedded struct flattening (inherits parent prefix)
- Pointer-to-struct fields: allocated only when at least one child env var is set
- Struct tag support: `env`, `envconfig`, `required`, `default`/`envDefault`, `envSeparator`, `envExpand`, `ignored`, `desc`
- CamelCase field name → `UPPER_SNAKE_CASE` env var conversion (envconfig-compatible)
- Multi-error collection: all field errors reported in a single `ProcessError`
- `ParseError` with `Unwrap()` for error chain inspection
- Zero external dependencies

[0.1.0]: https://github.com/agentine/envstruct/releases/tag/v0.1.0
