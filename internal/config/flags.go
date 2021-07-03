package config

import "flag"

// parseFlags parses command line flags
func parseFlags() {
	// We use fields of Opts as default, because we want to overwrite values of only setted flags
	flag.StringVar(&Opts.Port, "port", Opts.Port, "port for website (with ':')")
	flag.BoolVar(&Opts.Debug, "debug", Opts.Debug, "debug mode")
	flag.IntVar(&Opts.Offset, "offset", Opts.Offset, "number of records on single page")
	flag.BoolVar(&Opts.OpenBrowser, "openBrowser", Opts.OpenBrowser, "should the program open a browser automatically")
	flag.BoolVar(&Opts.NeatWindow, "neatWindow", Opts.NeatWindow, "should the program open a neat window")
	flag.Parse()

	// Checking of ':' before port
	if Opts.Port[0] != ':' {
		Opts.Port = ":" + Opts.Port
	}
}
