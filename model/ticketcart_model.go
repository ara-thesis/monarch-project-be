package model

type TicketCartModel struct {
	Ticket_id   string  `json:"ticket_id" xml:"ticket_id" form:"ticket_id"`
	Qty         int     `json:"qty" xml:"qty" form:"qty"`
	Total_price float64 `json:"total_price" xml:"total_price" form:"total_price"`
	Wait        bool    `json:"wait" xml:"wait" form:"wait"`
}
