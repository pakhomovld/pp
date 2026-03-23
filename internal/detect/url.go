package detect

import (
	"bytes"
	"strings"
)

// URLDetector detects URL-encoded query strings (key=value&key=value).
type URLDetector struct{}

func (d *URLDetector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) == 0 {
		return Result{Format: URLEncode, Confidence: None}
	}

	// Must be single-line.
	if bytes.ContainsAny(trimmed, "\n\r") {
		return Result{Format: URLEncode, Confidence: None}
	}

	s := string(trimmed)

	// Must contain at least one = and one &.
	if !strings.Contains(s, "=") || !strings.Contains(s, "&") {
		return Result{Format: URLEncode, Confidence: None}
	}

	// Split by & and check each pair has a key=value structure.
	pairs := strings.Split(s, "&")
	validPairs := 0
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 && len(parts[0]) > 0 {
			validPairs++
		}
	}

	if validPairs >= 3 {
		return Result{Format: URLEncode, Confidence: High}
	}
	if validPairs >= 2 {
		return Result{Format: URLEncode, Confidence: Medium}
	}

	return Result{Format: URLEncode, Confidence: None}
}
