package envstruct

import "testing"

func TestDecoderInterface(t *testing.T) {
	// Verify Decoder and Setter interfaces are defined.
	var _ Decoder = (*testDecoder)(nil)
	var _ Setter = (*testSetter)(nil)
}

type testDecoder struct{ val string }

func (d *testDecoder) Decode(value string) error { d.val = value; return nil }

type testSetter struct{ val string }

func (s *testSetter) Set(value string) error { s.val = value; return nil }
