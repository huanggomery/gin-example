package models

import (
    "fmt"
    "log"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/plugin/soft_delete"

    "gin-example/gin-blog/setting"
)

var db *gorm.DB

type Model struct {
    ID         int                   `gorm:"primary_key" json:"id"`
    CreatedOn  int                   `json:"created_on"`
    ModifiedOn int                   `json:"modified_on"`
    DeletedOn  soft_delete.DeletedAt `json:"deleted_on"`
}

// 读取setting模块配置，连接数据库
func Setup() {
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        setting.DatabaseSetting.User,
        setting.DatabaseSetting.Password,
        setting.DatabaseSetting.Host,
        setting.DatabaseSetting.DbName,
    )

    // 连接到数据库
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 设置最大连接数
    mysqlDB, _ := db.DB()
    // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
    mysqlDB.SetMaxIdleConns(10)
    // SetMaxOpenConns sets the maximum number of open connections to the database.
    mysqlDB.SetMaxOpenConns(100)
}
