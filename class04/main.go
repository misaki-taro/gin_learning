package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 创建路由分组
	api := r.Group("/api")

	// 使用中间件
	api.Use(logger())

	// GET
	api.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": "get user",
		})
	})

	// POST
	// PUT
	// DELETE

	r.Run(":8000")
}

// 中间件
func logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()

		//请求前的日志
		log.Printf("[%s] %s %s \n", t.Format("2006-01-02 15:04:05"), ctx.Request.Method, ctx.Request.URL.Path)

		//处理请求
		ctx.Next()

		//请求后的日志
		log.Printf("[%s] %s %s %s \n", t.Format("2006-01-02 15:04:05"), ctx.Request.Method, ctx.Request.URL.Path, time.Since(t))
	}
}
