package model

type BookingModel struct {
	UserId  string
	PlaceId string `json:"placeid" xml:"placeid" form:"placeid"`
}
