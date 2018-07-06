package config

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// parseFile parses config.ini
// If config.ini doesn't exist, program will create file with default values.
// If there's an error, program will exit with code 2.
func parseFile() {
	if _, err := os.Open("config.ini"); os.IsNotExist(err) {
		// If we have to create a new file, we don't need to parse it
		err = createFile()
		if err != nil {
			panic(err)
		}

		return
	}

	// Read data from the file
	file, _ := os.Open("config.ini")

	scanner := bufio.NewScanner(file)

	opts := make(map[string]string)

	for scanner.Scan() {
		line := string(scanner.Bytes())
		// Skip empty strings and comments
		if line != "" && line[0] != '#' {
			data := strings.Split(line, "=")
			if len(data) == 2 {
				opts[data[0]] = data[1]
			}
		}
	}

	file.Close()

	// Set values
	var values []field

	t := reflect.TypeOf(Opts)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		configName := f.Tag.Get("name")
		fieldName := f.Name

		// If there's no such name, skip
		if _, ok := opts[configName]; !ok {
			continue
		}

		switch f.Type.Kind() {
		case reflect.Int:
			v, _ := strconv.ParseInt(opts[configName], 10, 64)
			values = append(values, field{name: fieldName, value: int(v)})
		case reflect.String:
			values = append(values, field{name: fieldName, value: opts[configName]})
		case reflect.Bool:
			v, _ := strconv.ParseBool(opts[configName])
			values = append(values, field{name: fieldName, value: v})
		default:
			panicf("Bad type of a field of Opts. Type: %s", f.Type.Kind())
		}

	}

	setValues(values)
}

// createFile creates config.ini and writes Opts.
func createFile() error {
	file, err := os.Create("config.ini")
	if err != nil {
		return err
	}

	// We can use Opts, because we already set default values by calling of setDefaultValues()
	v := reflect.ValueOf(Opts) // for values
	t := v.Type()              // for tags

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		configName := f.Tag.Get("name")
		comment := f.Tag.Get("comment")
		defValue := fmt.Sprint(v.Field(i).Interface())

		var line string
		if comment != "" {
			line = "# " + comment + "\n"
		}

		line += configName + "=" + defValue + "\n"

		file.Write([]byte(line))
	}

	return nil
}
