package models

type ReportOrderResponse struct {
    ID              int    `json:"id"`
    Code            string `json:"code"` // TRV-001
    UserID          int    `json:"user_id"`
    UserName        string `json:"user_name"`
    DestinationName string `json:"destination_name"`
    StartDate       string `json:"start_date"`
    EndDate         string `json:"end_date"`
    PeopleCount     int    `json:"people_count"`
    TransportPrice  int    `json:"transport_price"`
    TotalPrice      int    `json:"total_price"`
    StatusName      string `json:"status_name"`
}

type IncomeReportResponse struct {
    TotalIncome int64 `json:"total_income"`
}