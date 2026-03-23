package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestXMLFormatter_NoColor(t *testing.T) {
	input := `<root><child id="1"><name>Alice</name></child></root>`

	f := &XMLFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "  <child") {
		t.Error("expected indented child element")
	}
	if !strings.Contains(out, "    <name>") || !strings.Contains(out, "Alice") {
		t.Error("expected indented name element with content")
	}
}
