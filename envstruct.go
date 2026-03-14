// Package envstruct populates struct fields from environment variables.
//
// It is a drop-in replacement for kelseyhightower/envconfig with the same
// Process(prefix, &spec) function signature.
package envstruct

import "io"

// Process populates the struct pointed to by spec with values from
// environment variables. The prefix is prepended to each field name
// (or tag override) when looking up environment variables.
func Process(prefix string, spec interface{}) error {
	return nil
}

// MustProcess is like Process but panics on error.
func MustProcess(prefix string, spec interface{}) {
	if err := Process(prefix, spec); err != nil {
		panic(err)
	}
}

// Usage writes a usage message describing the environment variables
// expected by spec to the given writer.
func Usage(prefix string, spec interface{}, out io.Writer) error {
	return nil
}
