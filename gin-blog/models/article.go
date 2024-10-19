package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title      string `json:"title"`
	Desc       string `json:"describe"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 结构体 Article -> 表 blog_article
func (Article) TableName() string {
	return "blog_article"
}

// 自动添加创建时间
func (article *Article) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("created_on", int(time.Now().Unix()))
	return nil
}

// 自动添加修改时间
func (article *Article) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("modified_on", int(time.Now().Unix()))
	return nil
}

func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id=?", id).First(&article)
	if article.ID > 0 {
		return true
	} else {
		return false
	}
}

func GetArticle(id int) (result Article) {
	db.Preload("Tag").Where("id=?", id).First(&result)
	return
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (results []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&results)
	return
}

func GetArticleCount(maps interface{}) int {
	var count int64
	db.Model(&Article{}).Where(maps).Count(&count)
	return int(count)
}

func AddArticle(data map[string]interface{}) bool {
	article := Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	}
	db.Create(&article)

	return true
}

func EditArticle(id int, data map[string]interface{}) bool {
	db.Model(&Article{}).Where("id=?", id).Updates(data)
	return true
}

func DeleteArticle(id int) bool {
	var article Article
	article.ID = id
	db.Delete(&article)
	return true
}
