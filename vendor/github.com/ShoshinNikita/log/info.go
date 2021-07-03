package log

import (
	"fmt"
)

// Info prints info message
// Output pattern: (?time) [INFO] msg
func (l Logger) Info(v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprint(v...), l.getTime, l.getInfoMsg))
}

// Infof prints info message
// Output pattern: (?time) [INFO] msg
func (l Logger) Infof(format string, v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprintf(format, v...), l.getTime, l.getInfoMsg))
}

// Infoln prints info message
// Output pattern: (?time) [INFO] msg
func (l Logger) Infoln(v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprintln(v...), l.getTime, l.getInfoMsg))
}
