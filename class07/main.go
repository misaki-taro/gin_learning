package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var secretKey string = ""

type User struct {
	gorm.Model
	ID       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

type Todo struct {
	gorm.Model
	ID     int    `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Status string `json:"status"`
	UserID int    `json:"user_id"`
}

// 连接数据库
func initDB() (*gorm.DB, error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 建表
	err = db.AutoMigrate(&User{}, &Todo{})
	if err != nil {
		panic("failed to connect database")
	}

	return db, nil
}

// 用户注册
func signUp(ctx *gin.Context, db *gorm.DB) {
	// 判断数据是否正确
	var user User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 判断用户名是否存在
	var existingUser User
	db.Where("username = ?", user.Username).First(existingUser)
	if existingUser.ID != 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "username already exist...",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Password = string(hashedPassword)
	db.Create(&user)
	ctx.JSON(http.StatusOK, user)
}

// 用户登录
func signIn(ctx *gin.Context, db *gorm.DB) {
	// 判断数据是否正确
	var user User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 判断用户名是否存在
	var existingUser User
	db.Where("username = ?", user.Username).First(&existingUser)
	if existingUser.ID == 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	//jwt生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       existingUser.ID,
		"username": existingUser.Username,
	})

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": signedToken,
	})
}

// 登陆验证的中间件
func authenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("token")

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header not provided",
			})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			ctx.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		ctx.Set("userID", int(claims["id"].(float64)))
		ctx.Next()

	}
}

// 创建待办事项
func createTodo(ctx *gin.Context, db *gorm.DB) {
	var todo Todo

	if err := ctx.BindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo.UserID = ctx.GetInt("userID")
	db.Create(&todo)
	ctx.JSON(http.StatusOK, todo)

}

// 查找待办事项
func getTodo(ctx *gin.Context, db *gorm.DB) {
	var todos []Todo
	db.Where("user_id = ?", ctx.GetInt("userID")).Find(&todos)
	ctx.JSON(http.StatusOK, todos)
}

// 更新待办事项
func updateTodo(ctx *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedTodo Todo
	if err := ctx.BindJSON(&updatedTodo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&Todo{}).Where("id = ?", id).Updates(updatedTodo)
	ctx.JSON(http.StatusOK, updatedTodo)
}

// 删除待办事项
func deleteTodo(ctx *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var todo Todo
	db.Where("id = ? AND user_id = ?", id, ctx.GetInt("userID")).First(&todo)
	if todo.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	db.Where("id = ?", id).Delete(&Todo{})
	ctx.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to connect database: %v", err)
	}

	r := gin.Default()

	r.POST("/signup", func(ctx *gin.Context) {
		signUp(ctx, db)
	})

	r.POST("/signin", func(ctx *gin.Context) {
		signIn(ctx, db)
	})

	authorized := r.Group("/")
	authorized.Use(authenticationMiddleware())
	{
		authorized.POST("/todos", func(ctx *gin.Context) {
			createTodo(ctx, db)
		})

		authorized.GET("/todos", func(ctx *gin.Context) {
			getTodo(ctx, db)
		})

		authorized.PUT("/todos/:id", func(ctx *gin.Context) {
			updateTodo(ctx, db)
		})

		authorized.DELETE("/todos/:id", func(ctx *gin.Context) {
			deleteTodo(ctx, db)
		})

		authorized.POST("/logout", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Logged out"})
		})

	}

	r.Run(":8000")

}
