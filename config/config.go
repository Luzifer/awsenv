package config

import "flag"

type Config struct {
	AWSEnv        string
	Debug         bool
	Export        bool
	OverrideShell string
}

func New() *Config {
	var (
		debug         = flag.Bool("d", false, "Enable debug output")
		export        = flag.Bool("x", false, "Set vars as exported")
		overrideShell = flag.String("shell", "", "Overrides shell detection from $SHELL var")
	)

	flag.Parse()

	return &Config{
		AWSEnv:        flag.Arg(0),
		Debug:         *debug,
		Export:        *export,
		OverrideShell: *overrideShell,
	}
}
