package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

func getLogFileFullPath() string {
	prefixName := LogSavePath
	suffixName := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)
	return fmt.Sprintf("%s%s", prefixName, suffixName)
}

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission: %v\n", err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile: %v\n", err)
	}

	return file
}

func mkDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(path.Join(dir, LogSavePath), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
