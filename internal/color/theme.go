package color

import (
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// Token types for syntax coloring.
type Token int

const (
	Key Token = iota
	String
	Number
	Boolean
	Null
	Bracket
	Colon
	Comma
	Comment
)

// Theme maps token types to color attributes.
type Theme struct {
	colors map[Token]*color.Color
}

// DefaultTheme returns a color theme for terminal output.
func DefaultTheme() *Theme {
	return &Theme{
		colors: map[Token]*color.Color{
			Key:     color.New(color.FgCyan, color.Bold),
			String:  color.New(color.FgGreen),
			Number:  color.New(color.FgYellow),
			Boolean: color.New(color.FgMagenta),
			Null:    color.New(color.FgRed),
			Bracket: color.New(color.FgWhite),
			Colon:   color.New(color.FgWhite),
			Comma:   color.New(color.FgWhite),
			Comment: color.New(color.FgHiBlack),
		},
	}
}

// Sprint returns the text colored for the given token type.
// If the theme is nil, the text is returned unmodified.
func (t *Theme) Sprint(tok Token, text string) string {
	if t == nil {
		return text
	}
	if c, ok := t.colors[tok]; ok {
		return c.Sprint(text)
	}
	return text
}

// ShouldColor returns true if output should be colored.
func ShouldColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}
