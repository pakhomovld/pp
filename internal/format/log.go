package format

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pakhomovld/ppp/internal/color"
)

var (
	logTimestampRe = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}[.\d]*)(\S*)`)
	logLevelRe     = regexp.MustCompile(`(?i)\b(INFO|ERROR|WARN|WARNING|DEBUG|TRACE|FATAL)\b`)
	logBracketedRe = regexp.MustCompile(`(?i)\[(INFO|ERROR|WARN|WARNING|DEBUG|TRACE|FATAL)\]`)
)

// LogFormatter colorizes log lines, processing them one at a time (streaming).
type LogFormatter struct{}

func (f *LogFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		// Try to extract and format embedded JSON.
		if strings.ContainsAny(line, "{[") {
			if prefix, v, found := extractJSON(line); found {
				prefix = strings.TrimRight(prefix, " \t")
				if theme != nil {
					prefix = colorizeLogLine(prefix, theme)
				}
				if _, err := fmt.Fprintln(w, prefix); err != nil {
					return err
				}
				out := formatValue(v, 1, theme)
				if _, err := fmt.Fprintln(w, out); err != nil {
					return err
				}
				continue
			}
		}

		if theme != nil {
			line = colorizeLogLine(line, theme)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func colorizeLogLine(line string, theme *color.Theme) string {
	result := line

	// Colorize timestamp.
	result = logTimestampRe.ReplaceAllStringFunc(result, func(match string) string {
		return theme.Sprint(color.Comment, match)
	})

	// Colorize bracketed levels: [ERROR], [INFO], etc.
	result = logBracketedRe.ReplaceAllStringFunc(result, func(match string) string {
		return colorizeLevel(match, theme)
	})

	// Colorize bare levels: ERROR, INFO, etc. (only if not already colored).
	if !logBracketedRe.MatchString(line) {
		result = logLevelRe.ReplaceAllStringFunc(result, func(match string) string {
			return colorizeLevel(match, theme)
		})
	}

	return result
}

// extractJSON finds the first valid JSON object or array embedded in a line.
// Returns the prefix before the JSON, the parsed value, and whether JSON was found.
func extractJSON(line string) (prefix string, v any, found bool) {
	for i, ch := range line {
		if ch != '{' && ch != '[' {
			continue
		}
		candidate := line[i:]
		var parsed any
		if err := json.Unmarshal([]byte(candidate), &parsed); err == nil {
			return line[:i], parsed, true
		}
	}
	return "", nil, false
}

func colorizeLevel(level string, theme *color.Theme) string {
	upper := strings.ToUpper(level)
	switch {
	case strings.Contains(upper, "ERROR"), strings.Contains(upper, "FATAL"):
		return theme.Sprint(color.Null, level) // Red.
	case strings.Contains(upper, "WARN"):
		return theme.Sprint(color.Number, level) // Yellow.
	case strings.Contains(upper, "INFO"):
		return theme.Sprint(color.Boolean, level) // Magenta (stands out).
	case strings.Contains(upper, "DEBUG"), strings.Contains(upper, "TRACE"):
		return theme.Sprint(color.Comment, level) // Dim.
	default:
		return level
	}
}
