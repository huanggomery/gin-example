package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 结构体 Tag -> 表 blog_tag
func (Tag) TableName() string {
	return "blog_tag"
}

// 新增和修改前设置改动时间
func (tag *Tag) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("created_on", int(time.Now().Unix()))
	return nil
}

func (tag *Tag) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("modified_on", int(time.Now().Unix()))
	return nil
}

func GetTags(pageNum int, pageSize int, maps interface{}) (results []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&results)
	return
}

func GetTagCount(maps interface{}) int {
	var count int64
	db.Model(&Tag{}).Where(maps).Count(&count)
	return int(count)
}

func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name=?", name).First(&tag)
	if tag.ID > 0 {
		return true
	} else {
		return false
	}
}

func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id=?", id).First(&tag)
	if tag.ID > 0 {
		return true
	} else {
		return false
	}
}

func AddTag(name string, state int, created_by string) bool {
	tag := Tag{
		Name:      name,
		CreatedBy: created_by,
		State:     state,
	}
	db.Create(&tag)
	return true
}

func DeleteTag(id int) bool {
	tag := Tag{}
	tag.ID = id
	db.Delete(&tag)
	return true
}

func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id=?", id).Updates(data)
	return true
}
