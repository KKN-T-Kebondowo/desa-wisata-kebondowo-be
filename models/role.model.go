package models

type Role struct {
	ID   uint32 `gorm:"column:id;primary_key"`
	Name string `gorm:"column:name;type:varchar(255);not null"`
}

type RoleInput struct {
	Name string `json:"name" binding:"required"`
}
