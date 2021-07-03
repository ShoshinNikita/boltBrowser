package log

import (
	"fmt"
)

// Print prints msg
// Output pattern: (?time) msg
func (l Logger) Print(v ...interface{}) {
	l.printText(fmt.Sprint(v...))
}

// Printf prints msg
// Output pattern: (?time) msg
func (l Logger) Printf(format string, v ...interface{}) {
	l.printText(fmt.Sprintf(format, v...))
}

// Println prints msg
// Output pattern: (?time) msg
func (l Logger) Println(v ...interface{}) {
	l.printText(fmt.Sprintln(v...))
}
