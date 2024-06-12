package handlers

import (
	"class08/config"
	"class08/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handlers struct {
	DB *gorm.DB
}

// 用户注册
func (h *Handlers) SignUp(ctx *gin.Context) {
	// 判断数据是否正确
	var user model.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 判断用户名是否存在
	var existingUser model.User
	h.DB.Where("username = ?", user.Username).First(existingUser)
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
	h.DB.Create(&user)
	ctx.JSON(http.StatusOK, user)
}

// 用户登录
func (h *Handlers) SignIn(ctx *gin.Context) {
	// 判断数据是否正确
	var user model.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 判断用户名是否存在
	var existingUser model.User
	h.DB.Where("username = ?", user.Username).First(&existingUser)
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

	config.Init()

	signedToken, err := token.SignedString([]byte(config.Cfg.SecretKey))
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
