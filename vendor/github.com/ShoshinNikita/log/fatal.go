package log

import (
	"fmt"
	"os"
)

// Fatal prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func (l Logger) Fatal(v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprint(v...), l.getTime, l.getFatalMsg, l.getCaller))
	os.Exit(1)
}

// Fatalf prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func (l Logger) Fatalf(format string, v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprintf(format, v...), l.getTime, l.getFatalMsg, l.getCaller))
	os.Exit(1)
}

// Fatalln prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func (l Logger) Fatalln(v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprint(v...), l.getTime, l.getFatalMsg, l.getCaller))
	os.Exit(1)
}
