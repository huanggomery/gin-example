package v1

import (
    "gin-example/gin-blog/e"
    "gin-example/gin-blog/logging"
    "gin-example/gin-blog/models"
    "gin-example/gin-blog/routers/api"
    "gin-example/gin-blog/service"
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
    state := -1

    if arg := c.Query("state"); arg != "" {
        state = com.StrTo(arg).MustInt()
    }

    tagService := service.TagService{
        Name:     name,
        State:    state,
        PageNum:  util.GetPage(c),
        PageSize: setting.AppSetting.PageSize,
    }

    tags, err := tagService.GetAll()
    if err != nil {
        api.Response(c, http.StatusInternalServerError, e.ERROR, nil)
        return
    }
    count := tagService.Count()

    data := map[string]interface{}{
        "list":  tags,
        "count": count,
    }

    api.Response(c, http.StatusOK, e.SUCCESS, data)
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

    if valid.HasErrors() {
        // 参数错误
        api.LogErrors(valid.Errors)
        api.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
        return
    }
    if models.ExistTagByName(name) {
        // 参数正确，但标签已存在，无法新增
        api.Response(c, http.StatusOK, e.ERROR_EXIST_TAG, nil)
        return
    }

    // 只有一切正常，才向表中新增标签
    tagService := service.TagService{
        Name:      name,
        State:     state,
        CreatedBy: createdBy,
    }
    if err := tagService.Add(); err != nil {
        logging.Error(err)
        api.Response(c, http.StatusInternalServerError, e.ERROR, nil)
        return
    }

    api.Response(c, http.StatusOK, e.SUCCESS, nil)
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

    // 参数错误
    if valid.HasErrors() {
        api.LogErrors(valid.Errors)
        api.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
        return
    }

    tagService := service.TagService{
        ID:         id,
        Name:       name,
        State:      state,
        ModifiedBy: modifiedBy,
    }

    // 标签不存在
    if !tagService.ExistById() {
        api.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
        return
    }

    // 标签存在，可以修改
    if err := tagService.Edit(); err != nil {
        logging.Error(err)
        api.Response(c, http.StatusInternalServerError, e.ERROR, nil)
        return
    }

    api.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
    id := com.StrTo(c.Param("id")).MustInt()
    valid := validation.Validation{}
    valid.Min(id, 1, "id").Message("ID必须大于0")

    if valid.HasErrors() {
        api.LogErrors(valid.Errors)
        api.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
        return
    }

    tagService := service.TagService{
        ID: id,
    }

    // 标签不存在
    if !tagService.ExistById() {
        api.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
        return
    }

    // 标签存在，可以删除
    if err := tagService.Delete(); err != nil {
        logging.Error(err)
        api.Response(c, http.StatusInternalServerError, e.ERROR, nil)
        return
    }

    api.Response(c, http.StatusOK, e.SUCCESS, nil)
}
