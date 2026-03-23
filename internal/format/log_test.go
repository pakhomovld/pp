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

	// No color = passthrough (no embedded JSON).
	if buf.String() != input {
		t.Errorf("no-color logs should pass through\ngot: %q\nwant: %q", buf.String(), input)
	}
}

func TestLogFormatter_EmbeddedJSON(t *testing.T) {
	f := &LogFormatter{}

	tests := []struct {
		name        string
		input       string
		wantContain []string
		wantMissing []string
	}{
		{
			"json object in log line",
			"2024-01-15T10:30:00Z INFO {\"status\":\"ok\",\"code\":200}\n",
			[]string{"2024-01-15T10:30:00Z INFO\n", "\"status\": \"ok\"", "\"code\": 200"},
			nil,
		},
		{
			"no json passes through",
			"2024-01-15T10:30:00Z INFO Starting server\n",
			[]string{"2024-01-15T10:30:00Z INFO Starting server"},
			nil,
		},
		{
			"invalid json brace unchanged",
			"2024-01-15T10:30:00Z INFO {not json\n",
			[]string{"{not json"},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := f.Format(&buf, strings.NewReader(tt.input), nil)
			if err != nil {
				t.Fatal(err)
			}
			got := buf.String()
			for _, s := range tt.wantContain {
				if !strings.Contains(got, s) {
					t.Errorf("output should contain %q\ngot:\n%s", s, got)
				}
			}
			for _, s := range tt.wantMissing {
				if strings.Contains(got, s) {
					t.Errorf("output should not contain %q\ngot:\n%s", s, got)
				}
			}
		})
	}
}

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		line      string
		wantFound bool
		wantPfx   string
	}{
		{"INFO {\"a\":1}", true, "INFO "},
		{"no json here", false, ""},
		{"{\"a\":1}", true, ""},
		{"INFO {broken", false, ""},
		{"INFO [1,2,3]", true, "INFO "},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			prefix, _, found := extractJSON(tt.line)
			if found != tt.wantFound {
				t.Errorf("extractJSON(%q) found=%v, want %v", tt.line, found, tt.wantFound)
			}
			if found && prefix != tt.wantPfx {
				t.Errorf("extractJSON(%q) prefix=%q, want %q", tt.line, prefix, tt.wantPfx)
			}
		})
	}
}
