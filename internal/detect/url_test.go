package detect

import "testing"

func TestURLDetector(t *testing.T) {
	d := &URLDetector{}

	tests := []struct {
		name     string
		input    string
		wantConf Confidence
	}{
		{"3+ pairs", "name=Alice&age=30&city=NYC", High},
		{"2 pairs", "name=Alice&age=30", Medium},
		{"1 pair no ampersand", "name=Alice", None},
		{"multi-line", "name=Alice\nage=30", None},
		{"no equals", "just&some&text", None},
		{"empty", "", None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := d.Detect([]byte(tt.input))
			if r.Confidence != tt.wantConf {
				t.Errorf("confidence = %v, want %v", r.Confidence, tt.wantConf)
			}
		})
	}
}
