package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestTOMLFormatter_NoColor(t *testing.T) {
	input := "[server]\nhost = \"localhost\"\nport = 8080\n"

	f := &TOMLFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	// No color = passthrough.
	if buf.String() != input {
		t.Errorf("no-color TOML should pass through\ngot: %q\nwant: %q", buf.String(), input)
	}
}
