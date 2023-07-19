package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

// 设置密钥
var secret = []byte("your-secret-key")

func CreateJwt(name string) string {
	// 创建自定义声明（Payload）
	claims := jwt.MapClaims{
		"username": name,
		"exp":      time.Now().Add(time.Hour * 24 * 365).Unix(), // 设置过期时间为1年后过期
	}

	// 使用密钥签名JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成最终的JWT字符串
	tokenString, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("Failed to generate JWT:", err)
		return ""
	}
	fmt.Println("JWT:", tokenString)
	return tokenString
}

func ValidJwt(tokenString string) bool {
	// 解析和校验JWT
	parsedToken, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	return parsedToken.Valid
}

//func main() {
//	tokenString := CreateJwt("TestUser")
//	ValidJwt(tokenString + "1")
//}
