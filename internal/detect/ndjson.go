package detect

import (
	"bytes"
	"encoding/json"
)

// NDJSONDetector detects newline-delimited JSON (one JSON object per line).
type NDJSONDetector struct{}

func (d *NDJSONDetector) Detect(sample []byte) Result {
	lines := bytes.Split(sample, []byte("\n"))

	var nonEmpty [][]byte
	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) > 0 {
			nonEmpty = append(nonEmpty, trimmed)
		}
	}

	// Require at least 2 non-empty lines; a single JSON object stays as JSON.
	if len(nonEmpty) < 2 {
		return Result{Format: NDJSON, Confidence: None}
	}

	// Each non-empty line must start with '{'.
	for _, line := range nonEmpty {
		if line[0] != '{' {
			return Result{Format: NDJSON, Confidence: None}
		}
	}

	// Count lines that are valid JSON.
	valid := 0
	for _, line := range nonEmpty {
		if json.Valid(line) {
			valid++
		}
	}

	if valid == len(nonEmpty) {
		return Result{Format: NDJSON, Confidence: High}
	}

	// All start with '{' but some fail parse (possibly truncated sniff buffer).
	return Result{Format: NDJSON, Confidence: Medium}
}
