package log

import (
	"io"
)

var globalLogger *Logger

// init inits globalLogger with NewLogger()
func init() {
	globalLogger = NewLogger()

	globalLogger.PrintTime(false)
	globalLogger.PrintColor(false)
	globalLogger.PrintErrorLine(false)

	globalLogger.global = true
}

// PrintTime sets globalLogger.PrintTime
// Time isn't printed by default
func PrintTime(b bool) {
	globalLogger.PrintTime(b)
}

// ShowTime sets printTime
// Time isn't printed by default
//
// It was left for backwards compatibility
var ShowTime = PrintTime

// PrintColor sets printColor
// printColor is false by default
func PrintColor(b bool) {
	globalLogger.PrintColor(b)
}

// PrintErrorLine sets PrintErrorLine
// If PrintErrorLine is true, log.Error(), log.Errorf(), log.Errorln() will print file and line,
// where functions were called.
// PrintErrorLine is false by default
func PrintErrorLine(b bool) {
	globalLogger.PrintErrorLine(b)
}

// ChangeOutput changes Logger.output writer.
// Default Logger.output is github.com/fatih/color.Output
func ChangeOutput(w io.Writer) {
	globalLogger.ChangeOutput(w)
}

// ChangeTimeLayout changes Logger.timeLayout
// Default Logger.timeLayout is DefaultTimeLayout
func ChangeTimeLayout(layout string) {
	globalLogger.ChangeTimeLayout(layout)
}

/* Print */

// Print prints msg
// Output pattern: (?time) msg
func Print(v ...interface{}) {
	globalLogger.Print(v...)
}

// Printf prints msg
// Output pattern: (?time) msg
func Printf(format string, v ...interface{}) {
	globalLogger.Printf(format, v...)
}

// Println prints msg
// Output pattern: (?time) msg
func Println(v ...interface{}) {
	globalLogger.Println(v...)
}

/* Info */

// Info prints info message
// Output pattern: (?time) [INFO] msg
func Info(v ...interface{}) {
	globalLogger.Info(v...)
}

// Infof prints info message
// Output pattern: (?time) [INFO] msg
func Infof(format string, v ...interface{}) {
	globalLogger.Infof(format, v...)
}

// Infoln prints info message
// Output pattern: (?time) [INFO] msg
func Infoln(v ...interface{}) {
	globalLogger.Infoln(v...)
}

/* Warn */

// Warn prints warning
// Output pattern: (?time) [WARN] warning
func Warn(v ...interface{}) {
	globalLogger.Warn(v...)
}

// Warnf prints warning
// Output pattern: (?time) [WARN] warning
func Warnf(format string, v ...interface{}) {
	globalLogger.Warnf(format, v...)
}

// Warnln prints warning
// Output pattern: (?time) [WARN] warning
func Warnln(v ...interface{}) {
	globalLogger.Warnln(v...)
}

/* Error */

// Error prints error
// Output pattern: (?time) [ERR] (?file:line) error
func Error(v ...interface{}) {
	globalLogger.Error(v...)
}

// Errorf prints error
// Output pattern: (?time) [ERR] (?file:line) error
func Errorf(format string, v ...interface{}) {
	globalLogger.Errorf(format, v...)
}

// Errorln prints error
// Output pattern: (?time) [ERR] (?file:line) error
func Errorln(v ...interface{}) {
	globalLogger.Errorln(v...)
}

/* Fatal */

// Fatal prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func Fatal(v ...interface{}) {
	globalLogger.Fatal(v...)
}

// Fatalf prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func Fatalf(format string, v ...interface{}) {
	globalLogger.Fatalf(format, v...)
}

// Fatalln prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func Fatalln(v ...interface{}) {
	globalLogger.Fatalln(v...)
}
