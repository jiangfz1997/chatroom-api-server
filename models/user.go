package models

// // User 表结构，用于 GORM 创建数据表
//
//	type User struct {
//		ID       uint   `gorm:"primaryKey"` // 主键，自动增长
//		Username string `gorm:"unique"`     // 用户名唯一
//		Password string // 密码（此处明文存储，仅开发用，生产应加密）
//	}
type User struct {
	Username string `dynamodbav:"username"` // 👈 主键
	Password string `dynamodbav:"password"`
}
