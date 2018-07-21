package log

import (
	"fmt"
)

func Info(v ...interface{}) {
	if showTime {
		printTime()
	}
	printInfoMsg()
	fmt.Print(v...)
}

func Infof(format string, v ...interface{}) {
	if showTime {
		printTime()
	}
	printInfoMsg()
	fmt.Printf(format, v...)
}

func Infoln(v ...interface{}) {
	if showTime {
		printTime()
	}
	printInfoMsg()
	fmt.Println(v...)
}
