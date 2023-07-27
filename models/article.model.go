package models

import (
	"time"
)

type Article struct {
	ID         uint32    `gorm:"column:id;primary_key" json:"id"`
	Title      string    `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Slug       string    `gorm:"column:slug;type:varchar(255);not null" json:"slug"`
	Author     string    `gorm:"column:author;type:varchar(255);not null" json:"author"`
	Content    string    `gorm:"column:content;type:text;not null" json:"content"`
	PictureUrl string    `gorm:"column:picture_url;type:varchar(255);not null" json:"picture_url"`
	Visitor    uint32    `gorm:"column:visitor;not null;default:0" json:"visitor"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

type ArticleInput struct {
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug" binding:"required"`
	Author     string `json:"author" binding:"required"`
	Content    string `json:"content" binding:"required"`
	PictureUrl string `json:"picture_url" binding:"required"`
}

type ArticleUpdate struct {
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug" binding:"required"`
	Author     string `json:"author" binding:"required"`
	Content    string `json:"content" binding:"required"`
	PictureUrl string `json:"picture_url"`
}

type ArticleResponse struct {
	ID         uint32    `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Slug       string    `json:"slug,omitempty"`
	Author     string    `json:"author,omitempty"`
	Content    string    `json:"content,omitempty"`
	PictureUrl string    `json:"picture_url,omitempty"`
	Visitor    uint32    `json:"visitor,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
