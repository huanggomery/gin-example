package api

import (
    "gin-example/gin-blog/e"
    "gin-example/gin-blog/logging"
    "github.com/astaxie/beego/validation"
    "github.com/gin-gonic/gin"
)

func Response(c *gin.Context, httpCode int, errCode int, data interface{}) {
    c.JSON(httpCode, gin.H{
        "code": errCode,
        "msg":  e.GetMsg(errCode),
        "data": data,
    })
}

func LogErrors(errs []*validation.Error) {
    for _, err := range errs {
        logging.Info(err.Key, err.Message)
    }
}
