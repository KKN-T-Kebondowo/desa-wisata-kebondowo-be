package models

import "time"

type TourismPicture struct {
	ID         uint32    `gorm:"column:id;primary_key"`
	PictureUrl string    `gorm:"column:picture_url;type:varchar(255);not null"`
	Caption    string    `gorm:"column:caption;type:varchar(255)"`
	TourismID  uint32    `gorm:"column:tourism_id;not null;references:ID"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`
	Tourism    Tourism   `gorm:"foreignKey:TourismID"`
}

type TourismPictureInput struct {
	PictureUrl string `json:"picture_url" binding:"required"`
	Caption    string `json:"caption"`
	TourismID  uint32 `json:"tourism_id" binding:"required"`
}

type TourismPictureResponse struct {
	ID         uint32    `json:"id,omitempty"`
	PictureUrl string    `json:"picture_url,omitempty"`
	Caption    string    `json:"caption,omitempty"`
	TourismID  uint32    `json:"tourism_id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
