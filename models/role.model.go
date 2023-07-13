package models

type Role struct {
	ID   uint32 `gorm:"primary_key"`
	Name string `gorm:"type:varchar(255);not null"`
}

type RoleInput struct {
	Name string `json:"name" binding:"required"`
}
