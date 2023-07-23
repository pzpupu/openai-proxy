package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

// JWT 密钥
var Secret = []byte("")

func CreateJwt(name string) string {
	// 创建自定义声明（Payload）
	claims := jwt.MapClaims{
		"name": name,
		"exp":  time.Now().Add(time.Hour * 24 * 365).Unix(), // 设置过期时间为1年后过期
	}

	// 使用密钥签名JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成最终的JWT字符串
	tokenString, err := token.SignedString(Secret)
	if err != nil {
		fmt.Println("Failed to generate JWT:", err)
		return ""
	}
	fmt.Println("JWT:", tokenString)
	return tokenString
}

func ValidJwt(tokenString string) *string {
	// 解析和校验JWT
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		fmt.Println("Failed to parse JWT:", tokenString, err)
		return nil
	}
	if parsedToken.Valid {
		claims := parsedToken.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		return &name
	}
	return nil

}
