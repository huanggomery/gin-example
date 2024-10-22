package logging

import (
	"fmt"
	"gin-example/gin-blog/file"
	"gin-example/gin-blog/setting"
	"log"
	"os"
	"path"
	"time"
)

// 获取日志文件的文件夹路径 （例如 runtime/logs/）
func getLogDirPath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.LogSetting.LogSavePath)
}

// 获取日志文件的文件名 （例如 log20241022.log）
func getLogFileName() string {
	fileName := fmt.Sprintf(
		"%s%s.%s",
		setting.LogSetting.LogSaveName,
		time.Now().Format(setting.LogSetting.TimeFormat),
		setting.LogSetting.LogFileExt,
	)
	return fileName
}

// 打开日志文件，如果没有则创建文件并打开。发生异常会终止进程
func openLogFile(fileName, dirPath string) *os.File {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dirPath = path.Join(dir, dirPath)
	if !file.CheckPermission(dirPath) {
		log.Fatalf("Permission: %v", err)
	}
	if !file.CheckExist(dirPath) {
		if err = file.Mkdir(dirPath); err != nil {
			log.Fatalf("Mkdir: %v", err)
		}
	}

	filePath := path.Join(dirPath, fileName)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile: %v", err)
	}

	return file
}
