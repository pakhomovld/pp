package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestSQLFormatter_SimpleSelect(t *testing.T) {
	input := "select id, name from users where active = true order by name"

	f := &SQLFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	// Keywords should be uppercased.
	if !strings.Contains(got, "SELECT") {
		t.Error("expected SELECT to be uppercased")
	}
	if !strings.Contains(got, "FROM") {
		t.Error("expected FROM to be uppercased")
	}
	if !strings.Contains(got, "WHERE") {
		t.Error("expected WHERE to be uppercased")
	}

	// FROM and WHERE should be on new lines.
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines, got %d: %q", len(lines), got)
	}
}

func TestSQLFormatter_Subquery(t *testing.T) {
	input := "SELECT * FROM (SELECT id FROM users WHERE active = true) AS sub"

	f := &SQLFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	// Should contain parentheses.
	if !strings.Contains(got, "(") || !strings.Contains(got, ")") {
		t.Error("expected parentheses in output")
	}
}

func TestSQLFormatter_Insert(t *testing.T) {
	input := "INSERT INTO users VALUES (1, 'Alice')"

	f := &SQLFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	if !strings.Contains(got, "INSERT INTO") {
		t.Error("expected INSERT INTO")
	}
	if !strings.Contains(got, "VALUES") {
		t.Error("expected VALUES")
	}
}

func TestSQLFormatter_CommentPreserved(t *testing.T) {
	input := "SELECT * -- get all\nFROM users"

	f := &SQLFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	if !strings.Contains(got, "-- get all") {
		t.Errorf("expected comment to be preserved, got: %q", got)
	}
}

func TestSQLFormatter_NoColor(t *testing.T) {
	input := "SELECT id FROM users WHERE id = 1"

	f := &SQLFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	// Should still be reformatted (keywords uppercased, newlines added).
	if !strings.Contains(got, "SELECT") {
		t.Error("expected SELECT uppercased even without color")
	}
	if !strings.Contains(got, "\n") {
		t.Error("expected newlines in formatted output")
	}
}
