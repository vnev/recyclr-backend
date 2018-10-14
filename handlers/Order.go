package handlers

//Order : struct to hold order information
type Order struct {
	ID        int  `json:"order_id"`
	UserID    int  `json:"user_id"`
	CompanyID int  `json:"company_id"`
	Total     int  `json:"total"`
	Confirmed bool `json:"confirmed"`
}
