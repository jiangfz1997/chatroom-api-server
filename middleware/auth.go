package middleware

import (
	"chatroom-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 鉴权中间件：验证 JWT Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少或无效的 Authorization 头"})
			c.Abort()
			return
		}

		// 提取 token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析 token 获取用户名
		username, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效或已过期"})
			c.Abort()
			return
		}

		// 设置用户名到上下文中，供后续 handler 使用
		c.Set("username", username)

		c.Next()
	}
}
