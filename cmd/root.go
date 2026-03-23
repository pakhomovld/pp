package cmd

import (
	"flag"
	"fmt"
	"os"
)

// Version is set via ldflags at build time.
var Version = "dev"

// Config holds CLI flags.
type Config struct {
	ForceFormat string
	NoColor     bool
	Version     bool
}

// ParseFlags parses CLI arguments.
func ParseFlags() Config {
	var cfg Config

	flag.StringVar(&cfg.ForceFormat, "format", "", "force a specific format (json, yaml, csv, ...)")
	flag.StringVar(&cfg.ForceFormat, "f", "", "force a specific format (shorthand)")
	flag.BoolVar(&cfg.NoColor, "no-color", false, "disable colored output")
	flag.BoolVar(&cfg.Version, "version", false, "print version")
	flag.BoolVar(&cfg.Version, "v", false, "print version (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "pp — universal pipe pretty-printer\n\n")
		fmt.Fprintf(os.Stderr, "Usage: <command> | pp [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	return cfg
}
