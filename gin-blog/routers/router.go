package routers

import (
	"gin-example/gin-blog/middleware"
	"gin-example/gin-blog/routers/api"
	v1 "gin-example/gin-blog/routers/api/v1"
	"gin-example/gin-blog/setting"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(setting.ServerSetting.RunMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/auth", api.GetAuth) // 提供用户名和密码，获取token

	apiv1 := r.Group("api/v1")
	apiv1.Use(middleware.JWT()) // 使用中间件校验token
	{
		// 获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		// 新建标签
		apiv1.POST("/tags", v1.AddTag)
		// 修改指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		// 删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		// 获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		// 获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		// 新增文章
		apiv1.POST("/articles", v1.AddArticle)
		// 修改指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		// 删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}
	return r
}
