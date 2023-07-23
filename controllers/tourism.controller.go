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
	var tourismPictures []models.TourismPicture
	var TourismResponse []models.TourismResponse

	// get all tourisms
	result := tc.DB.Find(&tourisms)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	// get all tourism pictures
	result2 := tc.DB.Find(&tourismPictures)

	if result2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result2.Error.Error()})
		return
	}

	// merge tourisms and tourism pictures with tourism response struct on same tourism id
	for _, tourism := range tourisms {
		var tourismResponse models.TourismResponse

		tourismResponse.ID = tourism.ID
		tourismResponse.Title = tourism.Title
		tourismResponse.Slug = tourism.Slug
		tourismResponse.Description = tourism.Description
		tourismResponse.Latitude = tourism.Latitude
		tourismResponse.Longitude = tourism.Longitude
		tourismResponse.CoverPictureUrl = tourism.CoverPictureUrl
		tourismResponse.Pictures = []models.TourismPictureResponse{}

		for _, tourismPicture := range tourismPictures {
			if tourismPicture.TourismID == tourism.ID {
				tourismResponse.Pictures = append(tourismResponse.Pictures, models.TourismPictureResponse{
					ID:         tourismPicture.ID,
					PictureUrl: tourismPicture.PictureUrl,
					TourismID:  tourismPicture.TourismID,
				})
			}
		}

		TourismResponse = append(TourismResponse, tourismResponse)
	}

	ctx.JSON(http.StatusOK, gin.H{"tourisms": TourismResponse})
}

func (tc *TourismController) GetOne(ctx *gin.Context) {
	var tourism models.Tourism
	var tourismPictures []models.TourismPicture
	var tourismResponse models.TourismResponse

	// Get one tourism by slug instead of id
	result := tc.DB.Where("slug = ?", ctx.Param("slug")).First(&tourism)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	// Get all tourism pictures with the same tourism id
	result2 := tc.DB.Where("tourism_id = ?", tourism.ID).Find(&tourismPictures)

	if result2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result2.Error.Error()})
		return
	}

	// Merge tourism and tourism pictures with the tourism response struct on the same tourism id
	tourismResponse.ID = tourism.ID
	tourismResponse.Title = tourism.Title
	tourismResponse.Slug = tourism.Slug
	tourismResponse.Description = tourism.Description
	tourismResponse.Latitude = tourism.Latitude
	tourismResponse.Longitude = tourism.Longitude
	tourismResponse.CoverPictureUrl = tourism.CoverPictureUrl
	tourismResponse.Pictures = []models.TourismPictureResponse{}

	for _, tourismPicture := range tourismPictures {
		if tourismPicture.TourismID == tourism.ID {
			tourismResponse.Pictures = append(tourismResponse.Pictures, models.TourismPictureResponse{
				ID:         tourismPicture.ID,
				PictureUrl: tourismPicture.PictureUrl,
				TourismID:  tourismPicture.TourismID,
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{"tourism": tourismResponse}})
}

func (tc *TourismController) Create(ctx *gin.Context) {
	var payload *models.TourismInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// User Make new tourism with tourismpicture
	// Make new tourism with tourismpicture
	// Return tourism with tourismpicture json response
	newTourism := models.Tourism{
		Title:           payload.Title,
		Description:     payload.Description,
		Slug:            payload.Slug,
		Latitude:        payload.Latitude,
		Longitude:       payload.Longitude,
		CoverPictureUrl: payload.CoverPictureUrl,
	}

	result := tc.DB.Create(&newTourism)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	// look for payload, if payload exist, loop through array, then create new tourismpicture
	// if payload not exist, return newTourism
	if payload.Pictures != nil {
		for _, picture := range payload.Pictures {
			tourismPicture := models.TourismPicture{
				PictureUrl: picture.PictureUrl,
				TourismID:  newTourism.ID,
			}
			tc.DB.Create(&tourismPicture)

		}

	}

	response := &models.TourismResponse{
		ID:              newTourism.ID,
		Title:           newTourism.Title,
		Description:     newTourism.Description,
		Slug:            newTourism.Slug,
		CoverPictureUrl: newTourism.CoverPictureUrl,
		Latitude:        newTourism.Latitude,
		Longitude:       newTourism.Longitude,
		CreatedAt:       newTourism.CreatedAt,
		UpdatedAt:       newTourism.UpdatedAt,
		Pictures:        []models.TourismPictureResponse{},
	}

	ctx.JSON(http.StatusOK, gin.H{"tourism": response})
}

// Update tourism with tourismpicture
// Return tourism with tourismpicture json response
func (tc *TourismController) Update(ctx *gin.Context) {
	var payload *models.TourismInput
	var tourism models.Tourism

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result := tc.DB.Where("id = ?", ctx.Param("id")).First(&tourism)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	tourism.Title = payload.Title
	tourism.Description = payload.Description
	tourism.Slug = payload.Slug
	tourism.Latitude = payload.Latitude
	tourism.Longitude = payload.Longitude
	tourism.CoverPictureUrl = payload.CoverPictureUrl

	result = tc.DB.Save(&tourism)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	response := &models.TourismResponse{
		ID:              tourism.ID,
		Title:           tourism.Title,
		Description:     tourism.Description,
		Slug:            tourism.Slug,
		CoverPictureUrl: tourism.CoverPictureUrl,
		Latitude:        tourism.Latitude,
		Longitude:       tourism.Longitude,
		CreatedAt:       tourism.CreatedAt,
		UpdatedAt:       tourism.UpdatedAt,
		Pictures:        []models.TourismPictureResponse{},
	}

	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{"tourism": response}})
}

// Delete tourism with tourismpicture
// Return tourism with tourismpicture json response
func (tc *TourismController) Delete(ctx *gin.Context) {
	var tourism models.Tourism
	var tourismpicture models.TourismPicture

	result := tc.DB.Where("tourism_id = ?", ctx.Param("id")).Delete(&tourismpicture)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	result2 := tc.DB.Where("id = ?", ctx.Param("id")).First(&tourism)

	if result2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	result2 = tc.DB.Delete(&tourism)

	if result2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tourism": tourism})
}
