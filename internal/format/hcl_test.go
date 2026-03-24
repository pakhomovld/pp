package format

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pakhomovld/ppp/internal/color"
)

func TestHCLFormatter_NoColor(t *testing.T) {
	input := "resource \"aws_instance\" \"web\" {\n  ami = \"ami-123\"\n}\n"

	f := &HCLFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	// No color = passthrough.
	if buf.String() != input {
		t.Errorf("no-color HCL should pass through\ngot: %q\nwant: %q", buf.String(), input)
	}
}

func TestHCLFormatter_WithColor(t *testing.T) {
	input := "resource \"aws_instance\" \"web\" {\n  ami = \"ami-123\"\n}\n"

	f := &HCLFormatter{}
	var buf bytes.Buffer
	theme := color.DefaultTheme()
	err := f.Format(&buf, strings.NewReader(input), theme)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	// With a theme, output should still contain the key content.
	if !strings.Contains(got, "resource") {
		t.Error("expected 'resource' keyword in output")
	}
	if !strings.Contains(got, "ami") {
		t.Error("expected 'ami' key in output")
	}
}

func TestHCLFormatter_Comment(t *testing.T) {
	input := "# This is a comment\nresource \"aws_instance\" \"web\" {\n}\n"

	f := &HCLFormatter{}
	var buf bytes.Buffer
	theme := color.DefaultTheme()
	err := f.Format(&buf, strings.NewReader(input), theme)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	if !strings.Contains(got, "This is a comment") {
		t.Error("expected comment to be preserved")
	}
}

func TestHCLFormatter_SlashSlashComment(t *testing.T) {
	input := "// Another comment\nvariable \"region\" {\n}\n"

	f := &HCLFormatter{}
	var buf bytes.Buffer
	theme := color.DefaultTheme()
	err := f.Format(&buf, strings.NewReader(input), theme)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	if !strings.Contains(got, "Another comment") {
		t.Error("expected // comment to be preserved")
	}
}
