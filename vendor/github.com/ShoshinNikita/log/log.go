// Package log provides functions for pretty print
//
// Patterns of functions print:
// * Print(), Printf(), Println():
//   (?time) msg
// * Info(), Infof(), Infoln():
//   (?time) [INFO] msg
// * Warn(), Warnf(), Warnln():
//   (?time) [WARN] warning
// * Error(), Errorf(), Errorln():
//   (?time) [ERR] (?file:line) error
// * Fatal(), Fatalf(), Fatalln():
//   (?time) [FATAL] (?file:line) error
//
// Time pattern: MM.dd.yyyy hh:mm:ss (01.30.2018 05:5:59)
//
package log

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

const (
	DefaultTimeLayout = "01.02.2006 15:04:05"
)

type textStruct struct {
	text string
	ch   chan struct{}
}

func newText(text string) textStruct {
	return textStruct{text: text, ch: make(chan struct{})}
}

func (t *textStruct) done() {
	close(t.ch)
}

type Logger struct {
	printTime      bool
	printColor     bool
	printErrorLine bool

	printChan chan textStruct
	global    bool

	output     io.Writer
	timeLayout string
}

// NewLogger creates *Logger and run goroutine (Logger.printer())
func NewLogger() *Logger {
	l := new(Logger)
	l.output = color.Output
	l.timeLayout = DefaultTimeLayout
	l.printChan = make(chan textStruct, 200)
	go l.printer()
	return l
}

func (l *Logger) printer() {
	for text := range l.printChan {
		fmt.Fprint(l.output, text.text)
		text.done()
	}
}

func (l *Logger) printText(text string) {
	t := newText(text)
	l.printChan <- t
	<-t.ch
}

// PrintTime sets Logger.printTime to b
func (l *Logger) PrintTime(b bool) {
	l.printTime = b
}

// PrintColor sets Logger.printColor to b
func (l *Logger) PrintColor(b bool) {
	l.printColor = b
}

// PrintErrorLine sets Logger.printErrorLine to b
func (l *Logger) PrintErrorLine(b bool) {
	l.printErrorLine = b
}

// ChangeOutput changes Logger.output writer.
// Default Logger.output is github.com/fatih/color.Output
func (l *Logger) ChangeOutput(w io.Writer) {
	l.output = w
}

// ChangeTimeLayout changes Logger.timeLayout
// Default Logger.timeLayout is DefaultTimeLayout
func (l *Logger) ChangeTimeLayout(layout string) {
	l.timeLayout = layout
}
