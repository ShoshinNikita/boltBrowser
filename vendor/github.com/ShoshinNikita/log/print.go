package log

import (
	"fmt"
)

func Print(v ...interface{}) {
	fmt.Print(v...)
}

func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func Println(v ...interface{}) {
	fmt.Println(v...)
}
