package api

import (
	"gin-example/gin-blog/e"
	"gin-example/gin-blog/logging"
	"gin-example/gin-blog/models"
	"gin-example/gin-blog/util"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 校验参数
	valid := validation.Validation{}
	ok, _ := valid.Valid(&auth{username, password})

	code := e.SUCCESS
	data := make(map[string]interface{})

	if ok {
		if models.CheckAuth(username, password) {
			tokenString, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = tokenString
				code = e.SUCCESS
			}
		} else {
			// 用户名或密码错误
			code = e.ERROR_AUTH
		}
	} else {
		// 参数错误
		code = e.INVALID_PARAMS

		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
