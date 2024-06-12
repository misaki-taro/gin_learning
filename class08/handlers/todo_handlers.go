package handlers

import (
	"class08/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 创建待办事项
func (h *Handlers) CreateTodo(ctx *gin.Context) {
	var todo model.Todo

	if err := ctx.BindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo.UserID = ctx.GetInt("userID")
	h.DB.Create(&todo)
	ctx.JSON(http.StatusOK, todo)

}

// 查找待办事项
func (h *Handlers) GetTodo(ctx *gin.Context) {
	var todos []model.Todo
	h.DB.Where("user_id = ?", ctx.GetInt("userID")).Find(&todos)
	ctx.JSON(http.StatusOK, todos)
}

// 更新待办事项
func (h *Handlers) UpdateTodo(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedTodo model.Todo
	if err := ctx.BindJSON(&updatedTodo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Model(&model.Todo{}).Where("id = ?", id).Updates(updatedTodo)
	ctx.JSON(http.StatusOK, updatedTodo)
}

// 删除待办事项
func (h *Handlers) DeleteTodo(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var todo model.Todo
	h.DB.Where("id = ? AND user_id = ?", id, ctx.GetInt("userID")).First(&todo)
	if todo.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	h.DB.Where("id = ?", id).Delete(&model.Todo{})
	ctx.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
