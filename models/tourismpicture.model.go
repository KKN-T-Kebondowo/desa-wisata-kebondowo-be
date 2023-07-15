package models

type TourismPicture struct {
	ID         uint32  `gorm:"primary_key"`
	PictureUrl string  `gorm:"type:varchar(255);not null"`
	Latitude   float64 `gorm:"not null"`
	Longitude  float64 `gorm:"not null"`
	TourismID  uint32  `gorm:"not null;references:ID"`
	Tourism    Tourism `gorm:"foreignKey:TourismID"`
}

type TourismPictureInput struct {
	PictureUrl string  `json:"picture_url" binding:"required"`
	Latitude   float64 `json:"latitude" binding:"required"`
	Longitude  float64 `json:"longitude" binding:"required"`
	TourismID  uint32  `json:"tourism_id" binding:"required"`
}
