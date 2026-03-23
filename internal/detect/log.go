package detect

import (
	"bytes"
	"regexp"
)

var logPatterns = []*regexp.Regexp{
	// ISO8601 timestamp: 2024-01-15T10:30:00 or 2024-01-15 10:30:00
	regexp.MustCompile(`^\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}`),
	// Bracketed level: [INFO], [ERROR], [WARN], [DEBUG]
	regexp.MustCompile(`(?i)\[(INFO|ERROR|WARN|WARNING|DEBUG|TRACE|FATAL)\]`),
	// Level prefix: INFO ..., ERROR ...
	regexp.MustCompile(`(?m)^(?:\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}\S*\s+)?(?:INFO|ERROR|WARN|WARNING|DEBUG|TRACE|FATAL)\s`),
}

// LogDetector detects log-formatted text.
type LogDetector struct{}

func (d *LogDetector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) == 0 {
		return Result{Format: LogLine, Confidence: None}
	}

	// Skip structured formats (but not bracketed log levels like [INFO]).
	if trimmed[0] == '{' || trimmed[0] == '<' {
		return Result{Format: LogLine, Confidence: None}
	}

	lines := bytes.SplitN(trimmed, []byte("\n"), 6)
	matchCount := 0

	for _, line := range lines {
		for _, pat := range logPatterns {
			if pat.Match(line) {
				matchCount++
				break
			}
		}
	}

	if matchCount >= 3 {
		return Result{Format: LogLine, Confidence: High}
	}
	if matchCount >= 2 {
		return Result{Format: LogLine, Confidence: Medium}
	}
	if matchCount >= 1 {
		return Result{Format: LogLine, Confidence: Low}
	}

	return Result{Format: LogLine, Confidence: None}
}
