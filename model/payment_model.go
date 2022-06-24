package model

type PaymentModel struct {
	Total_price float64
	Image       interface{} `json:"image" xml:"image" form:"image"`
}
