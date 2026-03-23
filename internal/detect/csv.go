package detect

import (
	"bytes"
	"strings"
)

// CSVDetector detects CSV and TSV data by checking for consistent
// delimiter counts across the first several lines.
type CSVDetector struct{}

func (d *CSVDetector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) == 0 {
		return Result{Format: CSV, Confidence: None}
	}

	// Skip if it looks like JSON or XML.
	if trimmed[0] == '{' || trimmed[0] == '[' || trimmed[0] == '<' {
		return Result{Format: CSV, Confidence: None}
	}

	lines := strings.Split(string(trimmed), "\n")

	// Need at least 2 lines (header + data).
	if len(lines) < 2 {
		return Result{Format: CSV, Confidence: None}
	}

	// Check comma first, then tab.
	if f, c := checkDelimiter(lines, ','); c > None {
		return Result{Format: f, Confidence: c}
	}
	if f, c := checkDelimiter(lines, '\t'); c > None {
		return Result{Format: f, Confidence: c}
	}

	return Result{Format: CSV, Confidence: None}
}

func checkDelimiter(lines []string, delim byte) (Format, Confidence) {
	format := CSV
	if delim == '\t' {
		format = TSV
	}

	// Check up to 10 non-empty lines.
	maxLines := 10
	counts := make([]int, 0, maxLines)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		counts = append(counts, strings.Count(line, string(delim)))
		if len(counts) >= maxLines {
			break
		}
	}

	if len(counts) < 2 {
		return format, None
	}

	// All lines must have the same delimiter count, and it must be > 0.
	expected := counts[0]
	if expected == 0 {
		return format, None
	}

	consistent := 0
	for _, c := range counts[1:] {
		if c == expected {
			consistent++
		}
	}

	total := len(counts) - 1
	if consistent == total && total >= 4 {
		return format, High
	}
	if consistent == total && total >= 1 {
		return format, Medium
	}

	return format, None
}
