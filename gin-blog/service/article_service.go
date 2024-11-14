package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "gin-example/gin-blog/e"
    "gin-example/gin-blog/logging"
    "gin-example/gin-blog/models"
    "gin-example/gredis"
    "strings"
    "time"
)

type ArticleService struct {
    ID         int
    TagID      int
    Title      string
    Desc       string
    Content    string
    CreatedBy  string
    ModifiedBy string
    State      int

    PageNum  int
    PageSize int
}

func (article *ArticleService) getKey() string {
    keys := []string{e.CACHE_ARTICLE}

    if article.ID != 0 {
        keys = append(keys, fmt.Sprintf("id%d", article.ID))
    }
    if article.TagID != 0 {
        keys = append(keys, fmt.Sprintf("tag%d", article.TagID))
    }
    if article.CreatedBy != "" {
        keys = append(keys, article.CreatedBy)
    }
    if article.State == 0 {
        keys = append(keys, "off")
    } else if article.State == 1 {
        keys = append(keys, "on")
    }
    if article.PageNum > 0 {
        keys = append(keys, fmt.Sprintf("from%d", article.PageNum))
    }
    if article.PageSize > 0 {
        keys = append(keys, fmt.Sprintf("size%d", article.PageSize))
    }

    return strings.Join(keys, "_")
}

func (article *ArticleService) getMaps() map[string]interface{} {
    maps := make(map[string]interface{})
    if article.TagID != 0 {
        maps["tag_id"] = article.TagID
    }
    if article.Title != "" {
        maps["title"] = article.Title
    }
    if article.Desc != "" {
        maps["desc"] = article.Desc
    }
    if article.Content != "" {
        maps["content"] = article.Content
    }
    if article.CreatedBy != "" {
        maps["created_by"] = article.CreatedBy
    }
    if article.ModifiedBy != "" {
        maps["modified_by"] = article.ModifiedBy
    }
    if article.State >= 0 {
        maps["state"] = article.State
    }

    return maps
}

// getFromCache 从缓存获取数据
func (article *ArticleService) getFromCache() ([]models.Article, error) {
    key := article.getKey()

    if !gredis.Exists(key) {
        return nil, errors.New("cache not exists")
    }

    data, err := gredis.Get(key)
    if err != nil {
        logging.Error(err)
        return nil, err
    }

    var articles []models.Article
    if data == "" {
        return articles, nil
    }
    if err = json.Unmarshal([]byte(data), &articles); err != nil {
        logging.Error(err)
        return nil, err
    }

    return articles, nil
}

func (article *ArticleService) ExistByID() bool {
    // 查看是否缓存
    if data, err := article.getFromCache(); err == nil {
        // 缓存中存在该key，判断value是否为空
        return len(data) > 0
    }

    // 缓存中不存在该key，查找数据库
    if !models.ExistArticleByID(article.ID) {
        // 数据库中也不存在，在缓存中设置空值
        if err := gredis.Set(article.getKey(), "", time.Second*3600); err != nil {
            logging.Error(err)
        }
        return false
    }

    // 数据库中存在，设置缓存
    data := models.GetArticle(article.ID)
    if err := gredis.Set(article.getKey(), data, time.Second*3600); err != nil {
        logging.Error(err)
    }
    return true
}

// Get 根据ID获取单个文章
func (article *ArticleService) Get() (models.Article, *e.Error) {
    // 判断是否存在，存在的话加载到缓存
    if !article.ExistByID() {
        return models.Article{}, e.NewError(e.ERROR_NOT_EXIST_ARTICLE, "article not exists")
    }

    // 尝试从缓存获取数据
    if data, err := article.getFromCache(); err == nil {
        return data[0], nil
    }

    // 从数据库中获取， 应该不会走以下逻辑
    data := models.GetArticle(article.ID)

    // 设置缓存
    if err := gredis.Set(article.getKey(), data, time.Second*3600); err != nil {
        logging.Error(err)
    }

    return data, nil
}

