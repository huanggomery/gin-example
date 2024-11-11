package service

import (
    "encoding/json"
    "errors"
    "gin-example/gin-blog/e"
    "gin-example/gin-blog/logging"
    "gin-example/gin-blog/models"
    "gin-example/gredis"
    "strconv"
    "strings"
    "time"
)

type TagService struct {
    ID         int
    Name       string
    CreatedBy  string
    ModifiedBy string
    State      int

    PageNum  int
    PageSize int
}

// getKey 获取访问缓存的key
func (t *TagService) getKey() string {
    keys := []string{e.CACHE_TAG}

    if t.Name != "" {
        keys = append(keys, t.Name)
    }
    if t.State == 0 {
        keys = append(keys, "off")
    } else if t.State == 1 {
        keys = append(keys, "on")
    }
    if t.PageNum >= 0 {
        keys = append(keys, strconv.Itoa(t.PageNum))
    }
    if t.PageSize >= 0 {
        keys = append(keys, strconv.Itoa(t.PageSize))
    }

    return strings.Join(keys, "_")
}

// getMaps 获取映射字典
func (t *TagService) getMaps() map[string]interface{} {
    maps := make(map[string]interface{})

    if t.ID > 0 {
        maps["id"] = t.ID
    }
    if t.Name != "" {
        maps["name"] = t.Name
    }
    if t.CreatedBy != "" {
        maps["created_by"] = t.CreatedBy
    }
    if t.ModifiedBy != "" {
        maps["modified_by"] = t.ModifiedBy
    }
    if t.State >= 0 {
        maps["state"] = t.State
    }

    return maps
}

// getFromCache 从缓存获取数据
func (t *TagService) getFromCache() ([]models.Tag, error) {
    key := t.getKey()

    if !gredis.Exists(key) {
        return nil, errors.New("cache not exists")
    }

    data, err := gredis.Get(key)
    if err != nil {
        logging.Error(err)
        return nil, err
    }

    tags := make([]models.Tag, 0)
    err = json.Unmarshal([]byte(data), &tags)
    if err != nil {
        logging.Error(err)
        return nil, err
    }

    return tags, nil
}

func (t *TagService) ExistByName() bool {
    return models.ExistTagByName(t.Name)
}

func (t *TagService) ExistById() bool {
    return models.ExistTagByID(t.ID)
}

// GetAll 获得所有符合条件的标签
func (t *TagService) GetAll() ([]models.Tag, error) {
    tags := make([]models.Tag, 0)

    // 尝试从缓存获取数据
    tags, err := t.getFromCache()
    if err == nil {
        return tags, nil
    }

    // 从数据库中获取
    tags = models.GetTags(t.PageNum, t.PageSize, t.getMaps())

    // 设置缓存
    if err = gredis.Set(t.getKey(), tags, time.Second*3600); err != nil {
        logging.Error(err)
    }

    return tags, nil
}

func (t *TagService) Count() int {
    return models.GetTagCount(t.getMaps())
}

// Add 添加新标签
func (t *TagService) Add() error {
    if !models.AddTag(t.Name, t.State, t.CreatedBy) {
        return errors.New("add tag failed")
    }

    // 删除缓存
    if t.State == 0 {
        gredis.LikeDel("off")
    } else if t.State == 1 {
        gredis.LikeDel("on")
    }
    gredis.LikeDel(t.Name)

    return nil
}

// Edit 修改标签
func (t *TagService) Edit() error {
    tag, _ := models.GetTagByID(t.ID)

    data := map[string]interface{}{
        "modified_by": t.ModifiedBy,
    }

    if t.Name != "" {
        data["name"] = t.Name
    }
    if t.State >= 0 {
        data["state"] = t.State
    }

    if !models.EditTag(t.ID, data) {
        return errors.New("edit tag failed")
    }

    // 删除缓存
    if t.State == 0 {
        gredis.LikeDel("off")
    } else if t.State == 1 {
        gredis.LikeDel("on")
    }
    gredis.LikeDel(tag.Name)

    return nil
}

func (t *TagService) Delete() error {
    tag, _ := models.GetTagByID(t.ID)

    if !models.DeleteTag(t.ID) {
        return errors.New("delete tag failed")
    }

    // 删除缓存
    if tag.State == 0 {
        gredis.LikeDel("off")
    } else if tag.State == 1 {
        gredis.LikeDel("on")
    }
    gredis.LikeDel(tag.Name)

    return nil
}
