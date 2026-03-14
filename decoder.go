package envstruct

// Decoder is implemented by types that can decode themselves from a string.
type Decoder interface {
	Decode(value string) error
}

// Setter is implemented by types that can set themselves from a string.
// This interface exists for compatibility with kelseyhightower/envconfig.
type Setter interface {
	Set(value string) error
}
