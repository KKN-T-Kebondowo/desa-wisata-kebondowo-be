package models

import (
	"time"
)

type Tourism struct {
	ID   uint32 `gorm:"primary_key"`
	Title string `gorm:"type:varchar(255);not null"`
	Slug string `gorm:"type:varchar(255);not null"`
	CoverPictureUrl string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type TourismInput struct {
	Title string `json:"title" binding:"required"`
	Slug string `json:"slug" binding:"required"`
	CoverPictureUrl string `json:"cover_picture_url" binding:"required"`
	Description string `json:"description" binding:"required"`
	
}

type TourismResponse struct {
	ID   uint32 `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Slug string `json:"slug,omitempty"`
	CoverPictureUrl string `json:"cover_picture_url,omitempty"`
	Description string `json:"description,omitempty"`
	Pictures []TourismPicture `json:"pictures,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
