package web

import (
	"errors"
	"html"
	"reflect"
)

// escapeRecords escapes fields "Key" and "Value" of a field "Records" and field "Path".
// s must be a pointer to a struct.
//
// Only fields "Records" and "Path" can have js-injected data
// See tests for details
func escapeRecords(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("s isn't a ptr")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("s isn't a struct")
	}

	// Change the field "Records"
	field := v.FieldByName("Records")
	// Skip if the field isn't valid and it isn't slice of structs
	if field.IsValid() && field.Kind() == reflect.Slice &&
		field.Type().Elem().Kind() == reflect.Struct {

		for i := 0; i < field.Len(); i++ {
			elem := field.Index(i)

			// Change fields "Key" and "Value"
			key := elem.FieldByName("Key")
			if err := checkString(key); err != nil {
				return err
			}

			// We can use SetString(), because we already have checked kind of 'key'
			key.SetString(escape(key.String()))

			value := elem.FieldByName("Value")
			if err := checkString(value); err != nil {
				return err
			}

			// We can use SetString(), because we already have checked kind of 'value'
			value.SetString(escape(value.String()))
		}
	}

	// Change the field "Path"
	field = v.FieldByName("Path")
	// Skip if the field isn't valid and it isn't string
	if field.IsValid() && field.Kind() == reflect.String && field.CanSet() {
		field.SetString(escape(field.String()))
	}

	return nil
}

// escape is a wrapper for html.EscapeString()
func escape(s string) string {
	return html.EscapeString(s)
}

// checkString checks is v valid and is v.Kind() reflect.String
func checkString(v reflect.Value) error {
	if !v.IsValid() {
		return errors.New("field isn't valid")
	}

	if v.Kind() != reflect.String {
		return errors.New("field's kind isn't a string")
	}

	return nil
}
