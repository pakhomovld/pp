package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestURLFormatter_NoColor(t *testing.T) {
	input := "name=Alice&age=30&city=New+York"

	f := &URLFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "name") || !strings.Contains(out, "Alice") {
		t.Error("output should contain key 'name' with value 'Alice'")
	}
	if !strings.Contains(out, "New York") {
		t.Error("output should decode '+' to space in 'New York'")
	}
}
