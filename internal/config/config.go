// Package config
//
// Program can get Opts from:
// 1. Command line flags
// 2. config.ini
//    Types of lines:
//      * # - comment
//      * some option=value of the option
//
package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Opts keeps vars for global use
//
// Tags:
// * default - default value
// * name - name in the file config.ini
// * comment - comment in the file config.ini
var Opts struct {
	// Port for website (with ':')
	Port string `default:":500" name:"port" comment:"Port for website"`
	// Debug mode
	Debug bool `default:"false" name:"debug"`
	// Offset - number of records on single screen
	Offset int `default:"100" name:"offset" comment:"number of records on a single screen"`
	// CheckVer - should the program check check is there a new version
	CheckVer bool `default:"true" name:"should_check_version"`
	// IsWriteMode - can program edit databases
	IsWriteMode bool `default:"true" name:"is_write_mode"`
	// OpenBrowser - should the program open a browser automatically
	OpenBrowser bool `default:"true" name:"open_browser"`
	// NeatWindow - should the program open the special neat window
	NeatWindow bool `default:"true" name:"open_neat_window" comment:"has effect only if 'open browser' is true"`
}

type field struct {
	name  string
	value interface{}
}

// ParseConfig parses flags like -port, -debug, -offset and etc.
// If there's no any flags, it tries to parse config file "config.ini"
func ParseConfig() (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case string:
				err = fmt.Errorf(r.(string))
			case error:
				err = r.(error)
			default:
				err = fmt.Errorf("undefined error")
			}
		}
	}()

	setDefaultValues(&Opts)

	// At first, we parse the file
	parseFile()

	// If user set any flag, the program will overwrite Opts
	if len(os.Args) > 1 {
		parseFlags()
	}

	return nil
}

// setDefaultValues sets default values of Opts's fields.
// s - pointer to a struct.
// If tag default was missed it panics
// If type of field isn't [int, string, bool] it panics
func setDefaultValues(s interface{}) {
	// value of a pointer
	value := reflect.ValueOf(s)
	t := value.Type()
	if t.Kind() != reflect.Ptr {
		panicf("s must be a pointer")
	}

	// value of struct (ptr -> struct)
	value = value.Elem()
	t = value.Type()

	if t.Kind() != reflect.Struct {
		panicf("s isn't structure, but ", t.Kind().String())
	}

	var defValues []field

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		def := f.Tag.Get("default")
		// If there's no tag default, skip this field
		if def == "" {
			continue
		}

		switch f.Type.Kind() {
		case reflect.Bool:
			b := func() bool {
				return def == "true"
			}()
			defValues = append(defValues, field{name: f.Name, value: b})
		case reflect.String:
			defValues = append(defValues, field{name: f.Name, value: def})
		case reflect.Int:
			i, _ := strconv.ParseInt(def, 10, 64)
			defValues = append(defValues, field{name: f.Name, value: int(i)})
		default:
			panicf("Bad type of a field of Opts. Type: %s", f.Type.Kind())
		}
	}

	setValues(s, defValues)
}

// s - pointer to a struct
func setValues(s interface{}, values []field) {
	// Checking of type
	// value of a pointer
	value := reflect.ValueOf(s)
	t := value.Type()
	if t.Kind() != reflect.Ptr {
		panicf("s must be a pointer")
	}

	// value of struct (ptr -> struct)
	if value.Elem().Type().Kind() != reflect.Struct {
		panicf("s isn't structure, but ", t.Kind().String())
	}

	// value of the struct (ptr -> struct)
	opts := reflect.ValueOf(s).Elem()

	for _, v := range values {
		f := opts.FieldByName(v.name)
		if f.IsValid() {
			if f.Kind() != reflect.TypeOf(v.value).Kind() {
				panicf("Different types of field and value: field type - %s, value type - %s", f.Kind().String(), reflect.TypeOf(v.value).Kind().String())
			}

			f.Set(reflect.ValueOf(v.value))
		}
	}
}

func panicf(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v))
}
