package main

import (
	"class06/posts"
	"class06/users"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	users.Router(r)
	posts.Router(r)

	r.Run(":8000")
}
