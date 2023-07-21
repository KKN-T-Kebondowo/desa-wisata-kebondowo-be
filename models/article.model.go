package models

import (
	"time"
)

type Article struct {
	ID         uint32    `gorm:"column:id;primary_key"`
	Title      string    `gorm:"column:title;type:varchar(255);not null"`
	Slug       string    `gorm:"column:slug;type:varchar(255);not null"`
	Author     string    `gorm:"column:author;type:varchar(255);not null"`
	Content    string    `gorm:"column:content;type:varchar(255);not null"`
	PictureURL string    `gorm:"column:picture_url;type:varchar(255);not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`
}

type ArticleInput struct {
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug" binding:"required"`
	Author     string `json:"author" binding:"required"`
	Content    string `json:"content" binding:"required"`
	PictureURL string `json:"picture_url" binding:"required"`
}

type ArticleResponse struct {
	ID         uint32    `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Slug       string    `json:"slug,omitempty"`
	Author     string    `json:"author,omitempty"`
	Content    string    `json:"content,omitempty"`
	PictureURL string    `json:"picture_url,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
