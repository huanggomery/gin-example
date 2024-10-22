package file

import (
	"os"
	"path"
)

// 活得文件内容字节数
func GetSize(filePath string) (int, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return int(info.Size()), nil
}

// 获取后缀名，如果没有的话返回空字符串
func GetExt(filePath string) string {
	return path.Ext(filePath)
}

// 判断是否存在
func CheckExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// 判断权限是否允许
func CheckPermission(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsPermission(err)
}

// 递归创建目录，权限为0777
func Mkdir(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	return err
}
