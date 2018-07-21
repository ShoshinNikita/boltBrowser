package log

import (
	"time"

	"github.com/fatih/color"
)

var (
	showTime   bool
	timeLayout = "01.02.2006 15:04:05"

	// For [ERR]
	red   = color.New(color.FgRed).PrintFunc()
	redf  = color.New(color.FgRed).PrintfFunc()
	redln = color.New(color.FgRed).PrintlnFunc()

	// For [INFO]
	cyan   = color.New(color.FgCyan).PrintFunc()
	cyanf  = color.New(color.FgCyan).PrintfFunc()
	cyanln = color.New(color.FgCyan).PrintlnFunc()

	// For time
	yellow   = color.New(color.FgYellow).PrintFunc()
	yellowf  = color.New(color.FgYellow).PrintfFunc()
	yellowln = color.New(color.FgYellow).PrintlnFunc()
)

// ShowTime enables showing of time
// Time isn't printed by default
func ShowTime() {
	showTime = true
}

// HideTime disable showing of time
// Time isn't printed by default
func HideTime() {
	showTime = false
}

func printTime() {
	yellowf("%s ", time.Now().Format(timeLayout))
}

func printErrMsg() {
	red("[ERR] ")
}

func printInfoMsg() {
	cyan("[INFO] ")
}
