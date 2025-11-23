package models

type DashboardStatisticsResponse struct {
	TotalDestinations   int64 `json:"total_destinations"`
	TotalActiveOrders   int64 `json:"total_active_orders"`
	TotalRegisteredUser int64 `json:"total_registered_user"`
}

type MonthlySalesResponse struct {
	Month   int `json:"month"`
	Revenue int `json:"revenue"`
}