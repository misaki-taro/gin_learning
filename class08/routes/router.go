package routes

import (
	"class08/database"
	"class08/handlers"
	"class08/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Run(addr ...string) (err error)
}

func setupRouter() *gin.Engine {

	db, err := database.InitDB()
	if err != nil {
		log.Print(err)
	}

	h := handlers.Handlers{DB: db}

	r := gin.Default()

	r.POST("/signup", h.SignUp)

	r.POST("/signin", h.SignIn)

	authorized := r.Group("/")
	authorized.Use(middleware.Auth())
	{
		authorized.POST("/todos", h.CreateTodo)

		authorized.GET("/todos", h.GetTodo)

		authorized.PUT("/todos/:id", h.UpdateTodo)

		authorized.DELETE("/todos/:id", h.DeleteTodo)
		authorized.POST("/logout", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Logged out"})
		})

	}
	return r
}

func NewServer() Server {
	return setupRouter()
}
