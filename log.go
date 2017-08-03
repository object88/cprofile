package cprofile

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// LogLevel is the logging level: None, Error, Warn, Info, Verbose, or Debug
type LogLevel int

const (
	// None means that the log should never write
	None LogLevel = iota

	// Error means that only errors will be written
	Error

	// Warn means that errors and warnings will be written
	Warn
	Info
	Verbose
	Debug

	stdOutLogname = "__stdout"
	stdErrLogname = "__stderr"
)

var m sync.RWMutex
var ls = map[string]*Log{}

// Log is a fairly basic logger
type Log struct {
	w   io.Writer
	lvl LogLevel
}

func NewLog(name string, w io.Writer) *Log {
	return getLog(name, w)
}

func Stderr() *Log {
	return getLog(stdErrLogname, os.Stderr)
}

func Stdout() *Log {
	return getLog(stdOutLogname, os.Stdout)
}

func getLog(name string, w io.Writer) *Log {
	m.RLock()

	if l, ok := ls[name]; ok {
		m.RUnlock()
		return l
	}

	m.RUnlock()

	m.Lock()

	if l, ok := ls[name]; ok {
		m.Unlock()
		return l
	}

	l := &Log{w, Error}
	ls[name] = l

	m.Unlock()
	return l
}

// Printf will always write out.  If the pointer receiver is nil,
// the log for `os.Stdout` will be used.
func (l *Log) Printf(msg string, v ...interface{}) {
	if l == nil {
		l = Stdout()
	}

	l.write(msg, v...)
}

// SetLevel will adjust the logger's level.  If the pointer receiver is nil,
// the log for `os.Stdout` will be used.
func (l *Log) SetLevel(lvl LogLevel) {
	if l == nil {
		l = Stdout()
	}

	l.lvl = lvl
}

// Verbosef will write if the log level is at least Verbose
func (l *Log) Verbosef(msg string, v ...interface{}) {
	if l == nil {
		l = Stdout()
	}

	if l.lvl < Verbose {
		return
	}

	l.write(msg, v...)
}

func (l *Log) write(msg string, v ...interface{}) {
	if v == nil {
		l.w.Write([]byte(msg))
	} else {
		m := fmt.Sprintf(msg, v...)
		l.w.Write([]byte(m))
	}
}
