package envstruct

import (
	"bytes"
	"net/url"
	"strings"
	"testing"
)

func TestUsageBasic(t *testing.T) {
	type Config struct {
		Host string
		Port int
	}
	var buf bytes.Buffer
	if err := Usage("APP", &Config{}, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_HOST") {
		t.Fatalf("expected APP_HOST in output: %s", out)
	}
	if !strings.Contains(out, "APP_PORT") {
		t.Fatalf("expected APP_PORT in output: %s", out)
	}
	if !strings.Contains(out, "string") {
		t.Fatalf("expected 'string' type in output: %s", out)
	}
	if !strings.Contains(out, "int") {
		t.Fatalf("expected 'int' type in output: %s", out)
	}
}

func TestUsageRequired(t *testing.T) {
	type Config struct {
		Host string `env:"HOST,required"`
	}
	var buf bytes.Buffer
	if err := Usage("APP", &Config{}, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "[required]") {
		t.Fatalf("expected [required] in output: %s", out)
	}
}

func TestUsageDefault(t *testing.T) {
	type Config struct {
		Host string `default:"localhost"`
	}
	var buf bytes.Buffer
	if err := Usage("APP", &Config{}, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "[default: localhost]") {
		t.Fatalf("expected [default: localhost] in output: %s", out)
	}
}

func TestUsageDescription(t *testing.T) {
	type Config struct {
		Host string `desc:"Server hostname"`
	}
	var buf bytes.Buffer
	if err := Usage("APP", &Config{}, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Server hostname") {
		t.Fatalf("expected 'Server hostname' in output: %s", out)
	}
}

func TestUsageNested(t *testing.T) {
	type DB struct {
		Host string
	}
	type Config struct {
		Database DB
	}
	var buf bytes.Buffer
	if err := Usage("APP", &Config{}, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_DATABASE_HOST") {
		t.Fatalf("expected APP_DATABASE_HOST in output: %s", out)
	}
}

func TestUsageIgnored(t *testing.T) {
	type Config struct {
		Host   string
		Secret string `env:"-"`
	}
	var buf bytes.Buffer
	if err := Usage("APP", &Config{}, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "SECRET") {
		t.Fatalf("expected no SECRET in output: %s", out)
	}
}

func TestUsageURLType(t *testing.T) {
	type Config struct {
		Endpoint *url.URL `env:"URL,required" desc:"Database URL"`
	}
	var buf bytes.Buffer
	if err := Usage("APP", &Config{}, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "*url.URL") {
		t.Fatalf("expected '*url.URL' in output: %s", out)
	}
	if !strings.Contains(out, "[required]") {
		t.Fatalf("expected [required] in output: %s", out)
	}
}
