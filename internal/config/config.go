package config

import (
	"flag"
	"os"
)

var Opts struct {
	// Port for website (with ':')
	Port string
	// Debug mode
	Debug bool
	// Offset - number of records on single screen
	Offset int
	// CheckVer - should the program check check is there a new version
	CheckVer bool
	// IsWriteMode - can program edit databases
	IsWriteMode bool
	// OpenBrowser - should the program open a browser automatically
	OpenBrowser bool
	// NeatWindow - should the program open the special neat window
	NeatWindow bool
}

// ParseConfig parses flags like -port, -debug, -offset and etc.
// If there's no any flags, it tries to parse config file "config.ini"
func ParseConfig() {
	if len(os.Args) > 1 {
		parseFlags()
	}
}

// parseFlags parses command line flags
func parseFlags() {
	flag.StringVar(&Opts.Port, "port", ":500", "port for website (with ':')")
	flag.BoolVar(&Opts.Debug, "debug", false, "debug mode")
	flag.IntVar(&Opts.Offset, "offset", 100, "number of records on single page")
	flag.BoolVar(&Opts.CheckVer, "checkVer", true, "should program check is there a new version")
	flag.BoolVar(&Opts.IsWriteMode, "writeMode", true, "can program edit dbs")
	flag.BoolVar(&Opts.OpenBrowser, "openBrowser", true, "should the program open a browser automatically")
	flag.BoolVar(&Opts.NeatWindow, "neatWindow", true, "should the program open a neat window")
	flag.Parse()

	// Checking of ':' before port
	if Opts.Port[0] != ':' {
		Opts.Port = ":" + Opts.Port
	}
}

