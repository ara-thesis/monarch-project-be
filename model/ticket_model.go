package model

type TicketModel struct {
	Ticket_name    string `json:"ticket_name" xml:"ticket_name" form:"ticket_name"`
	Ticket_details string `json:"ticket_details" xml:"ticket_details" form:"ticket_details"`
	Ticket_placeid string `json:"ticket_placeid" xml:"ticket_placeid" form:"ticket_placeid"`
	Ticket_price   int    `json:"ticket_price" xml:"ticket_price" form:"ticket_price"`
}

type TicketBoughtModel struct {
	TicketBought_ticketId string
	TicketBought_code     string `json:"code" xml:"code" form:"code"`
	TicketBought_qty      int
	TicketBought_redeemed bool
	TicketBought_owner    string
}
