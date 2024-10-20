package util

import (
	"errors"
	"gin-example/gin-blog/setting"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// JWT签名算法的密钥
var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 将用户名和密码加密成为Token字符串
func GenerateToken(username, password string) (string, error) {
	timeExpire := time.Now().Add(time.Hour * 3).Unix()
	claims := Claims{
		Username: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: timeExpire,
			Issuer:    "gin-blog",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

// 解码和校验Token字符串，如果有效，则将其解析成Claims结构体
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("token invalid")
	}
}
