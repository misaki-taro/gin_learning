package users

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
	users := r.Group("/users")

	users.Use(Logger())

	//GET
	users.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": "users GET OK",
		})
	})
}
