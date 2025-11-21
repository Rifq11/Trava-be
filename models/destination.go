package models

type Destination struct {
	ID             int    `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID     int    `gorm:"not null;index" json:"category_id"`
	CreatedBy      int    `gorm:"not null;index" json:"created_by"`
	Name           string `gorm:"type:varchar(255);not null" json:"name"`
	Description    string `gorm:"type:text" json:"description"`
	Location       string `gorm:"type:varchar(255)" json:"location"`
	PricePerPerson int    `gorm:"not null" json:"price_per_person"`
	Image          string `gorm:"type:text" json:"image"`
}

type DestinationResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Location       string `json:"location"`
	PricePerPerson int    `json:"price_per_person"`
	Image          string `json:"image"`
	CategoryID     int    `json:"category_id"`
	CreatedBy      int    `json:"created_by"`
}

type DestinationDetailResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Location       string `json:"location"`
	PricePerPerson int    `json:"price_per_person"`
	Image          string `json:"image"`
	CategoryID     int    `json:"category_id"`
	CreatedBy      int    `json:"created_by"`
}

type CreateDestinationRequest struct {
	CategoryID     int    `json:"category_id" form:"category_id" binding:"required"`
	Name           string `json:"name" form:"name" binding:"required"`
	Description    string `json:"description" form:"description"`
	Location       string `json:"location" form:"location" binding:"required"`
	PricePerPerson int    `json:"price_per_person" form:"price_per_person" binding:"required"`
	Image          string `json:"image" form:"image"`
}

type UpdateDestinationRequest struct {
	CategoryID     *int    `json:"category_id" form:"category_id"`
	Name           *string `json:"name" form:"name"`
	Description    *string `json:"description" form:"description"`
	Location       *string `json:"location" form:"location"`
	PricePerPerson *int    `json:"price_per_person" form:"price_per_person"`
	Image          *string `json:"image" form:"image"`
}

type DestinationCategory struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(255);not null;index" json:"name"`
}
