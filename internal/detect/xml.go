package detect

import (
	"bytes"
	"regexp"
)

var htmlDoctype = regexp.MustCompile(`(?i)^<!doctype\s+html`)
var htmlTag = regexp.MustCompile(`(?i)^<html[\s>]`)

// XMLDetector detects XML and HTML documents.
type XMLDetector struct{}

func (d *XMLDetector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) == 0 || trimmed[0] != '<' {
		return Result{Format: XML, Confidence: None}
	}

	// Check for HTML first.
	if htmlDoctype.Match(trimmed) || htmlTag.Match(trimmed) {
		return Result{Format: HTML, Confidence: High}
	}

	// XML declaration.
	if bytes.HasPrefix(trimmed, []byte("<?xml")) {
		return Result{Format: XML, Confidence: High}
	}

	// Looks like a tag — check for matching close or self-closing.
	if trimmed[0] == '<' && len(trimmed) > 1 && isLetter(trimmed[1]) {
		return Result{Format: XML, Confidence: Medium}
	}

	return Result{Format: XML, Confidence: None}
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
