package format

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pakhomovld/ppp/internal/color"
)

var (
	hclBlockLineRe   = regexp.MustCompile(`^(\s*)(resource|variable|output|module|provider|data|locals|terraform)\b(.*)$`)
	hclKVLineRe      = regexp.MustCompile(`^(\s*)([\w\-]+)(\s*=\s*)(.+)$`)
	hclCommentLineRe = regexp.MustCompile(`^(\s*)(#.*|//.*)$`)
	hclBraceOpenRe   = regexp.MustCompile(`^(\s*)(.*)(\{)\s*$`)
	hclBraceCloseRe  = regexp.MustCompile(`^(\s*)(\})\s*$`)
)

// HCLFormatter colorizes HCL/Terraform output line by line.
type HCLFormatter struct{}

func (f *HCLFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	if theme == nil {
		_, err := io.Copy(w, r)
		return err
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		colored := colorizeHCLLine(line, theme)
		if _, err := fmt.Fprintln(w, colored); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func colorizeHCLLine(line string, theme *color.Theme) string {
	// Comment.
	if m := hclCommentLineRe.FindStringSubmatch(line); m != nil {
		return m[1] + theme.Sprint(color.Comment, m[2])
	}

	// Block header (resource "type" "name" {).
	if m := hclBlockLineRe.FindStringSubmatch(line); m != nil {
		indent, keyword, rest := m[1], m[2], m[3]
		result := indent + theme.Sprint(color.Key, keyword)
		result += colorizeHCLRest(rest, theme)
		return result
	}

	// Closing brace.
	if m := hclBraceCloseRe.FindStringSubmatch(line); m != nil {
		return m[1] + theme.Sprint(color.Bracket, m[2])
	}

	// Key = value.
	if m := hclKVLineRe.FindStringSubmatch(line); m != nil {
		indent, key, sep, val := m[1], m[2], m[3], m[4]
		return indent + theme.Sprint(color.Key, key) +
			theme.Sprint(color.Colon, sep) +
			colorizeHCLValue(val, theme)
	}

	return line
}

func colorizeHCLRest(rest string, theme *color.Theme) string {
	var b strings.Builder
	parts := strings.Fields(rest)
	for i, part := range parts {
		if i > 0 || len(rest) > 0 && rest[0] == ' ' {
			b.WriteString(" ")
		}
		trimmed := strings.TrimSpace(part)
		if trimmed == "{" || trimmed == "}" {
			b.WriteString(theme.Sprint(color.Bracket, part))
		} else if isQuotedString(trimmed) {
			b.WriteString(theme.Sprint(color.String, part))
		} else {
			b.WriteString(part)
		}
	}
	return b.String()
}

func colorizeHCLValue(val string, theme *color.Theme) string {
	trimmed := strings.TrimSpace(val)

	switch {
	case trimmed == "true" || trimmed == "false":
		return theme.Sprint(color.Boolean, val)
	case isQuotedString(trimmed):
		return theme.Sprint(color.String, val)
	case isNumeric(trimmed):
		return theme.Sprint(color.Number, val)
	case strings.HasPrefix(trimmed, "["):
		return theme.Sprint(color.Bracket, val)
	case strings.HasPrefix(trimmed, "{"):
		return theme.Sprint(color.Bracket, val)
	default:
		return val
	}
}
