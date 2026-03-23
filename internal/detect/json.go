package detect

import (
	"bytes"
	"encoding/json"
)

// JSONDetector detects JSON objects and arrays.
type JSONDetector struct{}

func (d *JSONDetector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) == 0 {
		return Result{Format: JSON, Confidence: None}
	}

	first := trimmed[0]
	if first != '{' && first != '[' {
		return Result{Format: JSON, Confidence: None}
	}

	if json.Valid(sample) {
		return Result{Format: JSON, Confidence: High}
	}

	// First char matches but invalid — could be truncated large JSON.
	return Result{Format: JSON, Confidence: Medium}
}
