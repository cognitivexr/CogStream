package log

import "log"

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
		log.Printf("["+level.Prefix()+"] "+format+"\n", v)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	LogWithLevel(ERROR, format, v)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	LogWithLevel(WARN, format, v)
}

func (l *Logger) Info(format string, v ...interface{}) {
	LogWithLevel(INFO, format, v)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	LogWithLevel(DEBUG, format, v)
}

func (l *Logger) Trace(format string, v ...interface{}) {
	LogWithLevel(TRACE, format, v)
}

var RootLogger = &Logger{"root", INFO}

func LogWithLevel(level Level, format string, v ...interface{}) {
	RootLogger.Log(level, format, v)
}

func Error(format string, v ...interface{}) {
	LogWithLevel(ERROR, format, v)
}

func Warn(format string, v ...interface{}) {
	LogWithLevel(WARN, format, v)
}

func Info(format string, v ...interface{}) {
	LogWithLevel(INFO, format, v)
}

func Debug(format string, v ...interface{}) {
	LogWithLevel(DEBUG, format, v)
}

func Trace(format string, v ...interface{}) {
	LogWithLevel(TRACE, format, v)
}
