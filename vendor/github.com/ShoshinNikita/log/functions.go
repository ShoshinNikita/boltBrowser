package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	usualINFO  = "[INFO]  "
	usualWARN  = "[WARN]  "
	usualERR   = "[ERR]   "
	usualFATAL = "[FATAL] "
)

var (
	timePrintf   = color.New(color.FgHiGreen).SprintfFunc()
	callerPrintf = color.RedString // color is the same as coloredErr

	coloredINFO  = color.CyanString(usualINFO)
	coloredWARN  = color.YellowString(usualWARN)
	coloredERR   = color.RedString(usualERR)
	coloredFATAL = color.New(color.BgRed).Sprint("[FATAL]") + " "
)

// getTime returns "file:line" if l.printErrorLine == true, else it returns empty string
func (l Logger) getCaller() string {

	if !l.printErrorLine {
		return ""
	}

	var (
		file string
		line int
		ok   bool
	)

	if l.global {
		_, file, line, ok = runtime.Caller(5)
	} else {
		_, file, line, ok = runtime.Caller(4)
	}
	if !ok {
		return ""
	}

	var shortFile string
	for i := len(file) - 1; i >= 0; i-- {
		if file[i] == '/' {
			shortFile = file[i+1:]
			break
		}
	}

	if l.printColor {
		return callerPrintf("%s:%d ", shortFile, line)
	}
	return fmt.Sprintf("%s:%d ", shortFile, line)
}

// getTime returns time if l.printTime == true, else it returns empty string
func (l Logger) getTime() string {
	if !l.printTime {
		return ""
	}

	if l.printColor {
		return timePrintf("%s ", time.Now().Format(l.timeLayout))
	}
	return fmt.Sprintf("%s ", time.Now().Format(l.timeLayout))
}

func (l Logger) getInfoMsg() string {
	if l.printColor {
		return coloredINFO
	}
	return usualINFO
}

func (l Logger) getWarnMsg() string {
	if l.printColor {
		return coloredWARN
	}
	return usualWARN
}

func (l Logger) getErrMsg() string {
	if l.printColor {
		return coloredERR
	}
	return usualERR
}

func (l Logger) getFatalMsg() (s string) {
	if l.printColor {
		return coloredFATAL
	}
	return usualFATAL
}

type prefixFunc func() string

// addPrefixes adds prefixes. It uses strings.Builder
func addPrefixes(str string, prefixes ...prefixFunc) string {
	b := strings.Builder{}

	for _, f := range prefixes {
		b.WriteString(f())
	}
	b.WriteString(str)

	return b.String()
}
