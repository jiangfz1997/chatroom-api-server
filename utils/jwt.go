package utils

import (
	log "chatroom-api/logger"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // 存入环境变量中

// 创建 Token（传入用户名）
func GenerateToken(username string) (string, error) {
	log.Log.Infof("生成 Token 请求: username=%s", username)
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token 有效期 24 小时
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Log.Infof("Token 生成成功: username=%s", username)
	return token.SignedString(jwtSecret)
}

// 验证 Token 并返回用户名
func ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 校验签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Log.Warn("Token 签名算法不匹配")
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		log.Log.Warnf("Token 无效或解析失败: %v", err)
		return "", errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		username, ok := claims["username"].(string)
		if !ok {
			log.Log.Warn("Token 中未包含用户名字段")
			return "", errors.New("username not found in token")
		}
		log.Log.Infof("Token 解析成功: username=%s", username)
		return username, nil
	}

	return "", errors.New("failed to parse token claims")
}
