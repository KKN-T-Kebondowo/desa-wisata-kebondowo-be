package controllers

import (
	"net/http"
	"time"

	"kebondowo/models"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type UMKMPictureController struct {
	DB *gorm.DB
}

func NewUMKMPictureController(DB *gorm.DB) UMKMPictureController {
	return UMKMPictureController{DB}
}

func (upc *UMKMPictureController) GetAll(ctx *gin.Context) {
	var umkmPictures []models.UMKMPicture
	var umkmPictureResponse []models.UMKMPictureResponse

	result := upc.DB.Find(&umkmPictures)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	for _, umkmPicture := range umkmPictures {
		umkmPictureResponse = append(umkmPictureResponse, models.UMKMPictureResponse{
			ID:         umkmPicture.ID,
			PictureUrl: umkmPicture.PictureUrl,
			UMKMID:     umkmPicture.UMKMID,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"umkm_pictures": umkmPictureResponse})

}

func (upc *UMKMPictureController) GetOne(ctx *gin.Context) {
	var umkmPicture models.UMKMPicture
	var umkmPictureResponse models.UMKMPictureResponse

	result := upc.DB.Where("id = ?", ctx.Param("id")).First(&umkmPicture)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	umkmPictureResponse.ID = umkmPicture.ID
	umkmPictureResponse.PictureUrl = umkmPicture.PictureUrl
	umkmPictureResponse.UMKMID = umkmPicture.UMKMID

	ctx.JSON(http.StatusOK, gin.H{"umkm_picture": umkmPictureResponse})
}

func (upc *UMKMPictureController) Create(ctx *gin.Context) {
	var payload *models.UMKMPictureInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	newUMKMPicture := models.UMKMPicture{
		PictureUrl: payload.PictureUrl,
		UMKMID:     payload.UMKMID,
		Caption:  payload.Caption,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := upc.DB.Create(&newUMKMPicture)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	response := &models.UMKMPictureResponse{
		ID:         newUMKMPicture.ID,
		PictureUrl: newUMKMPicture.PictureUrl,
		Caption: newUMKMPicture.Caption,
		UpdatedAt: newUMKMPicture.UpdatedAt,
		CreatedAt:  newUMKMPicture.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{"umkm_picture": response})

}

func (upc *UMKMPictureController) Delete(ctx *gin.Context) {
	var umkmPicture models.UMKMPicture

	result := upc.DB.Where("id = ?", ctx.Param("id")).First(&umkmPicture)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "UMKM picture deleted"})
}
