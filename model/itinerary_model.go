package model

import "time"

type ItineraryModel struct {
	Title  string        `json:"title" xml:"title" form:"title"`
	Detail string        `json:"detail" xml:"detail" form:"detail"`
	Items  []interface{} `json:"items" xml:"items" form:"items"`
}

type ItineraryItemModel struct {
	ItineraryId string
	PlaceId     string
	Detail      string
	Went_time   time.Time
}
