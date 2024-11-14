package v1

import (
    "gin-example/gin-blog/e"
    "gin-example/gin-blog/routers/api"
    "gin-example/gin-blog/service"
    "gin-example/gin-blog/setting"
    "gin-example/gin-blog/util"
    "net/http"

    "github.com/astaxie/beego/validation"
    "github.com/gin-gonic/gin"
    "github.com/unknwon/com"
)

// GetArticle 获取单个文章
func GetArticle(c *gin.Context) {
    id := com.StrTo(c.Param("id")).MustInt()

    valid := validation.Validation{}
    valid.Required(id, "id").Message("ID不能为空")
    valid.Min(id, 1, "id").Message("ID必须大于0")

    if valid.HasErrors() {
        api.LogErrors(valid.Errors)
        api.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
        return
    }

    articleService := service.ArticleService{ID: id}
    article, err := articleService.Get()
    if err != nil {
        api.Response(c, http.StatusOK, err.Code(), nil)
    }

    api.Response(c, http.StatusOK, e.SUCCESS, article)
}

// GetArticles 获取多个文章
func GetArticles(c *gin.Context) {
    valid := validation.Validation{}

    var state int = -1
    if arg := c.Query("state"); arg != "" {
        state = com.StrTo(arg).MustInt()
        valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
    }

    var tagId int = 0
    if arg := c.Query("tag_id"); arg != "" {
        tagId = com.StrTo(arg).MustInt()
        valid.Min(tagId, 1, "tag_id").Message("tag_id必须大于0")
    }

    createdBy := c.Query("created_by")
    if createdBy != "" {
        valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
    }

    if valid.HasErrors() {
        api.LogErrors(valid.Errors)
        api.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
        return
    }

    data := make(map[string]interface{})
    articleService := service.ArticleService{
        TagID:     tagId,
        State:     state,
        CreatedBy: createdBy,
        PageNum:   util.GetPage(c),
        PageSize:  setting.AppSetting.PageSize,
    }

    articles, err := articleService.GetAll()
    if err != nil {
        api.Response(c, http.StatusOK, err.Code(), nil)
        return
    }
    data["lists"] = articles
    data["count"] = articleService.Count()

    api.Response(c, http.StatusOK, e.SUCCESS, data)
}

// AddArticle 新增文章
func AddArticle(c *gin.Context) {
    valid := validation.Validation{}

    // 判断tag_id是否为空以及是否大于0
    tagIdArg := c.Query("tag_id")
    var tagId int
    valid.Required(tagIdArg, "tag_id").Message("tag_id不能为空")
    if tagIdArg != "" {
        tagId = com.StrTo(tagIdArg).MustInt()
        valid.Min(tagId, 1, "tag_id").Message("tag_id必须大于0")
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

    if valid.HasErrors() {
        api.LogErrors(valid.Errors)
        api.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
        return
    }

    articleService := service.ArticleService{
        TagID:     tagId,
        Title:     title,
        Desc:      desc,
        Content:   content,
        CreatedBy: createdBy,
        State:     state,
    }

    err := articleService.Add()
    if err != nil {
        api.Response(c, http.StatusOK, err.Code(), nil)
        return
    }

    api.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// EditArticle 修改文章
func EditArticle(c *gin.Context) {
    valid := validation.Validation{}

    // 获取修改参数
    id := com.StrTo(c.Param("id")).MustInt()
    var tagId int = 0
    if arg := c.Query("tag_id"); arg != "" {
        tagId = com.StrTo(arg).MustInt()
        valid.Min(tagId, 1, "tag_id").Message("tag_id必须大于0")
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

    if valid.HasErrors() {
        api.LogErrors(valid.Errors)
        api.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
        return
    }

    articleService := service.ArticleService{
        ID:         id,
        TagID:      tagId,
        Title:      title,
        Desc:       desc,
        Content:    content,
        ModifiedBy: modifiedBy,
        State:      state,
    }

    err := articleService.Edit()
    if err != nil {
        api.Response(c, http.StatusOK, err.Code(), nil)
        return
    }

    api.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// DeleteArticle 删除文章
func DeleteArticle(c *gin.Context) {
    id := com.StrTo(c.Param("id")).MustInt()
    valid := validation.Validation{}
    valid.Min(id, 1, "id").Message("ID必须大于0")

    if valid.HasErrors() {
        api.LogErrors(valid.Errors)
        api.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
        return
    }

    articleService := service.ArticleService{ID: id}

    err := articleService.Delete()
    if err != nil {
        api.Response(c, http.StatusOK, err.Code(), nil)
        return
    }

    api.Response(c, http.StatusOK, e.SUCCESS, nil)
}
