package log

import (
	"fmt"
)

// Warn prints warning
// Output pattern: (?time) [WARN] warning
func (l Logger) Warn(v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprint(v...), l.getTime, l.getWarnMsg))
}

// Warnf prints warning
// Output pattern: (?time) [WARN] warning
func (l Logger) Warnf(format string, v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprintf(format, v...), l.getTime, l.getWarnMsg))
}

// Warnln prints warning
// Output pattern: (?time) [WARN] warning
func (l Logger) Warnln(v ...interface{}) {
	l.printText(addPrefixes(fmt.Sprintln(v...), l.getTime, l.getWarnMsg))
}
