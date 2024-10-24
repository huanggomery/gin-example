package api

import (
	"gin-example/gin-blog/e"
	"gin-example/gin-blog/upload"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

// 从POST表单获取图片，校验图片，并保存图片
func UploadImage(c *gin.Context) {
	img, img_header, err := c.Request.FormFile("image")
	data := map[string]string{}
	code := e.SUCCESS

	if err != nil {
		code = e.ERROR
	}
	if img_header == nil {
		code = e.INVALID_PARAMS
	}
	if code != e.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
		return
	}

	imgName := upload.GetImageName(img_header.Filename)
	dirPath := upload.GetImageFullDir()
	imgPath := path.Join(dirPath, imgName)

	if upload.CheckImageExt(imgName) && upload.CheckImageSize(img) {
		if err = upload.CheckAndCreateImageDir(dirPath); err != nil {
			// 创建目录失败或权限不允许
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
		} else if err = c.SaveUploadedFile(img_header, imgPath); err != nil {
			// 保存图片失败
			code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
		} else {
			// 一切正常
			data["image_url"] = upload.GetImageFullUrl(imgName)
		}
	} else {
		code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
