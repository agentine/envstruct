package envstruct

import "fmt"

// ParseError is returned when an environment variable value cannot be
// parsed into the target type.
type ParseError struct {
	FieldName string
	EnvVar    string
	Value     string
	TypeName  string
	Err       error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("envstruct: parsing %q for field %s (env %s): %v",
		e.Value, e.FieldName, e.EnvVar, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// RequiredError is returned when a required environment variable is missing.
type RequiredError struct {
	FieldName string
	EnvVar    string
}

func (e *RequiredError) Error() string {
	return fmt.Sprintf("envstruct: required environment variable %s (field %s) is not set",
		e.EnvVar, e.FieldName)
}
