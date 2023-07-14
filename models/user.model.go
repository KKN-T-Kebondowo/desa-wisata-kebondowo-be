package models

import (
	"time"
)

//Role to foreign key to Role table

type User struct {
	ID        uint32 `gorm:"primary_key"`
	Username  string    `gorm:"type:varchar(255);not null"`
	Password  string    `gorm:"not null"`
	RoleID  uint32    `gorm:"not null;references:ID"`
	
	CreatedAt time.Time
	UpdatedAt time.Time
	Role 	Role `gorm:"foreignKey:RoleID"`
}

type SignUpInput struct {
	Username            string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	RoleID	  uint32 `json:"roleid" binding:"required"`
}

type SignInInput struct {
	Username    string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        uint32 `json:"id,omitempty"`
	Username      string    `json:"username,omitempty"`
	RoleID      uint32   `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

