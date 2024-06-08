package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type UpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	r := gin.Default()

	dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 建表
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to connect database")
	}

	// 逻辑
	// GET /users 查
	r.GET("/users", func(ctx *gin.Context) {
		var users []User
		if err = db.Find(&users).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to find users",
			})

			log.Println("Failed to find users")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": users,
		})
	})

	// Post 增
	r.POST("/users", func(ctx *gin.Context) {
		var input User
		if err = ctx.ShouldBindJSON(&input); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			log.Println("bind failed")
			return
		}

		if err = db.Create(&input).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
			log.Println("create failed")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": input,
		})
	})

	// Post 改
	r.PUT("/users/:id", func(ctx *gin.Context) {
		var input UpdateUserInput
		if err = ctx.ShouldBindJSON(&input); err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			log.Println("User not found")
			return
		}

		id := ctx.Param("id")
		var user User
		if err = db.First(&user, id).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update user",
			})
			log.Println("Failed to update user")
			return
		}

		if input.Name != "" {
			user.Name = input.Name
		}

		if input.Email != "" {
			user.Email = input.Email
		}

		if err = db.Save(&user).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save",
			})
			log.Println("Failed to save")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": user,
		})

	})

	// Post 删
	r.DELETE("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var user User
		if err = db.First(&user, id).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update user",
			})
			log.Println("Failed to delete User")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": "User delete OK",
		})
	})

	r.Run(":8000")

}
