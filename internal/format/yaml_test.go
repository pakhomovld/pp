package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestYAMLFormatter_NoColor(t *testing.T) {
	input := "---\nname: test\ncount: 42\nenabled: true\n# comment\nitems:\n  - one\n  - two\n"

	f := &YAMLFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Without color, YAML passes through unchanged.
	if buf.String() != input {
		t.Errorf("no-color YAML should pass through unchanged\ngot:\n%s\nwant:\n%s", buf.String(), input)
	}
}
