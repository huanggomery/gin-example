package models

type Auth struct {
    ID       int    `gorm:"primary_key" json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
}

func (Auth) TableName() string {
    return "blog_auth"
}

// 查询数据库，判断用户名和密码是否存在
func CheckAuth(username, password string) bool {
    var auth Auth
    maps := map[string]string{
        "username": username,
        "password": password,
    }
    db.Select("id").Where(maps).First(&auth)
    if auth.ID > 0 {
        return true
    } else {
        return false
    }
}
