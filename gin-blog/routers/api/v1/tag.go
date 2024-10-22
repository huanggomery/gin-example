package v1

import (
	"gin-example/gin-blog/e"
	"gin-example/gin-blog/logging"
	"gin-example/gin-blog/models"
	"gin-example/gin-blog/setting"
	"gin-example/gin-blog/util"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	if arg := c.Query("state"); arg != "" {
		state := com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	data["lists"] = models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["count"] = models.GetTagCount(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": data,
	})
}

// 新增文章标签
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.SUCCESS
	if valid.HasErrors() {
		// 参数错误
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
		code = e.INVALID_PARAMS
	} else if models.ExistTagByName(name) {
		// 参数正确，但标签已存在，无法新增
		code = e.ERROR_EXIST_TAG
	}

	if code == e.SUCCESS {
		// 只有一切正常，才向表中新增标签
		models.AddTag(name, state, createdBy)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// 修改文章标签
func EditTag(c *gin.Context) {
	// 获取修改参数
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if state_str := c.Query("state"); state_str != "" {
		state = com.StrTo(state_str).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	// 校验参数合法性
	valid.Required(id, "id").Message("ID不能为空")
	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	var code int
	if !valid.HasErrors() {
		if models.ExistTagByID(id) {
			// 标签存在，可以修改
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			models.EditTag(id, data)
			code = e.SUCCESS
		} else {
			// 标签不存在
			code = e.ERROR_NOT_EXIST_TAG
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

// 删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	var code int
	if !valid.HasErrors() {
		if models.ExistTagByID(id) {
			// 标签存在，可以删除
			models.DeleteTag(id)
			code = e.SUCCESS
		} else {
			// 标签不存在
			code = e.ERROR_NOT_EXIST_TAG
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
