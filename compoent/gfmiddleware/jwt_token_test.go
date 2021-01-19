package gfmiddleware

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gogf/gf/frame/g"
)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func TestJwt(t *testing.T) {

	//设置载荷数据
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: "admin",
		Password: "123456",
		StandardClaims: jwt.StandardClaims{
			// 默认精确到微妙
			ExpiresAt: jwt.At(time.Now().Add(60 * time.Second)),
		},
	})

	// 使用密钥进行签名
	tokenString, err := token.SignedString([]byte("secret_key"))

	fmt.Println(tokenString, err)

}

func TestParseJwt(t *testing.T) {

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIxMjM0NTYiLCJleHAiOjE2MTEwNjU4ODYuNzI0NzU4fQ.phVojCjOucpyyJhJW1prSZEA1gP6x3F1PBQkRLD7Jew"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// 校验期望的alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意想不到的签名方法：%v", token.Header["alg"])
		}

		return []byte("secret_key"), nil
	})

	if tokenError, ok := err.(*jwt.MalformedTokenError); ok {
		fmt.Println(tokenError)
	}

	if err != nil {
		fmt.Println("解析异常", err.Error())
		return
	}

	// 校验token并获取载荷数据
	if claims, ok := token.Claims.(jwt.Claims); ok && token.Valid {
		g.Dump(claims)
	}

}
