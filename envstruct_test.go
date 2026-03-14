package envstruct

import (
	"errors"
	"os"
	"testing"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func TestProcessBasicString(t *testing.T) {
	type Config struct {
		Host string
	}
	setEnv(t, "APP_HOST", "localhost")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Host != "localhost" {
		t.Fatalf("expected 'localhost', got %q", c.Host)
	}
}

func TestProcessNoPrefix(t *testing.T) {
	type Config struct {
		Host string
	}
	setEnv(t, "HOST", "example.com")
	var c Config
	if err := Process("", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Host != "example.com" {
		t.Fatalf("expected 'example.com', got %q", c.Host)
	}
}

func TestProcessEnvTag(t *testing.T) {
	type Config struct {
		Host string `env:"CUSTOM_HOST"`
	}
	setEnv(t, "APP_CUSTOM_HOST", "tagged.com")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Host != "tagged.com" {
		t.Fatalf("expected 'tagged.com', got %q", c.Host)
	}
}

func TestProcessEnvconfigTag(t *testing.T) {
	type Config struct {
		Host string `envconfig:"MY_HOST"`
	}
	setEnv(t, "APP_MY_HOST", "compat.com")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Host != "compat.com" {
		t.Fatalf("expected 'compat.com', got %q", c.Host)
	}
}

func TestProcessRequired(t *testing.T) {
	type Config struct {
		Host string `env:"HOST,required"`
	}
	var c Config
	err := Process("APP", &c)
	if err == nil {
		t.Fatal("expected error for missing required field")
	}
	var reqErr *RequiredError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected RequiredError, got %T: %v", err, err)
	}
	if reqErr.FieldName != "Host" {
		t.Fatalf("expected field 'Host', got %q", reqErr.FieldName)
	}
}

func TestProcessRequiredPresent(t *testing.T) {
	type Config struct {
		Host string `env:"HOST,required"`
	}
	setEnv(t, "APP_HOST", "present.com")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Host != "present.com" {
		t.Fatalf("expected 'present.com', got %q", c.Host)
	}
}

func TestProcessDefault(t *testing.T) {
	type Config struct {
		Host string `default:"default.com"`
	}
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Host != "default.com" {
		t.Fatalf("expected 'default.com', got %q", c.Host)
	}
}

func TestProcessDefaultOverriddenByEnv(t *testing.T) {
	type Config struct {
		Host string `default:"default.com"`
	}
	setEnv(t, "APP_HOST", "override.com")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Host != "override.com" {
		t.Fatalf("expected 'override.com', got %q", c.Host)
	}
}

func TestProcessIgnored(t *testing.T) {
	type Config struct {
		Host   string
		Secret string `env:"-"`
	}
	setEnv(t, "APP_HOST", "host.com")
	os.Setenv("APP_SECRET", "should-be-ignored")
	defer os.Unsetenv("APP_SECRET")

	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Secret != "" {
		t.Fatalf("expected empty string for ignored field, got %q", c.Secret)
	}
}

func TestProcessUnexportedFieldSkipped(t *testing.T) {
	type Config struct {
		Host   string
		secret string //nolint:unused
	}
	setEnv(t, "APP_HOST", "host.com")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Host != "host.com" {
		t.Fatalf("expected 'host.com', got %q", c.Host)
	}
}

func TestProcessNonPointer(t *testing.T) {
	type Config struct {
		Host string
	}
	var c Config
	err := Process("APP", c)
	if err == nil {
		t.Fatal("expected error for non-pointer")
	}
}

func TestProcessNilPointer(t *testing.T) {
	err := Process("APP", (*struct{})(nil))
	if err == nil {
		t.Fatal("expected error for nil pointer")
	}
}

func TestProcessPointerField(t *testing.T) {
	type Config struct {
		Host *string
	}
	setEnv(t, "APP_HOST", "ptr.com")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Host == nil || *c.Host != "ptr.com" {
		t.Fatalf("expected pointer to 'ptr.com', got %v", c.Host)
	}
}

func TestProcessPointerFieldUnset(t *testing.T) {
	type Config struct {
		Host *string
	}
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Host != nil {
		t.Fatalf("expected nil pointer, got %v", c.Host)
	}
}

func TestMustProcessPanics(t *testing.T) {
	type Config struct {
		Host string `env:"HOST,required"`
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	var c Config
	MustProcess("APP", &c)
}
