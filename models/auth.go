package models

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   *int   `json:"role_id"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	RoleID   int    `json:"role_id"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type RegisterResponse struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	RoleID   int    `json:"role_id"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type UpdateProfileRequest struct {
	FullName  *string `json:"full_name"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	BirthDate *string `json:"birth_date"`
	Password  *string `json:"password" binding:"omitempty,min=6"`
	UserPhoto *string `json:"user_photo"`
}

type UserProfile struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int    `gorm:"not null;index" json:"user_id"`
	Phone       string `gorm:"type:varchar(50)" json:"phone"`
	Address     string `gorm:"type:text" json:"address"`
	BirthDate   string `gorm:"type:date" json:"birth_date"`
	UserPhoto   string `gorm:"type:varchar(500)" json:"user_photo"`
	IsCompleted bool   `gorm:"default:false" json:"is_completed"`
}
