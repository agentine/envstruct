package envstruct

import (
	"bytes"
	"testing"
)

func TestUsageNoError(t *testing.T) {
	type Config struct {
		Host string
	}
	var buf bytes.Buffer
	if err := Usage("APP", &Config{}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
