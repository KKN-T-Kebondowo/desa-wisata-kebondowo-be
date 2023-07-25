package models

import (
	"time"
)

type UMKM struct {
	ID              uint32    `gorm:"column:id;primary_key"`
	Title           string    `gorm:"column:title;type:varchar(255);not null"`
	Slug            string    `gorm:"column:slug;type:varchar(255);not null"`
	CoverPictureUrl string    `gorm:"column:cover_picture_url;type:varchar(255);not null"`
	Description     string    `gorm:"column:description;type:varchar(255);not null"`
	Latitude        float64   `gorm:"column:latitude;type:float;not null"`
	Longitude       float64   `gorm:"column:longitude;type:float;not null"`
	Visitor         uint32    `gorm:"column:visitor;not null;default:0"`
	Contact         string    `gorm:"column:contact;type:varchar(255);not null"`
	ContactName     string    `gorm:"column:contact_name;type:varchar(255);not null"`
	CreatedAt       time.Time `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null"`
}

type UMKMInput struct {
	Title           string             `json:"title" binding:"required"`
	Slug            string             `json:"slug" binding:"required"`
	Latitude        float64            `json:"latitude" binding:"required"`
	Longitude       float64            `json:"longitude" binding:"required"`
	CoverPictureUrl string             `json:"cover_picture_url" binding:"required"`
	Description     string             `json:"description" binding:"required"`
	Contact         string             `json:"contact" binding:"required"`
	ContactName     string             `json:"contact_name" binding:"required"`
	Pictures        []UMKMPictureInput `json:"pictures"`
}

type UMKMUpdate struct {
	Title           string  `json:"title" binding:"required"`
	Slug            string  `json:"slug" binding:"required"`
	Latitude        float64 `json:"latitude" binding:"required"`
	Longitude       float64 `json:"longitude" binding:"required"`
	CoverPictureUrl string  `json:"cover_picture_url"`
	Description     string  `json:"description" binding:"required"`
	Contact         string  `json:"contact" binding:"required"`
	ContactName     string  `json:"contact_name" binding:"required"`
}

type UMKMResponse struct {
	ID              uint32                `json:"id,omitempty"`
	Title           string                `json:"title,omitempty"`
	Slug            string                `json:"slug,omitempty"`
	CoverPictureUrl string                `json:"cover_picture_url,omitempty"`
	Description     string                `json:"description,omitempty"`
	Latitude        float64               `json:"latitude,omitempty"`
	Longitude       float64               `json:"longitude,omitempty"`
	Visitor         uint32                `json:"visitor,omitempty"`
	Contact         string                `json:"contact,omitempty"`
	ContactName     string                `json:"contact_name,omitempty"`
	Pictures        []UMKMPictureResponse `json:"pictures,omitempty"`
	CreatedAt       time.Time             `json:"created_at,omitempty"`
	UpdatedAt       time.Time             `json:"updated_at,omitempty"`
}
