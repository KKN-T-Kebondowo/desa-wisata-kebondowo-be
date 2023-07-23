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

func circularShift(arr []int64, positions int) []int64 {
	n := len(arr)
	shiftedArr := make([]int64, n)
	for i := 0; i < n; i++ {
		shiftedArr[i] = arr[(i+positions)%n]
	}
	return shiftedArr
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

	now := time.Now()
	oneYearAgo := now.AddDate(-1, 0, 0)

	type MonthData struct {
		Month time.Month
		Count int64
	}

	// Fetch article data
	var articleData []MonthData
	dc.DB.Raw(`
		SELECT
			EXTRACT(MONTH FROM created_at) as month,
			COUNT(*) as count
		FROM
			articles
		WHERE
			created_at >= ? AND created_at < ?
		GROUP BY
			EXTRACT(MONTH FROM created_at)
		ORDER BY
			EXTRACT(MONTH FROM created_at)
	`, oneYearAgo, now).Scan(&articleData)

	articlePerMonth := make([]int64, 12)
	for _, data := range articleData {
		// Since the month ranges from 1 to 12, we need to subtract 1 to get the correct index in the array.
		articlePerMonth[data.Month-1] = data.Count
	}

	// Fetch gallery data
	var galleryData []MonthData
	dc.DB.Raw(`
		SELECT
			EXTRACT(MONTH FROM created_at) as month,
			COUNT(*) as count
		FROM
			galleries
		WHERE
			created_at >= ? AND created_at < ?
		GROUP BY
			EXTRACT(MONTH FROM created_at)
		ORDER BY
			EXTRACT(MONTH FROM created_at)
	`, oneYearAgo, now).Scan(&galleryData)

	galleryPerMonth := make([]int64, 12)
	for _, data := range galleryData {
		// Since the month ranges from 1 to 12, we need to subtract 1 to get the correct index in the array.
		galleryPerMonth[data.Month-1] = data.Count
	}

	// Fetch tourism data
	var tourismData []MonthData
	dc.DB.Raw(`
		SELECT
			EXTRACT(MONTH FROM created_at) as month,
			COUNT(*) as count
		FROM
			tourisms
		WHERE
			created_at >= ? AND created_at < ?
		GROUP BY
			EXTRACT(MONTH FROM created_at)
		ORDER BY
			EXTRACT(MONTH FROM created_at)
	`, oneYearAgo, now).Scan(&tourismData)

	tourismPerMonth := make([]int64, 12)
	for _, data := range tourismData {
		// Since the month ranges from 1 to 12, we need to subtract 1 to get the correct index in the array.
		tourismPerMonth[data.Month-1] = data.Count
	}

	currentMonth := int(now.Month())
	articlePerMonth = circularShift(articlePerMonth, currentMonth)
	galleryPerMonth = circularShift(galleryPerMonth, currentMonth)
	tourismPerMonth = circularShift(tourismPerMonth, currentMonth)

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
