package logging

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-ini/ini"
)

type LogLevelType int

var (
	Cfg         *ini.File
	LogLevel    LogLevelType
	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
	logger      *log.Logger
)

// log level
const (
	DEBUG LogLevelType = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var levelString = []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}

func init() {
	var err error
	Cfg, err = ini.Load("gin-blog/conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'gin-blog/conf/app.ini': %v", err)
	}

	loadLog()

	filepath := getLogFileFullPath()
	file := openLogFile(filepath)
	logger = log.New(file, "", log.Ldate|log.Lmicroseconds)

}

// 从ini文件加载日志参数
func loadLog() {
	sec, err := Cfg.GetSection("log")
	if err != nil {
		log.Fatalf("Fail to get section 'log': %v", err)
	}

	level := sec.Key("LOG_LEVEL").MustString("DEBUG")
	switch level {
	case "DEBUG":
		LogLevel = DEBUG
	case "INFO":
		LogLevel = INFO
	case "WARNING":
		LogLevel = WARNING
	case "ERROR":
		LogLevel = ERROR
	case "FATAL":
		LogLevel = FATAL
	default:
		log.Fatalf("invalid log level: %s", level)
	}
	LogSavePath = sec.Key("LOG_SAVE_PATH").MustString("runtime/logs/")
	LogSaveName = sec.Key("LOG_SAVE_NAME").MustString("log")
	LogFileExt = sec.Key("LOG_FILE_EXT").MustString("log")
	TimeFormat = sec.Key("TIME_FORMAT").MustString("20200810")
}

// 检查日志级别并设置前缀
func setPrefix(level LogLevelType) bool {
	if level < LogLevel {
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