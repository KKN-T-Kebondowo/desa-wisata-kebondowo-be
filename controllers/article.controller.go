package controllers

import (
	"net/http"
	"strconv"
	"time"

	"kebondowo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleController struct {
	DB *gorm.DB
}

func NewArticleController(DB *gorm.DB) ArticleController {
	return ArticleController{DB}
}

func (ac *ArticleController) GetAll(ctx *gin.Context) {
	var articles []models.Article

	// Parse query parameters
	limitStr := ctx.DefaultQuery("limit", "20")
	offsetStr := ctx.DefaultQuery("offset", "0")
	sortby := ctx.DefaultQuery("sortby", "created_at")
	orderedby := ctx.DefaultQuery("orderedby", "desc")

	// Convert limit and offset to integers
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		// Handle error, invalid limit value
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid limit value"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		// Handle error, invalid offset value
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid offset value"})
		return
	}

	// Query the database with the parsed parameters
	result := ac.DB.
		Order(sortby + " " + orderedby).
		Limit(limit).
		Offset(offset).
		Find(&articles)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	// Get the total count of articles
	var total int64
	ac.DB.Model(&models.Article{}).Count(&total)

	// Create the meta object
	meta := gin.H{
		"limit":  limit,
		"offset": offset,
		"total":  total,
	}

	// Create the response payload
	response := gin.H{
		"articles": articles,
		"meta":     meta,
	}

	ctx.JSON(http.StatusOK, response)
}

func (ac *ArticleController) GetOne(ctx *gin.Context) {
	var article models.Article

	result := ac.DB.Where("id = ?", ctx.Param("id")).First(&article)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"article": article})
}

func (ac *ArticleController) Create(ctx *gin.Context) {
	var payload *models.ArticleInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	newArticle := models.Article{
		Title:      payload.Title,
		Slug:       payload.Slug,
		Author:     payload.Author,
		Content:    payload.Content,
		PictureUrl: payload.PictureUrl,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result := ac.DB.Create(&newArticle)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	response := &models.ArticleResponse{
		ID:         newArticle.ID,
		Title:      newArticle.Title,
		Slug:       newArticle.Slug,
		Author:     newArticle.Author,
		Content:    newArticle.Content,
		PictureUrl: newArticle.PictureUrl,
		CreatedAt:  newArticle.CreatedAt,
		UpdatedAt:  newArticle.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"article": response})
}

func (ac *ArticleController) Update(ctx *gin.Context) {
	var payload *models.ArticleInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Get the article ID from the URL parameter
	articleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid article ID"})
		return
	}

	// Fetch the existing article from the database
	var existingArticle models.Article
	result := ac.DB.Where("id = ?", articleID).First(&existingArticle)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Article not found"})
		return
	}

	// Update the article fields with the new values
	existingArticle.Title = payload.Title
	existingArticle.Slug = payload.Slug
	existingArticle.Author = payload.Author
	existingArticle.Content = payload.Content
	existingArticle.PictureUrl = payload.PictureUrl
	existingArticle.UpdatedAt = time.Now()

	// Save the updated article in the database
	result = ac.DB.Save(&existingArticle)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	response := &models.ArticleResponse{
		ID:         existingArticle.ID,
		Title:      existingArticle.Title,
		Slug:       existingArticle.Slug,
		Author:     existingArticle.Author,
		Content:    existingArticle.Content,
		PictureUrl: existingArticle.PictureUrl,
		CreatedAt:  existingArticle.CreatedAt,
		UpdatedAt:  existingArticle.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"article": response})
}

func (ac *ArticleController) Delete(ctx *gin.Context) {
	var article models.Article

	result := ac.DB.Where("id = ?", ctx.Param("id")).Delete(&article)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
