package model

type TicketModel struct {
	Ticket_name    interface{} `json:"ticket_name" xml:"ticket_name" form:"ticket_name"`
	Ticket_details interface{} `json:"ticket_details" xml:"ticket_details" form:"ticket_details"`
	Ticket_placeid interface{} `json:"ticket_placeid" xml:"ticket_placeid" form:"ticket_placeid"`
	Ticket_price   interface{} `json:"ticket_price" xml:"ticket_price" form:"ticket_price"`
}
