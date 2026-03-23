package detect

import "testing"

func TestTOMLDetector(t *testing.T) {
	d := &TOMLDetector{}

	tests := []struct {
		name     string
		input    string
		wantConf Confidence
	}{
		{"section + kv", "[server]\nhost = \"localhost\"\nport = 8080", High},
		{"multiple kv", "name = \"pp\"\nversion = \"0.1\"\nauthor = \"test\"", Medium},
		{"single kv", "name = \"pp\"", None},
		{"yaml not toml", "name: pp\nversion: 0.1", None},
		{"json", `{"key": "value"}`, None},
		{"empty", "", None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := d.Detect([]byte(tt.input))
			if r.Confidence != tt.wantConf {
				t.Errorf("confidence for %q = %v, want %v", tt.name, r.Confidence, tt.wantConf)
			}
		})
	}
}
