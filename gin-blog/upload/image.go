package upload

import (
    "fmt"
    "gin-example/gin-blog/file"
    "gin-example/gin-blog/logging"
    "gin-example/gin-blog/setting"
    "gin-example/gin-blog/util"
    "mime/multipart"
    "os"
    "path"
    "strings"
)

// 获取图片完整的URL
func GetImageFullUrl(name string) string {
    return path.Join(
        setting.AppSetting.ImagePrefixUrl,
        setting.AppSetting.ImageSavePath,
        name,
    )
}

// 获取图片完整的本地路径
func GetImageFullDir() string {
    return path.Join(
        setting.AppSetting.RuntimeRootPath,
        setting.AppSetting.ImageSavePath,
    )
}

// 获取图片MD5编码后的名字
func GetImageName(name string) string {
    ext := file.GetExt(name)
    fileName := strings.TrimSuffix(name, ext)
    return util.EncodeMD5(fileName) + ext
}

// 检查图片的后缀名是否正确
func CheckImageExt(name string) bool {
    ext := file.GetExt(name)
    for _, allowed := range setting.AppSetting.ImageAllowExts {
        if strings.EqualFold(allowed, ext) {
            return true
        }
    }
    return false
}

// 检查图片大小是否超过限制
func CheckImageSize(img multipart.File) bool {
    sz, err := file.GetSize(img)
    if err != nil {
        logging.Warn(err)
        return false
    }
    return sz <= setting.AppSetting.ImageMaxSize
}

// 检查图片保存的目录是否存在，如果不存在则创建目录
func CheckAndCreateImageDir(dirPath string) (err error) {
    dir, err := os.Getwd()
    if err != nil {
        return
    }
    dirPath = path.Join(dir, dirPath)
    if !file.CheckExist(dirPath) {
        if err = file.Mkdir(dirPath); err != nil {
            return
        }
    }

    if !file.CheckPermission(dirPath) {
        err = fmt.Errorf("file.CheckPermission Permission denied src: %s", dirPath)
        return
    }
    return nil
}
