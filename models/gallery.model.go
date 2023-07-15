package models

import (
	"time"
)

//Role to foreign key to Role table

type Gallery struct {
	ID         uint32 `gorm:"primary_key"`
	PictureUrl string `gorm:"type:varchar(255);not null"`
	Caption    string `gorm:"type:varchar(255);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
