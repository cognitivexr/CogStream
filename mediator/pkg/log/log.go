package log

import (
	"log"
	"os"
)

type Level int

func (l Level) Prefix() string {
	switch l {
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	}
	return ""
}

var ERROR Level = 5
var WARN Level = 4
var INFO Level = 3
var DEBUG Level = 2
var TRACE Level = 1

type Logger struct {
	*log.Logger
	Name        string
	ActiveLevel Level
}

func (l *Logger) IsEnabledFor(level Level) bool {
	return level >= l.ActiveLevel
}

func (l *Logger) SetLevel(level Level) {
	l.ActiveLevel = level
}

func (l *Logger) Log(level Level, format string, v ...interface{}) {
	if l.IsEnabledFor(level) {
		log.Printf("["+level.Prefix()+"] "+format+"\n", v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	Log(ERROR, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	Log(WARN, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	Log(INFO, format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	Log(DEBUG, format, v...)
}

func (l *Logger) Trace(format string, v ...interface{}) {
	Log(TRACE, format, v...)
}

var RootLogger = &Logger{log.New(os.Stderr, "", log.LstdFlags), "root", INFO}

func Log(level Level, format string, v ...interface{}) {
	RootLogger.Log(level, format, v...)
}

func Error(format string, v ...interface{}) {
	Log(ERROR, format, v...)
}

func Warn(format string, v ...interface{}) {
	Log(WARN, format, v...)
}

func Info(format string, v ...interface{}) {
	Log(INFO, format, v...)
}

func Debug(format string, v ...interface{}) {
	Log(DEBUG, format, v...)
}

func Trace(format string, v ...interface{}) {
	Log(TRACE, format, v...)
}

// functions from log/log

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	log.Print(v...)
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	log.Println(v...)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	log.Panic(v...)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	log.Panicln(v...)
}
