package v1

import (
	"gin-example/gin-blog/e"
	"gin-example/gin-blog/models"
	"gin-example/gin-blog/setting"
	"gin-example/gin-blog/util"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取单个文章
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("ID不能为空")
	valid.Min(id, 1, "id").Message("ID必须大于0")

	var code int
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			// 文章存在，可以获取
			data = models.GetArticle(id)
			code = e.SUCCESS
		} else {
			// 文章不存在
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		// 参数错误
		code = e.INVALID_PARAMS
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// 获取多个文章
func GetArticles(c *gin.Context) {
	valid := validation.Validation{}
	maps := make(map[string]interface{})

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tag_id int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tag_id = com.StrTo(arg).MustInt()
		maps["tag_id"] = tag_id
		valid.Min(tag_id, 1, "tag_id").Message("tag_id必须大于0")
	}

	created_by := c.Query("created_by")
	if created_by != "" {
		maps["created_by"] = created_by
		valid.MaxSize(created_by, 100, "created_by").Message("创建人最长为100字符")
	}

	var code int
	data := make(map[string]interface{})

	if !valid.HasErrors() {
		// 参数正确，正常查询
		data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
		data["count"] = models.GetArticleCount(maps)
		code = e.SUCCESS
	} else {
		// 参数错误
		code = e.INVALID_PARAMS
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// 新增文章
func AddArticle(c *gin.Context) {
	valid := validation.Validation{}

	// 判断tag_id是否为空以及是否大于0
	tag_id_arg := c.Query("tag_id")
	var tag_id int
	valid.Required(tag_id_arg, "tag_id").Message("tag_id不能为空")
	if tag_id_arg != "" {
		tag_id = com.StrTo(tag_id_arg).MustInt()
		valid.Min(tag_id, 1, "tag_id").Message("tag_id必须大于0")
	}

	// 文章参数
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()

	// 校验参数合法性
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	var code int
	if !valid.HasErrors() {
		// 参数正确，可以新增文章
		models.AddArticle(map[string]interface{}{
			"tag_id":     tag_id,
			"title":      title,
			"desc":       desc,
			"content":    content,
			"created_by": createdBy,
			"state":      state,
		})
		code = e.SUCCESS
	} else {
		// 参数错误
		code = e.INVALID_PARAMS
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// 修改文章
func EditArticle(c *gin.Context) {
	valid := validation.Validation{}

	// 获取修改参数
	id := com.StrTo(c.Param("id")).MustInt()
	var tag_id int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tag_id = com.StrTo(arg).MustInt()
		valid.Min(tag_id, 1, "tag_id").Message("tag_id必须大于0")
	}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")

	// 校验参数合法性
	valid.Required(id, "id").Message("ID不能为空")
	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	var code int

	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			// 文章存在，可以修改
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if tag_id > 0 {
				data["tag_id"] = tag_id
			}
			if state != -1 {
				data["state"] = state
			}
			if title != "" {
				data["title"] = title
			}
			if desc != "" {
				data["desc"] = desc
			}
			if content != "" {
				data["content"] = content
			}
			models.EditArticle(id, data)
			code = e.SUCCESS
		} else {
			// 文章不存在
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		// 参数错误
		code = e.INVALID_PARAMS
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	var code int
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			// 文章存在，可以删除
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			// 文章不存在
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		// 参数错误
		code = e.INVALID_PARAMS
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