// GetAll 获取多篇文章
func (article *ArticleService) GetAll() ([]models.Article, *e.Error) {
    // 尝试从缓存获取数据
    if data, err := article.getFromCache(); err == nil {
        if len(data) == 0 {
            return []models.Article{}, e.NewError(e.ERROR_NOT_EXIST_ARTICLE, "can't find any article")
        }
        return data, nil
    }

    // 从数据库中获取
    data := models.GetArticles(article.PageNum, article.PageSize, article.getMaps())

    // 设置缓存
    if err := gredis.Set(article.getKey(), data, time.Second*3600); err != nil {
        logging.Error(err)
    }
    if len(data) == 0 {
        return []models.Article{}, e.NewError(e.ERROR_NOT_EXIST_ARTICLE, "can't find any article")
    }

    return data, nil
}

func (article *ArticleService) Count() int {
    return models.GetArticleCount(article.getMaps())
}

// Add 添加新文章
func (article *ArticleService) Add() *e.Error {
    if !models.AddArticle(article.getMaps()) {
        return e.NewError(e.ERROR, "add article fail")
    }

    // 删除缓存
    if article.State == 0 {
        if err := gredis.LikeDel(e.CACHE_ARTICLE, "off"); err != nil {
            logging.Error(err)
        }
    } else if article.State == 1 {
        if err := gredis.LikeDel(e.CACHE_ARTICLE, "on"); err != nil {
            logging.Error(err)
        }
    }
    if err := gredis.LikeDel(e.CACHE_ARTICLE, fmt.Sprintf("tag%d", article.TagID)); err != nil {
        logging.Error(err)
    }
    if err := gredis.LikeDel(e.CACHE_ARTICLE, article.CreatedBy); err != nil {
        logging.Error(err)
    }

    return nil
}

// Edit 编辑文章
func (article *ArticleService) Edit() *e.Error {
    if !article.ExistByID() {
        return e.NewError(e.ERROR_NOT_EXIST_ARTICLE, "article not exists")
    }

    if !models.EditArticle(article.ID, article.getMaps()) {
        return e.NewError(e.ERROR, "edit article fail")
    }

    // 删除缓存
    if article.State == 0 {
        if err := gredis.LikeDel(e.CACHE_ARTICLE, "off"); err != nil {
            logging.Error(err)
        }
    } else if article.State == 1 {
        if err := gredis.LikeDel(e.CACHE_ARTICLE, "on"); err != nil {
            logging.Error(err)
        }
    }
    if err := gredis.LikeDel(e.CACHE_ARTICLE, fmt.Sprintf("tag%d", article.TagID)); err != nil {
        logging.Error(err)
    }

    return nil
}

// Delete 删除文章
func (article *ArticleService) Delete() *e.Error {
    if !article.ExistByID() {
        return e.NewError(e.ERROR_NOT_EXIST_ARTICLE, "article not exists")
    }

    data, _ := article.Get()

    if !models.DeleteArticle(article.ID) {
        return e.NewError(e.ERROR, "delete article fail")
    }

    // 删除缓存
    if data.State == 0 {
        if err := gredis.LikeDel(e.CACHE_ARTICLE, "off"); err != nil {
            logging.Error(err)
        }
    } else if data.State == 1 {
        if err := gredis.LikeDel(e.CACHE_ARTICLE, "on"); err != nil {
            logging.Error(err)
        }
    }
    if err := gredis.LikeDel(e.CACHE_ARTICLE, fmt.Sprintf("id%d", article.ID)); err != nil {
        logging.Error(err)
    }
    if err := gredis.LikeDel(e.CACHE_ARTICLE, fmt.Sprintf("tag%d", article.TagID)); err != nil {
        logging.Error(err)
    }
    if err := gredis.LikeDel(e.CACHE_ARTICLE, article.CreatedBy); err != nil {
        logging.Error(err)
    }

    return nil
}
