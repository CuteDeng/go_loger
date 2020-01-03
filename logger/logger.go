package logger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

// LogLevel 日誌級別
type LogLevel uint16

// Logger ...
type Logger interface {
	Info(format string, a ...interface{})
	Trace(format string, a ...interface{})
	Debug(format string, a ...interface{})
	Warring(format string, a ...interface{})
	Error(format string, a ...interface{})
	Fatal(format string, a ...interface{})
}

// 日誌級別種類
const (
	UNKNOWN LogLevel = iota
	INFO
	TRACE
	DEBUG
	WARRING
	ERROR
	FATAL
)

func parseLevel(level string) (LogLevel, error) {
	level = strings.ToLower(level)
	switch level {
	case "info":
		return INFO, nil
	case "trace":
		return TRACE, nil
	case "debug":
		return DEBUG, nil
	case "warring":
		return WARRING, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		err := errors.New("未知等級")
		return UNKNOWN, err
	}
}

func levelName(level LogLevel) string {
	switch level {
	case INFO:
		return "INFO"
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case WARRING:
		return "WARRING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "DEBUG"
	}
}

func logDetail(skip int) (funcname, filename string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("error")
	}
	funcname = strings.Split(runtime.FuncForPC(pc).Name(), ".")[1]
	filename = path.Base(file)
	return
}
