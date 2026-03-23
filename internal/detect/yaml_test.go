package detect

import "testing"

func TestYAMLDetector(t *testing.T) {
	d := &YAMLDetector{}

	tests := []struct {
		name     string
		input    string
		wantConf Confidence
	}{
		{"document separator", "---\nname: test", High},
		{"key-value pairs", "name: test\nage: 30\ncity: NYC", Medium},
		{"single key", "name: test", Low},
		{"json object", `{"key": "value"}`, None},
		{"xml", `<root><child/></root>`, None},
		{"empty", "", None},
		{"plain text", "hello world", None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := d.Detect([]byte(tt.input))
			if r.Confidence != tt.wantConf {
				t.Errorf("confidence for %q = %v, want %v", tt.input, r.Confidence, tt.wantConf)
			}
		})
	}
}
