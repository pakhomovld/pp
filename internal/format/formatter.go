package format

import (
	"io"

	"github.com/pakhomovld/ppp/internal/color"
)

// Formatter pretty-prints data from a reader to a writer.
type Formatter interface {
	Format(w io.Writer, r io.Reader, theme *color.Theme) error
}
