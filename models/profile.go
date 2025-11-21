package models

type ProfileUserResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
}

type ProfileResponse struct {
	User    ProfileUserResponse `json:"user"`
	Profile interface{}         `json:"profile"` // Can be UserProfile or AdminProfile
}

type CompleteProfileRequest struct {
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	BirthDate *string `json:"birth_date"`
	UserPhoto *string `json:"user_photo"`
}

type ProfileDetailResponse struct {
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	BirthDate string `json:"birth_date"`
	UserPhoto string `json:"user_photo"`
	Password  string `json:"password"`
	RoleID    int    `json:"role_id"`
}

type AdminProfile struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int    `gorm:"not null;index" json:"user_id"`
	Phone       string `gorm:"type:varchar(50)" json:"phone"`
	Address     string `gorm:"type:text" json:"address"`
	BirthDate   string `gorm:"type:date" json:"birth_date"`
	UserPhoto   string `gorm:"type:varchar(500)" json:"user_photo"`
	IsCompleted bool   `gorm:"default:false" json:"is_completed"`
}
