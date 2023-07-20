package controllers

import (
	"net/http"

	"kebondowo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleController struct {
	DB *gorm.DB
}

func NewRoleController(DB *gorm.DB) RoleController {
	return RoleController{DB}
}

func (rc *RoleController) GetAll(ctx *gin.Context) {
	var roles []models.Role

	result := rc.DB.Find(&roles)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (rc *RoleController) GetOne(ctx *gin.Context) {
	var role models.Role

	result := rc.DB.Where("id = ?", ctx.Param("id")).First(&role)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"role": role})
}

func (rc *RoleController) Create(ctx *gin.Context) {
	var payload *models.RoleInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	newRole := models.Role{
		Name: payload.Name,
	}

	result := rc.DB.Create(&newRole)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"role": newRole})
}

// edit role
func (rc *RoleController) Update(ctx *gin.Context) {
	var payload *models.RoleInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var role models.Role

	result := rc.DB.Where("id = ?", ctx.Param("id")).First(&role)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	role.Name = payload.Name

	result = rc.DB.Save(&role)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"role": role})
}

// delete role
func (rc *RoleController) Delete(ctx *gin.Context){
	var role models.Role

	result := rc.DB.Where("id = ?", ctx.Param("id")).First(&role)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	result = rc.DB.Delete(&role)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"role": role})
}