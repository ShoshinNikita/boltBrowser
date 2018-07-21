package log

import (
	"fmt"
)

func Error(v ...interface{}) {
	if showTime {
		printTime()
	}
	printErrMsg()
	fmt.Print(v...)
}

func Errorf(format string, v ...interface{}) {
	if showTime {
		printTime()
	}
	printErrMsg()
	fmt.Printf(format, v...)
}

func Errorln(v ...interface{}) {
	if showTime {
		printTime()
	}
	printErrMsg()
	fmt.Println(v...)
}
