package models

import (
	"time"
)

type Article struct {
	ID   uint32 `gorm:"primary_key"`
	Title string `gorm:"type:varchar(255);not null"`
	Slug string `gorm:"type:varchar(255);not null"`
	Author string `gorm:"type:varchar(255);not null"`
	Content string `gorm:"type:varchar(255);not null"`
	PictureUrl string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type ArticleInput struct {
	Title string `json:"title" binding:"required"`
	Slug string `json:"slug" binding:"required"`
	Author string `json:"author" binding:"required"`
	Content string `json:"content" binding:"required"`
	PictureUrl string `json:"picture_url" binding:"required"`

}

type ArticleResponse struct {
	ID   uint32 `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Slug string `json:"slug,omitempty"`
	Author string `json:"author,omitempty"`
	Content string `json:"content,omitempty"`
	PictureUrl string `json:"picture_url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}