package controllers

import (
	"net/http"
	"strconv"
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

	// Parse query parameters
	limitStr := ctx.DefaultQuery("limit", "20")
	offsetStr := ctx.DefaultQuery("offset", "0")
	sortby := ctx.DefaultQuery("sortby", "created_at")
	orderedby := ctx.DefaultQuery("orderedby", "desc")

	// Convert limit and offset to integers
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		// Handle error, invalid limit value
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid limit value"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		// Handle error, invalid offset value
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid offset value"})
		return
	}

	// Query the database with the parsed parameters
	result := gc.DB.
		Order(sortby + " " + orderedby).
		Limit(limit).
		Offset(offset).
		Find(&galleries)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	// Get the total count of galleries
	var total int64
	gc.DB.Model(&models.Gallery{}).Count(&total)

	// Create the meta object
	meta := gin.H{
		"limit":  limit,
		"offset": offset,
		"total":  total,
	}

	// Create the response payload
	response := gin.H{
		"galleries": galleries,
		"meta":      meta,
	}

	ctx.JSON(http.StatusOK, response)
}

func (gc *GalleryController) GetOne(ctx *gin.Context) {
	var gallery models.Gallery

	result := gc.DB.Where("id = ?", ctx.Param("id")).First(&gallery)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"gallery": gallery})
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
