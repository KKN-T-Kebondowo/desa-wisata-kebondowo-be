package controllers

import (
	"net/http"
	"time"

	"kebondowo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardController struct {
	DB *gorm.DB
}

func NewDashboardController(DB *gorm.DB) DashboardController {
	return DashboardController{DB}
}

func (dc *DashboardController) GetDashboard(ctx *gin.Context) {
	// get total article
	var totalArticle int64
	dc.DB.Model(&models.Article{}).Count(&totalArticle)

	// get total gallery
	var totalGallery int64
	dc.DB.Model(&models.Gallery{}).Count(&totalGallery)

	// get total tourism
	var totalTourism int64
	dc.DB.Model(&models.Tourism{}).Count(&totalTourism)

	var monthNow = int(time.Now().Month())
	var yearNow = int(time.Now().Year())


	// get article per month
	var articlePerMonth []int32
	for i := monthNow+1; i <= 12; i++ {
		var count int32
		dc.DB.Raw("SELECT COUNT(*) FROM articles WHERE EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", i, yearNow-1).Scan(&count)
		articlePerMonth = append(articlePerMonth, count)
	}

	// ambil data dari tahun sekarang
	for i := 1; i <= monthNow; i++ {
		var count int32
		dc.DB.Raw("SELECT COUNT(*) FROM articles WHERE EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", i, yearNow).Scan(&count)
		articlePerMonth = append(articlePerMonth, count)
	}
	


	// get gallery per month
	var galleryPerMonth []int32
	for i := monthNow+1; i <= 12; i++ {
		var count int32
		dc.DB.Raw("SELECT COUNT(*) FROM galleries WHERE EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", i, yearNow-1).Scan(&count)
		galleryPerMonth = append(galleryPerMonth, count)
	}

	for i := 1; i <= monthNow; i++ {
		var count int32
		dc.DB.Raw("SELECT COUNT(*) FROM galleries WHERE EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", i, yearNow).Scan(&count)
		galleryPerMonth = append(galleryPerMonth, count)
	}

	// get tourism per month
	var tourismPerMonth []int32
	for i := monthNow+1; i <= 12; i++ {
		var count int32
		dc.DB.Raw("SELECT COUNT(*) FROM tourisms WHERE EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", i, yearNow-1).Scan(&count)
		tourismPerMonth = append(tourismPerMonth, count)
	}

	for i := 1; i <= monthNow; i++ {
		var count int32
		dc.DB.Raw("SELECT COUNT(*) FROM tourisms WHERE EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", i, yearNow).Scan(&count)
		tourismPerMonth = append(tourismPerMonth, count)
	}



	dashboardResponse := &models.DashboardResponse{
		TotalArticle:    totalArticle,
		TotalGallery:    totalGallery,
		TotalTourism:    totalTourism,
		ArticlePerMonth: articlePerMonth,
		GalleryPerMonth: galleryPerMonth,
		TourismPerMonth: tourismPerMonth,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": dashboardResponse})

}
