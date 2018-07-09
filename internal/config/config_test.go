package config_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ShoshinNikita/boltBrowser/internal/config"
)

// s - struct
func compare(s interface{}, values []interface{}) error {
	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != values[i] {
			return fmt.Errorf("Field #%d Want: '%v' Got: '%v'", i, values[i], v.Field(i).Interface())
		}
	}

	return nil
}

func TestSetDefaultValues(t *testing.T) {
	Test1 := struct {
		Hello string `default:"hello"`
	}{}
	values := []interface{}{"hello"}

	config.SetDefaultValues(&Test1)
	t.Logf("%+v\n", Test1)
	err := compare(Test1, values)
	if err != nil {
		t.Error(err)
	}

	Test2 := struct {
		A int    `default:"15"`
		B bool   `default:"false"`
		C string `default:"123"`
		D bool   `default:"true"`
	}{}
	values = []interface{}{15, false, "123", true}

	config.SetDefaultValues(&Test2)
	t.Logf("%+v\n", Test2)
	err = compare(Test2, values)
	if err != nil {
		t.Error(err)
	}

	Test3 := struct {
		A int
		B bool
		C string
		D bool
	}{}
	values = []interface{}{0, false, "", false}

	config.SetDefaultValues(&Test3)
	t.Logf("%+v\n", Test3)
	err = compare(Test3, values)
	if err != nil {
		t.Error(err)
	}

	Test4 := struct {
		A int  `default:"15"`
		B bool `default:"true"`
		C string
		D bool
	}{}
	values = []interface{}{15, true, "", false}

	config.SetDefaultValues(&Test4)
	t.Logf("%+v\n", Test4)
	err = compare(Test4, values)
	if err != nil {
		t.Error(err)
	}
}
