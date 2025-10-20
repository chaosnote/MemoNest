package config

import "flag"

type CLIFlags struct {
	Debug bool
}

func ParseCLIFlags() *CLIFlags {
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	return &CLIFlags{
		Debug: *debug,
	}
}
