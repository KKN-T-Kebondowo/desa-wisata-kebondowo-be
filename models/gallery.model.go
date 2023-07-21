package models

import (
	"time"
)

type Gallery struct {
	ID         uint32    `gorm:"column:id;primary_key"`
	PictureUrl string    `gorm:"column:picture_url;type:varchar(255);not null"`
	Caption    string    `gorm:"column:caption;type:varchar(255);not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`
}

type GalleryInput struct {
	PictureUrl string `json:"picture_url" binding:"required"`
	Caption    string `json:"caption" binding:"required"`
}

type GalleryResponse struct {
	ID         uint32    `json:"id,omitempty"`
	PictureUrl string    `json:"picture_url,omitempty"`
	Caption    string    `json:"caption,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
