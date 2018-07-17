package web

import (
	"errors"
	"html"
	"reflect"
)

// escapeRecords escapes fields "Key" and "Value" of a field "Records"
// See tests for details
func escapeRecords(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return errors.New("s isn't a struct")
	}

	// We want to change only field "Records"
	field := v.FieldByName("Records")
	if !field.IsValid() {
		return errors.New("Field 'Records' isn't valid")
	}
	if field.Kind() != reflect.Slice {
		return errors.New("Field 'Records' isn't a slice")
	}
	if field.Type().Elem().Kind() != reflect.Struct {
		return errors.New("Field 'Recrod' isn't a slice of structs")
	}

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

	return nil
}

// escape is a wrapper for html.EscapeString()
func escape(s string) string {
	return html.EscapeString(s)
}

// checkString checks is v valid and is v.Kind() reflect.String
func checkString(v reflect.Value) error {
	if !v.IsValid() {
		return errors.New("Field isn't valid")
	}

	if v.Kind() != reflect.String {
		return errors.New("Field's kind isn't a string")
	}

	return nil
}
