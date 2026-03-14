package envstruct

import (
	"errors"
	"net/url"
	"testing"
	"time"
)

func TestDecoderInterface(t *testing.T) {
	var _ Decoder = (*testDecoder)(nil)
	var _ Setter = (*testSetter)(nil)
}

type testDecoder struct{ val string }

func (d *testDecoder) Decode(value string) error { d.val = value; return nil }

type testSetter struct{ val string }

func (s *testSetter) Set(value string) error { s.val = value; return nil }

func TestDecodeString(t *testing.T) {
	type C struct{ Name string }
	t.Setenv("NAME", "hello")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Name != "hello" {
		t.Fatalf("got %q", c.Name)
	}
}

func TestDecodeBool(t *testing.T) {
	type C struct{ Debug bool }
	t.Setenv("DEBUG", "true")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if !c.Debug {
		t.Fatal("expected true")
	}
}

func TestDecodeBoolFalse(t *testing.T) {
	type C struct{ Debug bool }
	t.Setenv("DEBUG", "false")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Debug {
		t.Fatal("expected false")
	}
}

func TestDecodeInt(t *testing.T) {
	type C struct{ Port int }
	t.Setenv("PORT", "8080")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Port != 8080 {
		t.Fatalf("got %d", c.Port)
	}
}

func TestDecodeInt8(t *testing.T) {
	type C struct{ Val int8 }
	t.Setenv("VAL", "127")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Val != 127 {
		t.Fatalf("got %d", c.Val)
	}
}

func TestDecodeInt16(t *testing.T) {
	type C struct{ Val int16 }
	t.Setenv("VAL", "32000")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Val != 32000 {
		t.Fatalf("got %d", c.Val)
	}
}

func TestDecodeInt32(t *testing.T) {
	type C struct{ Val int32 }
	t.Setenv("VAL", "2000000")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Val != 2000000 {
		t.Fatalf("got %d", c.Val)
	}
}

func TestDecodeInt64(t *testing.T) {
	type C struct{ Val int64 }
	t.Setenv("VAL", "9000000000")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Val != 9000000000 {
		t.Fatalf("got %d", c.Val)
	}
}

func TestDecodeUint(t *testing.T) {
	type C struct{ Val uint }
	t.Setenv("VAL", "42")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Val != 42 {
		t.Fatalf("got %d", c.Val)
	}
}

func TestDecodeUint8(t *testing.T) {
	type C struct{ Val uint8 }
	t.Setenv("VAL", "255")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Val != 255 {
		t.Fatalf("got %d", c.Val)
	}
}

func TestDecodeUint16(t *testing.T) {
	type C struct{ Val uint16 }
	t.Setenv("VAL", "60000")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Val != 60000 {
		t.Fatalf("got %d", c.Val)
	}
}

func TestDecodeUint32(t *testing.T) {
	type C struct{ Val uint32 }
	t.Setenv("VAL", "4000000000")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Val != 4000000000 {
		t.Fatalf("got %d", c.Val)
	}
}

func TestDecodeUint64(t *testing.T) {
	type C struct{ Val uint64 }
	t.Setenv("VAL", "18000000000000000000")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Val != 18000000000000000000 {
		t.Fatalf("got %d", c.Val)
	}
}

func TestDecodeFloat32(t *testing.T) {
	type C struct{ Val float32 }
	t.Setenv("VAL", "3.14")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Val < 3.13 || c.Val > 3.15 {
		t.Fatalf("got %f", c.Val)
	}
}

func TestDecodeFloat64(t *testing.T) {
	type C struct{ Val float64 }
	t.Setenv("VAL", "3.141592653589793")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Val != 3.141592653589793 {
		t.Fatalf("got %f", c.Val)
	}
}

func TestDecodeDuration(t *testing.T) {
	type C struct{ Timeout time.Duration }
	t.Setenv("TIMEOUT", "5s")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Timeout != 5*time.Second {
		t.Fatalf("got %v", c.Timeout)
	}
}

func TestDecodeURL(t *testing.T) {
	type C struct{ Endpoint url.URL }
	t.Setenv("ENDPOINT", "https://example.com/path")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Endpoint.String() != "https://example.com/path" {
		t.Fatalf("got %v", c.Endpoint.String())
	}
}

func TestDecodeURLPointer(t *testing.T) {
	type C struct{ Endpoint *url.URL }
	t.Setenv("ENDPOINT", "https://example.com")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Endpoint == nil || c.Endpoint.Host != "example.com" {
		t.Fatalf("got %v", c.Endpoint)
	}
}

func TestDecodeSliceString(t *testing.T) {
	type C struct{ Hosts []string }
	t.Setenv("HOSTS", "a,b,c")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if len(c.Hosts) != 3 || c.Hosts[0] != "a" || c.Hosts[1] != "b" || c.Hosts[2] != "c" {
		t.Fatalf("got %v", c.Hosts)
	}
}

func TestDecodeSliceInt(t *testing.T) {
	type C struct{ Ports []int }
	t.Setenv("PORTS", "80,443,8080")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if len(c.Ports) != 3 || c.Ports[0] != 80 || c.Ports[1] != 443 || c.Ports[2] != 8080 {
		t.Fatalf("got %v", c.Ports)
	}
}

func TestDecodeMapStringString(t *testing.T) {
	type C struct{ Labels map[string]string }
	t.Setenv("LABELS", "env=prod,region=us")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Labels["env"] != "prod" || c.Labels["region"] != "us" {
		t.Fatalf("got %v", c.Labels)
	}
}

func TestDecodeMapStringInt(t *testing.T) {
	type C struct{ Limits map[string]int }
	t.Setenv("LIMITS", "cpu=4,mem=8192")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Limits["cpu"] != 4 || c.Limits["mem"] != 8192 {
		t.Fatalf("got %v", c.Limits)
	}
}

func TestDecodeCustomDecoder(t *testing.T) {
	type C struct{ Dec testDecoder }
	t.Setenv("DEC", "custom-value")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Dec.val != "custom-value" {
		t.Fatalf("got %q", c.Dec.val)
	}
}

func TestDecodeCustomSetter(t *testing.T) {
	type C struct{ Set testSetter }
	t.Setenv("SET", "setter-value")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Set.val != "setter-value" {
		t.Fatalf("got %q", c.Set.val)
	}
}

func TestDecodePointerInt(t *testing.T) {
	type C struct{ Count *int }
	t.Setenv("COUNT", "42")
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Count == nil || *c.Count != 42 {
		t.Fatalf("got %v", c.Count)
	}
}

func TestDecodePointerNilWhenUnset(t *testing.T) {
	type C struct{ Count *int }
	var c C
	if err := Process("", &c); err != nil {
		t.Fatal(err)
	}
	if c.Count != nil {
		t.Fatalf("expected nil, got %v", c.Count)
	}
}

func TestDecodeParseError(t *testing.T) {
	type C struct{ Port int }
	t.Setenv("PORT", "not-a-number")
	var c C
	err := Process("", &c)
	if err == nil {
		t.Fatal("expected error")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected ParseError, got %T", err)
	}
}

func TestDecodeBoolParseError(t *testing.T) {
	type C struct{ Debug bool }
	t.Setenv("DEBUG", "nope")
	var c C
	err := Process("", &c)
	if err == nil {
		t.Fatal("expected error")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected ParseError, got %T", err)
	}
}
