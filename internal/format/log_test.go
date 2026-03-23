package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogFormatter_NoColor(t *testing.T) {
	input := "2024-01-15T10:30:00Z INFO Starting\n2024-01-15T10:30:01Z ERROR Failed\n"

	f := &LogFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	// No color = passthrough.
	if buf.String() != input {
		t.Errorf("no-color logs should pass through\ngot: %q\nwant: %q", buf.String(), input)
	}
}
