package detect

import (
	"encoding/base64"
	"testing"
)

func TestBase64Detector(t *testing.T) {
	d := &Base64Detector{}

	tests := []struct {
		name     string
		input    string
		wantConf Confidence
	}{
		{
			"valid base64 of utf8",
			base64.StdEncoding.EncodeToString([]byte(`{"hello": "world"}`)),
			Low,
		},
		{
			"too short",
			"aGVsbG8=",
			None,
		},
		{
			"not base64 chars",
			"this is definitely not base64!!!",
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
