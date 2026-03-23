package detect

import (
	"bytes"
	"encoding/base64"
	"strings"
)

// JWTDetector detects JSON Web Tokens (three base64url-encoded segments separated by dots).
type JWTDetector struct{}

func (d *JWTDetector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) == 0 {
		return Result{Format: JWT, Confidence: None}
	}

	// JWT is a single line.
	if bytes.ContainsAny(trimmed, "\n\r") {
		return Result{Format: JWT, Confidence: None}
	}

	s := string(trimmed)
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return Result{Format: JWT, Confidence: None}
	}

	// All three segments must be valid base64url.
	for _, part := range parts {
		if len(part) == 0 {
			return Result{Format: JWT, Confidence: None}
		}
		if !isBase64URL(part) {
			return Result{Format: JWT, Confidence: None}
		}
	}

	// First segment (header) should decode to JSON with "alg" field.
	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Result{Format: JWT, Confidence: None}
	}
	if bytes.Contains(header, []byte(`"alg"`)) {
		return Result{Format: JWT, Confidence: High}
	}

	return Result{Format: JWT, Confidence: Medium}
}

func isBase64URL(s string) bool {
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '_' || c == '=') {
			return false
		}
	}
	return true
}
