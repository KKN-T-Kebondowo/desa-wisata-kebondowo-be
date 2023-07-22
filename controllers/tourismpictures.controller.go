package controllers

import (
	"net/http"

	"kebondowo/models"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type TourismPictureController struct {
	DB *gorm.DB
}

func NewTourismPictureController(DB *gorm.DB) TourismPictureController {
	return TourismPictureController{DB}
}

func (tpc *TourismPictureController) GetAll(ctx *gin.Context) {
	var tourismPictures []models.TourismPicture
	var tourismPictureResponse []models.TourismPictureResponse

	result := tpc.DB.Find(&tourismPictures)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	for _, tourismPicture := range tourismPictures {
		tourismPictureResponse = append(tourismPictureResponse, models.TourismPictureResponse{
			ID:         tourismPicture.ID,
			PictureUrl: tourismPicture.PictureUrl,
			TourismID:  tourismPicture.TourismID,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"tourism_pictures": tourismPictureResponse})

}

func (tpc *TourismPictureController) GetOne(ctx *gin.Context) {
	var tourismPicture models.TourismPicture
	var tourismPictureResponse models.TourismPictureResponse

	result := tpc.DB.Where("id = ?", ctx.Param("id")).First(&tourismPicture)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	tourismPictureResponse.ID = tourismPicture.ID
	tourismPictureResponse.PictureUrl = tourismPicture.PictureUrl
	tourismPictureResponse.TourismID = tourismPicture.TourismID

	ctx.JSON(http.StatusOK, gin.H{"tourism_picture": tourismPictureResponse})
	
}

func (tpc *TourismPictureController) Create(ctx *gin.Context) {
	var payload *models.TourismPictureInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tourismPicture := models.TourismPicture{
		TourismID:  payload.TourismID,
		PictureUrl: payload.PictureUrl,
	}

	result := tpc.DB.Create(&tourismPicture)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tourism_picture": tourismPicture})
}

// delete
func (tpc *TourismPictureController) Delete(ctx *gin.Context) {
	var tourismPicture models.TourismPicture

	result := tpc.DB.Where("id = ?", ctx.Param("id")).Delete(&tourismPicture)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
