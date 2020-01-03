package logger

import (
	"fmt"
	"time"
)

// ConsoleLogger 结构体xxx
type ConsoleLogger struct {
	level LogLevel
}

// NewConsoleLogger 构造函数
func NewConsoleLogger(level string) *ConsoleLogger {
	l, err := parseLevel(level)
	if err != nil {
		panic(err)
	}
	return &ConsoleLogger{
		level: l,
	}
}

func (log *ConsoleLogger) levelEnable(level LogLevel) bool {
	return level > log.level
}

func (log *ConsoleLogger) logInfo(level LogLevel, format string, a ...interface{}) {
	if log.levelEnable(level) {
		msg := fmt.Sprintf(format, a...)
		new := time.Now()
		funcname, filename, line := logDetail(1)
		fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", new.Format("2006-01-02 15:03:04"), levelName(level), filename, funcname, line, msg)
	}
}

// Info 消息
func (log *ConsoleLogger) Info(format string, a ...interface{}) {
	log.logInfo(INFO, format, a...)
}

// Trace 消息
func (log *ConsoleLogger) Trace(format string, a ...interface{}) {
	log.logInfo(TRACE, format, a...)
}

// Debug 消息
func (log *ConsoleLogger) Debug(format string, a ...interface{}) {
	log.logInfo(DEBUG, format, a...)
}

// Warring 警告
func (log *ConsoleLogger) Warring(format string, a ...interface{}) {
	log.logInfo(WARRING, format, a...)
}

// Error 消息
func (log *ConsoleLogger) Error(format string, a ...interface{}) {
	log.logInfo(ERROR, format, a...)

}

// Fatal 消息
func (log *ConsoleLogger) Fatal(format string, a ...interface{}) {
	log.logInfo(FATAL, format, a...)
}
