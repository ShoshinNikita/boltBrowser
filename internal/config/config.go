package config

import (
	"flag"
)

// TODO: don't use global options
var Opts struct {
	// Port for website
	Port int
	// Offset defines the number of records on single screen
	// TODO: remove
	Offset int
	// OpenBrowser defines should the program open a browser automatically
	OpenBrowser bool
}

func ParseConfig() error {
	flag.IntVar(&Opts.Port, "port", 8080, "port for website")
	flag.IntVar(&Opts.Offset, "offset", 100, "number of records on single page")
	flag.BoolVar(&Opts.OpenBrowser, "openBrowser", true, "should the program open a browser automatically")

	flag.Parse()

	return nil
}
