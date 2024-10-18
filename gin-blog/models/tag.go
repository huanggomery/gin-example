package models

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
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
