package logging

import (
	"fmt"
	"gin-example/gin-blog/setting"
	"log"
	"runtime"
)

var (
	logger   *log.Logger
	loglevel LogLevelType
)

type LogLevelType int

// log level
const (
	DEBUG LogLevelType = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var levelString = []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}

func Setup() {
	// 打开日志文件，初始化日志句柄
	fileName := getLogFileName()
	dirPath := getLogDirPath()
	file := openLogFile(fileName, dirPath)
	logger = log.New(file, "", log.Ldate|log.Lmicroseconds)

	// 设置日志级别
	switch setting.LogSetting.LogLeval {
	case "DEBUG":
		loglevel = DEBUG
	case "INFO":
		loglevel = INFO
	case "WARNING":
		loglevel = WARNING
	case "ERROR":
		loglevel = ERROR
	case "FATAL":
		loglevel = FATAL
	default:
		log.Fatalf("invalid log level: %s", setting.LogSetting.LogLeval)
	}
}

// 检查日志级别并设置前缀
func setPrefix(level LogLevelType) bool {
	if level < loglevel {
		return false
	}

	_, file, line, ok := runtime.Caller(2) // 获取调用位置信息
	if ok {
		logger.SetPrefix(fmt.Sprintf("[%s][%s:%d]", levelString[level], file, line))
	} else {
		logger.SetPrefix(fmt.Sprintf("[%s]", levelString[level]))
	}
	return true
}

func Debug(v ...any) {
	if setPrefix(DEBUG) {
		logger.Println(v...)
	}
}

func Debugf(format string, v ...any) {
	if setPrefix(DEBUG) {
		logger.Printf(format, v...)
	}
}

func Info(v ...any) {
	if setPrefix(INFO) {
		logger.Println(v...)
	}
}

func Infof(format string, v ...any) {
	if setPrefix(INFO) {
		logger.Printf(format, v...)
	}
}

func Warn(v ...any) {
	if setPrefix(WARNING) {
		logger.Println(v...)
	}
}

func Warnf(format string, v ...any) {
	if setPrefix(WARNING) {
		logger.Printf(format, v...)
	}
}

func Error(v ...any) {
	if setPrefix(ERROR) {
		logger.Println(v...)
	}
}

func Errorf(format string, v ...any) {
	if setPrefix(ERROR) {
		logger.Printf(format, v...)
	}
}

func Fatal(v ...any) {
	if setPrefix(FATAL) {
		logger.Fatalln(v...)
	}
}

func Fatalf(format string, v ...any) {
	if setPrefix(FATAL) {
		logger.Fatalf(format, v...)
	}
}
