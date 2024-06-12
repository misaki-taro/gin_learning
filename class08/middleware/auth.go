package middleware

import (
	"class08/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// 登陆验证的中间件
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		config.Init()
		tokenString := ctx.GetHeader("token")

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header not provided",
			})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.SecretKey), nil
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
