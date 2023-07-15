package controllers

import (
	"net/http"

	"kebondowo/models"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type TourismController struct {
	DB *gorm.DB
}

func NewTourismController(DB *gorm.DB) TourismController {
	return TourismController{DB}
}

func (tc *TourismController) GetAll(ctx *gin.Context) {
	var tourisms []models.Tourism

	result := tc.DB.Find(&tourisms)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"tourisms": tourisms}})
}

func (tc *TourismController) GetOne(ctx *gin.Context) {
	var tourism models.Tourism

	result := tc.DB.Where("id = ?", ctx.Param("id")).First(&tourism)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"tourism": tourism}})
}



func (tc *TourismController) Create(ctx *gin.Context) {
	var payload *models.TourismInput
	

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// User Make new tourism with tourismpicture
	// Make new tourism with tourismpicture
	// Return tourism with tourismpicture json response
	newTourism := models.Tourism{
		Title: payload.Title,
		Description: payload.Description,
		Slug : payload.Slug,
		CoverPictureUrl: payload.CoverPictureUrl,
	}

	result := tc.DB.Create(&newTourism)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	// make tourism picture with tourism id
	// return tourism picture json response
	tourismPicture := models.TourismPicture{
		TourismID: newTourism.ID,
		PictureUrl: "nyoba",
		Latitude: 7.1,
		Longitude: 7.2,
	}

	result2 := tc.DB.Create(&tourismPicture)
	


	if result2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result2.Error.Error()})
		return
	}

	response := &models.TourismResponse{
		ID : newTourism.ID,
		Title : newTourism.Title,
		Description : newTourism.Description,
		Slug : newTourism.Slug,
		CoverPictureUrl : newTourism.CoverPictureUrl,
		Pictures : []models.TourismPicture{tourismPicture},
	}

	


	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"tourism": response}})
}


