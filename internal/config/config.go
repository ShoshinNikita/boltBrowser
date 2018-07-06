package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

var Opts struct {
	// Port for website (with ':')
	Port string `default:":500"`
	// Debug mode
	Debug bool `default:"false"`
	// Offset - number of records on single screen
	Offset int `default:"100"`
	// CheckVer - should the program check check is there a new version
	CheckVer bool `default:"true"`
	// IsWriteMode - can program edit databases
	IsWriteMode bool `default:"true"`
	// OpenBrowser - should the program open a browser automatically
	OpenBrowser bool `default:"true"`
	// NeatWindow - should the program open the special neat window
	NeatWindow bool `default:"true"`
}

// ParseConfig parses flags like -port, -debug, -offset and etc.
// If there's no any flags, it tries to parse config file "config.ini"
func ParseConfig() {
	defer func() {
		if r := recover(); r != nil {
			exit(r)
		}
	}()

	setDefaultValues()

	if len(os.Args) > 1 {
		parseFlags()
	}
}

// setDefaultValues sets default values of Opts's fields.
// If tag default was missed it panics
// If type of field isn't [int, string, bool] it panics
func setDefaultValues() {
	// For change value
	v := reflect.ValueOf(&Opts).Elem()
	// For getting tags
	t := v.Type()

	// Opts is always struct, so we shouldn't check if tp.Kind() == reflect.Struct
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		def := t.Field(i).Tag.Get("default")
		if def == "" {
			panicf("default tag of field %s is empty", t.Field(i).Name)
		}

		switch f.Kind() {
		case reflect.Bool:
			b := func() bool {
				if def == "true" {
					return true
				}
				return false
			}()
			f.SetBool(b)
		case reflect.String:
			f.SetString(def)
		case reflect.Int:
			i, _ := strconv.ParseInt(def, 10, 64)
			f.SetInt(i)
		default:
			panicf("Bad type of a field if Opts. Type: %s", f.Kind().String())
		}
	}
}

// parseFlags parses command line flags
func parseFlags() {
	// We can use fields of Opts as default, because we already set default values by calling of setDefaultValues()
	flag.StringVar(&Opts.Port, "port", Opts.Port, "port for website (with ':')")
	flag.BoolVar(&Opts.Debug, "debug", Opts.Debug, "debug mode")
	flag.IntVar(&Opts.Offset, "offset", Opts.Offset, "number of records on single page")
	flag.BoolVar(&Opts.CheckVer, "checkVer", Opts.CheckVer, "should program check is there a new version")
	flag.BoolVar(&Opts.IsWriteMode, "writeMode", Opts.IsWriteMode, "can program edit dbs")
	flag.BoolVar(&Opts.OpenBrowser, "openBrowser", Opts.OpenBrowser, "should the program open a browser automatically")
	flag.BoolVar(&Opts.NeatWindow, "neatWindow", Opts.NeatWindow, "should the program open a neat window")
	flag.Parse()

	// Checking of ':' before port
	if Opts.Port[0] != ':' {
		Opts.Port = ":" + Opts.Port
	}
}

func panicf(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v))
}

func exit(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(2)
}
