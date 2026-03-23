package detect

import "testing"

func TestXMLDetector(t *testing.T) {
	d := &XMLDetector{}

	tests := []struct {
		name     string
		input    string
		wantFmt  Format
		wantConf Confidence
	}{
		{"xml declaration", `<?xml version="1.0"?><root/>`, XML, High},
		{"html doctype", `<!DOCTYPE html><html></html>`, HTML, High},
		{"html tag", `<html lang="en"><body></body></html>`, HTML, High},
		{"xml tag", `<root><child/></root>`, XML, Medium},
		{"json", `{"key": "value"}`, XML, None},
		{"plain text", "hello world", XML, None},
		{"empty", "", XML, None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := d.Detect([]byte(tt.input))
			if r.Confidence != tt.wantConf {
				t.Errorf("confidence = %v, want %v", r.Confidence, tt.wantConf)
			}
			if r.Confidence > None && r.Format != tt.wantFmt {
				t.Errorf("format = %v, want %v", r.Format, tt.wantFmt)
			}
		})
	}
}
