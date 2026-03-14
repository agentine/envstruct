package envstruct

import (
	"errors"
	"os"
	"strings"
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

// --- Nested struct tests ---

func TestNestedStruct(t *testing.T) {
	type DB struct {
		Host string
		Port int
	}
	type Config struct {
		Database DB
	}
	setEnv(t, "APP_DATABASE_HOST", "db.local")
	setEnv(t, "APP_DATABASE_PORT", "5432")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Database.Host != "db.local" {
		t.Fatalf("expected 'db.local', got %q", c.Database.Host)
	}
	if c.Database.Port != 5432 {
		t.Fatalf("expected 5432, got %d", c.Database.Port)
	}
}

func TestThreeLevelNesting(t *testing.T) {
	type Conn struct {
		Host string
	}
	type DB struct {
		Primary Conn
	}
	type Config struct {
		Database DB
	}
	setEnv(t, "APP_DATABASE_PRIMARY_HOST", "deep.local")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Database.Primary.Host != "deep.local" {
		t.Fatalf("expected 'deep.local', got %q", c.Database.Primary.Host)
	}
}

func TestPointerToStruct(t *testing.T) {
	type DB struct {
		Host string
	}
	type Config struct {
		Database *DB
	}
	setEnv(t, "APP_DATABASE_HOST", "ptr-db.local")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Database == nil || c.Database.Host != "ptr-db.local" {
		t.Fatalf("expected 'ptr-db.local', got %v", c.Database)
	}
}

func TestEmbeddedStruct(t *testing.T) {
	type Common struct {
		Debug bool
	}
	type Config struct {
		Common
		Host string
	}
	setEnv(t, "APP_HOST", "embedded.local")
	setEnv(t, "APP_DEBUG", "true")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Host != "embedded.local" {
		t.Fatalf("expected 'embedded.local', got %q", c.Host)
	}
	if !c.Debug {
		t.Fatal("expected Debug=true")
	}
}

func TestNestedCustomPrefix(t *testing.T) {
	type DB struct {
		Host string
	}
	type Config struct {
		Database DB `env:"DB"`
	}
	setEnv(t, "APP_DB_HOST", "custom.local")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Database.Host != "custom.local" {
		t.Fatalf("expected 'custom.local', got %q", c.Database.Host)
	}
}

func TestMultipleErrorsCollected(t *testing.T) {
	type Config struct {
		Host string `env:"HOST,required"`
		Port int    `env:"PORT,required"`
	}
	var c Config
	err := Process("APP", &c)
	if err == nil {
		t.Fatal("expected error")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "HOST") || !strings.Contains(errStr, "PORT") {
		t.Fatalf("expected both HOST and PORT errors, got: %s", errStr)
	}
}

// --- Pointer-to-struct nil when no env vars set (Task #428) ---

func TestPointerToStructNilWhenNoEnvVars(t *testing.T) {
	type DB struct {
		Host string
		Port int
	}
	type Config struct {
		Database *DB
	}
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Database != nil {
		t.Fatalf("expected nil pointer when no env vars set, got %+v", c.Database)
	}
}

func TestPointerToStructSetWhenPartialEnvVars(t *testing.T) {
	type DB struct {
		Host string
		Port int
	}
	type Config struct {
		Database *DB
	}
	setEnv(t, "APP_DATABASE_HOST", "partial.local")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Database == nil {
		t.Fatal("expected non-nil pointer when env var is set")
	}
	if c.Database.Host != "partial.local" {
		t.Fatalf("expected 'partial.local', got %q", c.Database.Host)
	}
}

// --- CamelCase to UPPER_SNAKE_CASE (Task #429) ---

func TestCamelCaseFieldName(t *testing.T) {
	type Config struct {
		DatabaseURL string
	}
	setEnv(t, "APP_DATABASE_URL", "postgres://localhost")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.DatabaseURL != "postgres://localhost" {
		t.Fatalf("expected 'postgres://localhost', got %q", c.DatabaseURL)
	}
}

func TestCamelCaseMultiWord(t *testing.T) {
	type Config struct {
		MaxRetryCount int
	}
	setEnv(t, "MAX_RETRY_COUNT", "5")
	var c Config
	if err := Process("", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.MaxRetryCount != 5 {
		t.Fatalf("expected 5, got %d", c.MaxRetryCount)
	}
}

// --- envDefault tag (Task #429) ---

func TestEnvDefaultTag(t *testing.T) {
	type Config struct {
		Port int `envDefault:"5000"`
	}
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Port != 5000 {
		t.Fatalf("expected 5000, got %d", c.Port)
	}
}

func TestEnvDefaultOverriddenByEnv(t *testing.T) {
	type Config struct {
		Port int `envDefault:"5000"`
	}
	setEnv(t, "APP_PORT", "9090")
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Port != 9090 {
		t.Fatalf("expected 9090, got %d", c.Port)
	}
}

func TestDefaultTagTakesPrecedenceOverEnvDefault(t *testing.T) {
	type Config struct {
		Port int `default:"3000" envDefault:"5000"`
	}
	var c Config
	if err := Process("APP", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Port != 3000 {
		t.Fatalf("expected 3000 (default takes precedence), got %d", c.Port)
	}
}

// --- envSeparator tag (Task #429) ---

func TestEnvSeparatorTag(t *testing.T) {
	type Config struct {
		Paths []string `envSeparator:":"`
	}
	setEnv(t, "PATHS", "/usr/bin:/usr/local/bin:/home/user/bin")
	var c Config
	if err := Process("", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Paths) != 3 || c.Paths[0] != "/usr/bin" || c.Paths[2] != "/home/user/bin" {
		t.Fatalf("expected 3 paths, got %v", c.Paths)
	}
}

func TestEnvSeparatorDefaultComma(t *testing.T) {
	type Config struct {
		Tags []string
	}
	setEnv(t, "TAGS", "a,b,c")
	var c Config
	if err := Process("", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %v", c.Tags)
	}
}

// --- envExpand tag (Task #429) ---

func TestEnvExpandTag(t *testing.T) {
	type Config struct {
		DataDir string `envExpand:"true"`
	}
	setEnv(t, "HOME", "/home/testuser")
	setEnv(t, "DATA_DIR", "$HOME/data")
	var c Config
	if err := Process("", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.DataDir != "/home/testuser/data" {
		t.Fatalf("expected '/home/testuser/data', got %q", c.DataDir)
	}
}

func TestEnvExpandDisabledByDefault(t *testing.T) {
	type Config struct {
		Val string
	}
	setEnv(t, "VAL", "$HOME/data")
	var c Config
	if err := Process("", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Val != "$HOME/data" {
		t.Fatalf("expected literal '$HOME/data', got %q", c.Val)
	}
}

// --- Empty slice decode ---

func TestDecodeEmptySlice(t *testing.T) {
	type Config struct {
		Hosts []string
	}
	setEnv(t, "HOSTS", "")
	var c Config
	if err := Process("", &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Hosts) != 0 {
		t.Fatalf("expected empty slice, got %v", c.Hosts)
	}
}
