// Package params serve as environment
package params

import (
	"flag"
)

var (
	Port        string
	Debug       bool
	Offset      int
	CheckVer    bool
	IsWriteMode bool
)

func ParseFlags() {
	flag.StringVar(&Port, "port", ":500", "port for website (with ':')")
	flag.BoolVar(&Debug, "debug", false, "debug mode")
	flag.IntVar(&Offset, "offset", 100, "number of records on single page")
	flag.BoolVar(&CheckVer, "checkVer", true, "should program check is there a new version")
	flag.BoolVar(&IsWriteMode, "writeMode", true, "can program edit dbs")
	flag.Parse()

	// Checking of ':' before port
	if Port[0] != ':' {
		Port = ":" + Port
	}
}
