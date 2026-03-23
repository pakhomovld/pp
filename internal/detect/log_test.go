package detect

import "testing"

func TestLogDetector(t *testing.T) {
	d := &LogDetector{}

	tests := []struct {
		name     string
		input    string
		wantConf Confidence
	}{
		{
			"iso8601 + levels",
			"2024-01-15T10:30:00Z INFO Starting\n2024-01-15T10:30:01Z ERROR Failed\n2024-01-15T10:30:02Z WARN Retry",
			High,
		},
		{
			"bracketed levels",
			"[INFO] Starting server\n[ERROR] Connection failed\n[WARN] Retrying",
			High,
		},
		{
			"two log lines",
			"2024-01-15 10:30:00 INFO Starting\n2024-01-15 10:30:01 ERROR Failed",
			Medium,
		},
		{
			"single log line",
			"2024-01-15T10:30:00Z INFO Starting server",
			Low,
		},
		{
			"plain text",
			"just some text\nwith multiple lines\nnothing special",
			None,
		},
		{"json", `{"level": "info"}`, None},
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
