package detect

import "testing"

func TestJWTDetector(t *testing.T) {
	d := &JWTDetector{}

	tests := []struct {
		name     string
		input    string
		wantConf Confidence
	}{
		{
			"valid jwt",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			High,
		},
		{
			"not three segments",
			"eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0",
			None,
		},
		{
			"plain text",
			"just some text",
			None,
		},
		{
			"multi-line",
			"line1\nline2",
			None,
		},
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
