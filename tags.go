package envstruct

// fieldSpec holds parsed struct tag information for a single field.
type fieldSpec struct {
	Name         string
	Required     bool
	DefaultValue string
	HasDefault   bool
	Ignored      bool
}
