package models

import (
	"time"
)

type Gallery struct {
	ID         uint32    `gorm:"column:id;primary_key" json:"id"`
	PictureUrl string    `gorm:"column:picture_url;type:varchar(255);not null" json:"picture_url"`
	Caption    string    `gorm:"column:caption;type:varchar(255);not null" json:"caption"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
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
