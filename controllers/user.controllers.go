package controllers

import (
	"net/http"

	"kebondowo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	val, exists := ctx.Get("currentUser")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "user not found in context"})
		return
	}
	currentUser, ok := val.(models.User)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "invalid user context"})
		return
	}

	userResponse := &models.UserResponse{
		ID:        currentUser.ID,
		Username:  currentUser.Username,
		RoleID:    currentUser.RoleID,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"user": userResponse})
}

