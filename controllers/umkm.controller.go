package controllers

import (
	"net/http"
	"time"

	"kebondowo/models"

	"strconv"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type UMKMController struct {
	DB *gorm.DB
}

func NewUMKMController(DB *gorm.DB) UMKMController {
	return UMKMController{DB}
}

func (uc *UMKMController) GetAll(ctx *gin.Context) {
	var umkms []models.UMKM
	var umkmPictures []models.UMKMPicture
	var UMKMResponse []models.UMKMResponse

	// get all umkms
	result := uc.DB.Find(&umkms)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	// get all umkm pictures
	result2 := uc.DB.Find(&umkmPictures)

	if result2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result2.Error.Error()})
		return
	}

	// merge umkms and umkm pictures with umkm response struct on same umkm id
	for _, umkm := range umkms {
		var umkmResponse models.UMKMResponse

		umkmResponse.ID = umkm.ID
		umkmResponse.Title = umkm.Title
		umkmResponse.Slug = umkm.Slug
		umkmResponse.Description = umkm.Description
		umkmResponse.Latitude = umkm.Latitude
		umkmResponse.Longitude = umkm.Longitude
		umkmResponse.CoverPictureUrl = umkm.CoverPictureUrl
		umkmResponse.Visitor = umkm.Visitor
		umkmResponse.Contact = umkm.Contact
		umkmResponse.ContactName = umkm.ContactName
		umkmResponse.Pictures = []models.UMKMPictureResponse{}

		for _, umkmPicture := range umkmPictures {
			if umkmPicture.UMKMID == umkm.ID {
				umkmResponse.Pictures = append(umkmResponse.Pictures, models.UMKMPictureResponse{
					ID:         umkmPicture.ID,
					UMKMID:     umkmPicture.UMKMID,
					PictureUrl: umkmPicture.PictureUrl,
				})
			}
		}

		UMKMResponse = append(UMKMResponse, umkmResponse)
	}

	ctx.JSON(http.StatusOK, gin.H{"umkms": UMKMResponse})
}

func (uc *UMKMController) GetOne(ctx *gin.Context) {
	var umkm models.UMKM
	var umkmPictures []models.UMKMPicture
	var umkmResponse models.UMKMResponse

	// get umkm
	result := uc.DB.Where("slug = ?", ctx.Param("slug")).First(&umkm)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	umkm.Visitor++
	uc.DB.Model(&umkm).UpdateColumn("visitor", umkm.Visitor)

	// get umkm pictures
	result2 := uc.DB.Where("umkm_id = ?", umkm.ID).Find(&umkmPictures)

	if result2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result2.Error.Error()})
		return
	}

	// merge umkm and umkm pictures with umkm response struct on same umkm id
	umkmResponse.ID = umkm.ID
	umkmResponse.Title = umkm.Title
	umkmResponse.Slug = umkm.Slug
	umkmResponse.Description = umkm.Description
	umkmResponse.Latitude = umkm.Latitude
	umkmResponse.Longitude = umkm.Longitude
	umkmResponse.Visitor = umkm.Visitor
	umkmResponse.Contact = umkm.Contact
	umkmResponse.ContactName = umkm.ContactName
	umkmResponse.CoverPictureUrl = umkm.CoverPictureUrl
	umkmResponse.Pictures = []models.UMKMPictureResponse{}

	for _, umkmPicture := range umkmPictures {
		if umkmPicture.UMKMID == umkm.ID {
			umkmResponse.Pictures = append(umkmResponse.Pictures, models.UMKMPictureResponse{
				ID:         umkmPicture.ID,
				UMKMID:     umkmPicture.UMKMID,
				PictureUrl: umkmPicture.PictureUrl,
				Caption:   umkmPicture.Caption,
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"umkm": umkmResponse})
}

func (uc *UMKMController) Create(ctx *gin.Context) {
	var payload *models.UMKMInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var slug string
	var count int64

	uc.DB.Model(&models.UMKM{}).Where("slug = ?", payload.Slug).Count(&count)
	if count > 0 {
		slug = payload.Slug + "-" + strconv.FormatInt(count, 10)
	} else {
		slug = payload.Slug
	}
		
	

	payload.Slug = slug


	newUMKM := models.UMKM{
		Title: 		 payload.Title,
		Slug:        payload.Slug,
		Description: payload.Description,
		Latitude:    payload.Latitude,
		Longitude:   payload.Longitude,
		CoverPictureUrl: payload.CoverPictureUrl,
		Contact: payload.Contact,
		ContactName: payload.ContactName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := uc.DB.Create(&newUMKM)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	// create umkm pictures
	if payload.Pictures != nil {
		for _, picture := range payload.Pictures {
			newUMKMPicture := models.UMKMPicture{
				UMKMID:      newUMKM.ID,
				PictureUrl: picture.PictureUrl,
			}

			result2 := uc.DB.Create(&newUMKMPicture)

			if result2.Error != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"message": result2.Error.Error()})
				return
			}
		}
	}

	response := &models.UMKMResponse{
		ID:          newUMKM.ID,
		Title:       newUMKM.Title,
		Slug:        newUMKM.Slug,
		Description: newUMKM.Description,
		Latitude:    newUMKM.Latitude,
		Longitude:   newUMKM.Longitude,
		Contact:  newUMKM.Contact,
		ContactName:  newUMKM.ContactName,
		CoverPictureUrl: newUMKM.CoverPictureUrl,
		CreatedAt:   newUMKM.CreatedAt,
		UpdatedAt:   newUMKM.UpdatedAt,
		Pictures:   []models.UMKMPictureResponse{},
	}

	ctx.JSON(http.StatusOK, gin.H{"umkm": response})
}


func (uc *UMKMController) Update(ctx *gin.Context) {
	var payload *models.UMKMUpdate
	var umkm models.UMKM

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result := uc.DB.Where("id = ?", ctx.Param("id")).First(&umkm)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	umkm.Title = payload.Title
	umkm.Description = payload.Description
	umkm.Latitude = payload.Latitude
	umkm.Longitude = payload.Longitude
	umkm.Slug = payload.Slug
	if payload.CoverPictureUrl != "" {
		umkm.CoverPictureUrl = payload.CoverPictureUrl
	}

	result2 := uc.DB.Save(&umkm)

	if result2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result2.Error.Error()})
		return
	}

	response := &models.UMKMResponse{
		ID:          umkm.ID,
		Title:       umkm.Title,
		Slug:        umkm.Slug,
		Description: umkm.Description,
		Latitude:    umkm.Latitude,
		Longitude:   umkm.Longitude,
		CoverPictureUrl: umkm.CoverPictureUrl,
		CreatedAt:   umkm.CreatedAt,
		UpdatedAt:   umkm.UpdatedAt,
		Pictures:   []models.UMKMPictureResponse{},
	}

	ctx.JSON(http.StatusOK, gin.H{"umkm": response})
}

func (uc *UMKMController) Delete(ctx *gin.Context) {
	var umkm models.UMKM

	result := uc.DB.Where("slug = ?", ctx.Param("slug")).First(&umkm)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	uc.DB.Delete(&umkm)

	ctx.JSON(http.StatusOK, gin.H{"message": "UMKM deleted"})
	
	result2 := uc.DB.Where("id = ?", ctx.Param("id")).First(&umkm)

	if result2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result2.Error.Error()})
		return
	}

	uc.DB.Delete(&umkm)

	ctx.JSON(http.StatusOK, gin.H{"message": "UMKM deleted"})


}



