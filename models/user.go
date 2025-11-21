package models

import (
	"time"
)

type Role struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;index" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID    int       `gorm:"type:int;not null;index" json:"role_id"`
	FullName  string    `gorm:"type:varchar(255);not null" json:"full_name"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	RoleID   int    `json:"role_id"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   *int   `json:"role_id"`
}

type UpdateUserRequest struct {
	FullName *string `json:"full_name"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6"`
	RoleID   *int    `json:"role_id"`
}
