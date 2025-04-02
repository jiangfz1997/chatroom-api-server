package handlers

//type RegisterRequest struct {
//	Username string `json:"username"`
//	Password string `json:"password"`
//}

//func Register(c *gin.Context) {
//	var req RegisterRequest
//	fmt.Println("收到注册请求：", req.Username, req.Password)
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "格式错误"})
//		return
//	}
//
//	// 检查用户名是否已存在
//	var exists int
//	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&exists)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
//		return
//	}
//	if exists > 0 {
//		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
//		return
//	}
//
//	// 插入用户数据
//	_, err = database.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", req.Username, req.Password)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
//	fmt.Println("写入数据库成功")
//}

//func Login(c *gin.Context) {
//	var req models.User
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
//		return
//	}
//
//	var storedPassword string
//	err := database.DB.QueryRow("SELECT password FROM users WHERE username = ?", req.Username).Scan(&storedPassword)
//
//	if err == sql.ErrNoRows {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名不存在"})
//		return
//	}
//
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
//		return
//	}
//
//	if storedPassword != req.Password {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "username": req.Username})
//}
