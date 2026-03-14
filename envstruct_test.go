package envstruct

import "testing"

func TestProcessNil(t *testing.T) {
	type Config struct {
		Host string
	}
	var c Config
	if err := Process("", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
