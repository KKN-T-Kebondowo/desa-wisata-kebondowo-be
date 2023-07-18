package controllers

import (
	"net/http"
	"time"

	"kebondowo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GalleryController struct {
	DB *gorm.DB
}

func NewGalleryController(DB *gorm.DB) GalleryController {
	return GalleryController{DB}
}

func (gc *GalleryController) GetAll(ctx *gin.Context) {
	var galleries []models.Gallery

	result := gc.DB.Find(&galleries)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"galleries": galleries}})
}

func (gc *GalleryController) GetOne(ctx *gin.Context) {
	var gallery models.Gallery

	result := gc.DB.Where("id = ?", ctx.Param("id")).First(&gallery)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"gallery": gallery}})
}

func (gc *GalleryController) Create(ctx *gin.Context) {
	var payload *models.GalleryInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	newGallery := models.Gallery{
		PictureUrl: payload.PictureUrl,
		Caption:    payload.Caption,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result := gc.DB.Create(&newGallery)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	response := &models.GalleryResponse{
		ID:         newGallery.ID,
		PictureUrl: newGallery.PictureUrl,
		Caption:    newGallery.Caption,
		CreatedAt:  newGallery.CreatedAt,
		UpdatedAt:  newGallery.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"gallery": response}})
}

// delete
func (gc *GalleryController) Delete(ctx *gin.Context) {
	var gallery models.Gallery

	result := gc.DB.Where("id = ?", ctx.Param("id")).Delete(&gallery)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
