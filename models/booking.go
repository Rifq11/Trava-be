package models

type Booking struct {
	ID               int    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int    `gorm:"not null;index" json:"user_id"`
	DestinationID    int    `gorm:"not null;index" json:"destination_id"`
	TransportationID int    `gorm:"not null;index" json:"transportation_id"`
	PaymentMethodID  int    `gorm:"not null;index" json:"payment_method_id"`
	StatusID         int    `gorm:"not null;index" json:"status_id"`
	PeopleCount      int    `gorm:"not null" json:"people_count"`
	StartDate        string `gorm:"type:datetime;not null" json:"start_date"`
	EndDate          string `gorm:"type:datetime;not null" json:"end_date"`
	TransportPrice   int    `gorm:"not null" json:"transport_price"`
	DestinationPrice int    `gorm:"not null" json:"destination_price"`
	TotalPrice       int    `gorm:"not null" json:"total_price"`
}

type BookingResponse struct {
	BookingID         int    `json:"booking_id"`
	DestinationID     int    `json:"destination_id"`
	DestinationName   string `json:"destination_name"`
	Location          string `json:"location"`
	PeopleCount       int    `json:"people_count"`
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
	TotalPrice        int    `json:"total_price"`
	TransportPrice    int    `json:"transport_price"`
	DestinationPrice  int    `json:"destination_price"`
	StatusName        string `json:"status_name"`
	PaymentMethodName string `json:"payment_method_name"`
	DestinationImage  string `json:"destination_image"`
}

type CreateBookingRequest struct {
	DestinationID    int    `json:"destination_id" binding:"required"`
	TransportationID int    `json:"transportation_id" binding:"required"`
	PaymentMethodID  int    `json:"payment_method_id" binding:"required"`
	PeopleCount      int    `json:"people_count" binding:"required"`
	StartDate        string `json:"start_date" binding:"required"`
	EndDate          string `json:"end_date" binding:"required"`
}

type BookingStatus struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(255);not null;index" json:"name"`
}

type PaymentMethod struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(255);not null;index" json:"name"`
}

type AdminBookingResponse struct {
	Code             string `json:"code"`
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	UserName         string `json:"user_name"`
	UserEmail        string `json:"user_email"`
	DestinationName  string `json:"destination_name"`
	StatusName       string `json:"status_name"`
	PaymentMethod    string `json:"payment_method"`
	PeopleCount      int    `json:"people_count"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	TotalPrice       int    `json:"total_price"`
	TransportPrice   int    `json:"transport_price"`
	DestinationPrice int    `json:"destination_price"`
}

type AdminBookingDetailResponse struct {
	Code             string `json:"code"`
	ID               int    `json:"id"`
	UserName         string `json:"user_name"`
	DestinationFull  string `json:"destination_full"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	TransportName    string `json:"transport_name"`
	PeopleCount      int    `json:"people_count"`
	TotalPrice       int    `json:"total_price"`
	PaymentMethod    string `json:"payment_method"`
	DestinationImage string `json:"destination_image"`
}
