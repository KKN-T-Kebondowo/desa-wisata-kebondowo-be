package controllers

import (
	"kebondowo/models"
	"net/http"

	"time"

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

	result := ac.DB.Find(&articles)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"articles": articles})
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

	now := time.Now()
	newArticle := models.Article{
		Title: payload.Title,
		Slug: payload.Slug,
		Author: payload.Author,
		PictureUrl: payload.PictureUrl,
		Content: payload.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newArticle)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	articleResponse := &models.ArticleResponse{
		ID: newArticle.ID,
		Title: newArticle.Title,
		Content: newArticle.Content,
		Slug: newArticle.Slug,
		Author: newArticle.Author,
		PictureUrl: newArticle.PictureUrl,
		CreatedAt: newArticle.CreatedAt,
		UpdatedAt: newArticle.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"article": articleResponse})
}

// Delete
func (ac *ArticleController) Delete(ctx *gin.Context) {
	var article models.Article

	result := ac.DB.Where("id = ?", ctx.Param("id")).Delete(&article)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"article": article})
}

// Update
func (ac *ArticleController) Update(ctx *gin.Context) {
	var payload *models.ArticleInput
	var article models.Article

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result := ac.DB.Where("id = ?", ctx.Param("id")).First(&article)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	article.Title = payload.Title
	article.Slug = payload.Slug
	article.Author = payload.Author
	article.PictureUrl = payload.PictureUrl
	article.Content = payload.Content
	article.UpdatedAt = time.Now()

	result = ac.DB.Save(&article)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
		return
	}

	articleResponse := &models.ArticleResponse{
		ID: article.ID,
		Title: article.Title,
		Content: article.Content,
		Slug: article.Slug,
		PictureUrl: article.PictureUrl,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"article": articleResponse})
}

