package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("/api")

	//二级路由
	users := v1.Group("/users")

	users.POST("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": "users POST OK",
		})
	})

	//二级路由
	v2 := v1.Group("/v2")

	v2.POST("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": "v2 post ok",
		})
	})

	r.Run(":8000")
}
