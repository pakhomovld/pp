package detect

import (
	"bytes"
	"encoding/base64"
	"unicode/utf8"
)

// Base64Detector detects base64-encoded data.
// Low priority — only fires if nothing else matched, since many strings are
// coincidentally valid base64.
type Base64Detector struct{}

func (d *Base64Detector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) < 20 {
		return Result{Format: Base64, Confidence: None}
	}

	// Must be single-line or have no whitespace (except padding lines in PEM).
	line := trimmed
	if bytes.ContainsAny(trimmed, "\n\r") {
		// Multi-line base64 (PEM-style): join lines and check.
		line = bytes.ReplaceAll(trimmed, []byte("\n"), nil)
		line = bytes.ReplaceAll(line, []byte("\r"), nil)
	}

	// Check character set.
	for _, b := range line {
		if !isBase64Char(b) {
			return Result{Format: Base64, Confidence: None}
		}
	}

	// Try to decode.
	decoded, err := base64.StdEncoding.DecodeString(string(line))
	if err != nil {
		// Try without padding.
		decoded, err = base64.RawStdEncoding.DecodeString(string(line))
		if err != nil {
			return Result{Format: Base64, Confidence: None}
		}
	}

	// Decoded content must be valid UTF-8 to be useful to display.
	if !utf8.Valid(decoded) {
		return Result{Format: Base64, Confidence: None}
	}

	return Result{Format: Base64, Confidence: Low}
}

func isBase64Char(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') || b == '+' || b == '/' || b == '='
}
