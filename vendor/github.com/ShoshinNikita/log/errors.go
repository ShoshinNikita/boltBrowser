package log

import (
	"fmt"
)

// Error prints error
// Output pattern: (?time) [ERR] (?file:line) error
func (l Logger) Error(v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprint(v...), l.getTime, l.getErrMsg, l.getCaller))
}

// Errorf prints error
// Output pattern: (?time) [ERR] (?file:line) error
func (l Logger) Errorf(format string, v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprintf(format, v...), l.getTime, l.getErrMsg, l.getCaller))
}

// Errorln prints error
// Output pattern: (?time) [ERR] (?file:line) error
func (l Logger) Errorln(v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprintln(v...), l.getTime, l.getErrMsg, l.getCaller))
}
