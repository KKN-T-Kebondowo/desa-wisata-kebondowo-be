package models

type TourismPicture struct {
	ID         uint32  `gorm:"column:id;primary_key"`
	PictureUrl string  `gorm:"column:picture_url;type:varchar(255);not null"`
	TourismID  uint32  `gorm:"column:tourism_id;not null;references:ID"`
	Tourism    Tourism `gorm:"foreignKey:TourismID"`
}

type TourismPictureInput struct {
	PictureUrl string `json:"picture_url" binding:"required"`
	TourismID  uint32 `json:"tourism_id" binding:"required"`
}

type TourismPictureResponse struct {
	ID         uint32 `json:"id,omitempty"`
	PictureUrl string `json:"picture_url,omitempty"`
	TourismID  uint32 `json:"tourism_id,omitempty"`
}
