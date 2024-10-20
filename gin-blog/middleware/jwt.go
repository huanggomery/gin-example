package middleware

import (
	"gin-example/gin-blog/e"
	"gin-example/gin-blog/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 校验token的中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")

		code := e.SUCCESS

		if tokenString != "" {
			claims, err := util.ParseToken(tokenString)
			if err != nil {
				// 解析验证失败
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if claims.ExpiresAt <= time.Now().Unix() {
				// token过期
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		} else {
			// 请求中未包含token
			code = e.INVALID_PARAMS
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": map[string]string{},
			})

			// 出现错误，返回数据并结束
			c.Abort()
			return
		}

		// 一切正常，继续往后执行
		c.Next()
	}
}
