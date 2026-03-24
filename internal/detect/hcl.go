package detect

import (
	"bytes"
	"regexp"
)

var (
	hclBlockHeaderRe = regexp.MustCompile(`(?m)^\s*(resource|variable|output|module|provider|data|locals|terraform)\s+"[^"]*"`)
	hclAssignmentRe  = regexp.MustCompile(`(?m)^\s*[\w\-]+\s*=\s*.+`)
	hclBraceBlockRe  = regexp.MustCompile(`\{`)
	tomlSectionGuard = regexp.MustCompile(`(?m)^\s*\[[\w.\-]+\]\s*$`)
)

// HCLDetector detects HCL/Terraform configuration.
type HCLDetector struct{}

func (d *HCLDetector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) == 0 {
		return Result{Format: HCL, Confidence: None}
	}

	// Skip JSON, JSON arrays, and XML.
	if trimmed[0] == '{' || trimmed[0] == '[' || trimmed[0] == '<' {
		return Result{Format: HCL, Confidence: None}
	}

	blockHeaders := hclBlockHeaderRe.FindAll(trimmed, -1)
	if len(blockHeaders) >= 2 {
		return Result{Format: HCL, Confidence: High}
	}

	assignments := hclAssignmentRe.FindAll(trimmed, -1)
	hasBraces := hclBraceBlockRe.Match(trimmed)
	hasTOMLSections := tomlSectionGuard.Match(trimmed)

	if len(blockHeaders) == 1 {
		return Result{Format: HCL, Confidence: Medium}
	}

	// HCL-style assignments with braces but no TOML sections.
	if len(assignments) >= 3 && hasBraces && !hasTOMLSections {
		return Result{Format: HCL, Confidence: Medium}
	}

	return Result{Format: HCL, Confidence: None}
}
