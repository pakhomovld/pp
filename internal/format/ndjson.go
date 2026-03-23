package format

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pakhomovld/ppp/internal/color"
)

// NDJSONFormatter pretty-prints newline-delimited JSON, one object per line.
type NDJSONFormatter struct{}

func (f *NDJSONFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()

		// Preserve empty lines.
		if len(bytes.TrimSpace(line)) == 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
			continue
		}

		var v any
		if err := json.Unmarshal(line, &v); err != nil {
			// Invalid line — pass through raw.
			if _, err := fmt.Fprintln(w, string(line)); err != nil {
				return err
			}
			continue
		}

		out := formatValue(v, 0, theme)
		if _, err := fmt.Fprintln(w, out); err != nil {
			return err
		}
	}
	return scanner.Err()
}
