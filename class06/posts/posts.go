package posts

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//记录日志
		fmt.Println("start...")

		//业务
		ctx.Next()
	}
}

func Router(r *gin.Engine) {
	posts := r.Group("/posts")

	posts.Use(Logger())

	//GET
	posts.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": "posts GET OK",
		})
	})
}
