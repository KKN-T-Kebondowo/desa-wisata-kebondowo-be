package models

import "time"


type UMKMPicture struct {
	ID         uint32    `gorm:"column:id;primary_key"`
	PictureUrl string    `gorm:"column:picture_url;type:varchar(255);not null"`
	Caption    string    `gorm:"column:caption;type:varchar(255)"`
	UMKMID     uint32    `gorm:"column:umkm_id;not null;references:ID"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`
	UMKM       UMKM      `gorm:"foreignKey:UMKMID"`
}

type UMKMPictureInput struct {
	PictureUrl string `json:"picture_url" binding:"required"`
	Caption    string `json:"caption"`
	UMKMID     uint32 `json:"umkm_id" binding:"required"`
}

type UMKMPictureResponse struct {
	ID         uint32    `json:"id,omitempty"`
	PictureUrl string    `json:"picture_url,omitempty"`
	Caption    string    `json:"caption,omitempty"`
	UMKMID     uint32    `json:"umkm_id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
