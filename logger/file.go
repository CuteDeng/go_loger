package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// FileLogger ...
type FileLogger struct {
	level       LogLevel
	filepath    string
	filename    string
	fileObj     *os.File
	errFileObj  *os.File
	maxFileSize int64
	logChan     chan *LogMsg
}

// LogMsg ...
type LogMsg struct {
	level     LogLevel
	message   string
	funcname  string
	filename  string
	timestamp string
	line      int
}

// NewFileLogger ...
func NewFileLogger(lv, fp, fn string, maxSize int64) *FileLogger {
	level, err := parseLevel(lv)
	if err != nil {
		panic(err)
	}
	fullFileName := path.Join(fp, fn)
	fullErrFileName := fullFileName + ".err"
	fo, err := os.OpenFile(fullFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println("open file err:", err)
	}
	efo, err := os.OpenFile(fullErrFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println("open errfile err:", err)
	}
	logc := make(chan *LogMsg, 5000)
	f1 := &FileLogger{
		level:       level,
		filepath:    fp,
		filename:    fn,
		fileObj:     fo,
		errFileObj:  efo,
		maxFileSize: maxSize,
		logChan:     logc,
	}
	go f1.writeTofile()
	return f1
}

func (log *FileLogger) writeTofile() {
	for {
		// 判断日志大小是否已经超过规定大小
		if log.checkFileSize(log.fileObj) {
			// 切分日志
			newFile, err := log.splitLogFile(log.fileObj)
			if err != nil {
				fmt.Println("log.splitLogFile ", err)
				return
			}
			log.fileObj = newFile
		}
		select {
		case logmsg := <-log.logChan:
			logInfo := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", logmsg.timestamp, levelName(logmsg.level), logmsg.filename, logmsg.funcname, logmsg.line, logmsg.message)
			fmt.Fprintf(log.fileObj, logInfo)
			if logmsg.level >= ERROR {
				// 判断日志大小是否已经超过规定大小
				if log.checkFileSize(log.errFileObj) {
					// 切分日志
					newErrFile, err := log.splitLogFile(log.errFileObj)
					if err != nil {
						fmt.Println("log.splitNewLogFile ", err)
						return
					}
					log.errFileObj = newErrFile
				}
				fmt.Fprintf(log.errFileObj, logInfo)
			}
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

// levelEnable ...
func (log *FileLogger) levelEnable(level LogLevel) bool {
	return level > log.level
}

// 判断日志大小
func (log *FileLogger) checkFileSize(file *os.File) bool {
	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println("file.Stat error ", err)
		return false
	}
	return fileinfo.Size() >= log.maxFileSize
}

// 切分日志
func (log *FileLogger) splitLogFile(file *os.File) (*os.File, error) {
	oldpath := path.Join(log.filepath, log.filename)
	now := time.Now().Format("20060102150405000")
	newpath := oldpath + ".bak" + now
	// 备份旧日志
	os.Rename(oldpath, newpath)
	// 创建新的日志文件
	newFile, err := os.OpenFile(oldpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	file.Close()
	return newFile, nil
}

func (log *FileLogger) writeLog(level LogLevel, format string, a ...interface{}) {
	if log.levelEnable(level) {
		msg := fmt.Sprintf(format, a...)
		funcname, filename, line := logDetail(1)
		// 将日志写到管道中
		logm := &LogMsg{
			level:     level,
			message:   msg,
			funcname:  funcname,
			filename:  filename,
			timestamp: time.Now().Format("2006-01-02 15:03:04"),
			line:      line,
		}
		select {
		case log.logChan <- logm:
		default:
			// 丢弃日志
		}

	}
}

// Info 消息
func (log *FileLogger) Info(format string, a ...interface{}) {
	log.writeLog(INFO, format, a...)
}

// Trace 消息
func (log *FileLogger) Trace(format string, a ...interface{}) {
	log.writeLog(TRACE, format, a...)
}

// Debug 消息
func (log *FileLogger) Debug(format string, a ...interface{}) {
	log.writeLog(DEBUG, format, a...)
}

// Warring 警告
func (log *FileLogger) Warring(format string, a ...interface{}) {
	log.writeLog(WARRING, format, a...)
}

// Error 消息
func (log *FileLogger) Error(format string, a ...interface{}) {
	log.writeLog(ERROR, format, a...)

}

// Fatal 消息
func (log *FileLogger) Fatal(format string, a ...interface{}) {
	log.writeLog(FATAL, format, a...)
}
