package handlers

import (
	//"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"chatroom-api/dynamodb"
	log "chatroom-api/logger"
	"chatroom-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Warn("注册参数格式错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式错误"})
		return
	}
	log.Log.Infof("用户注册请求: %s", req.Username)
	// 创建用户结构体
	user := dynamodb.User{
		Username: req.Username,
		Password: req.Password, // 注意：生产环境应加密！
	}

	err := dynamodb.CreateUser(user)
	if err != nil {
		log.Log.Warnf("用户创建失败: %v", err)

		if strings.Contains(err.Error(), "ConditionalCheckFailed") {
			log.Log.Infof("用户名已存在: %s", req.Username)
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		} else {
			log.Log.Errorf("注册失败（系统错误）: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		}
		return
	}
	log.Log.Infof("用户注册成功: %s", req.Username)
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Login(c *gin.Context) {
	log.Log.Info("🔥 Login Hit!")
	var req dynamodb.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
		return
	}

	// 查询用户
	user, err := dynamodb.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名不存在"})
		return
	}

	// 验证密码（生产中应加密）
	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}
	// 生成 JWT Token
	//临时调试
	//var req RegisterRequest
	//req.Username = "qqq"

	token, err := utils.GenerateToken(req.Username)
	if err != nil {
		log.Log.Errorf("Token 生成失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token 生成失败"})
		return
	}
	log.Log.Infof("登录成功: %s，Token 已生成", req.Username)

	c.JSON(http.StatusOK, gin.H{
		"message":  "登录成功",
		"username": req.Username,
		"token":    token, // 加上 token 字段
	})

}
