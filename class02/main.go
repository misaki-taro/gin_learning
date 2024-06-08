package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// GET POST PUT DELETE

	r.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": "Get all users",
		})
	})

	r.POST("/users", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": "create Users",
		})
	})

	r.PUT("/users/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.JSON(http.StatusOK, gin.H{
			"data": "Update User " + id,
		})
	})

	r.POST("users/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.JSON(http.StatusOK, gin.H{
			"data": "Delete User " + id,
		})
	})

	r.Run(":8000")

}
