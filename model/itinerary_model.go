package model

import (
	"github.com/google/uuid"
)

type ItineraryModel struct {
	Title  string        `json:"title" xml:"title" form:"title"`
	Detail string        `json:"detail" xml:"detail" form:"detail"`
	Items  []interface{} `json:"items" xml:"items" form:"items"`
}

type ItineraryDayModel struct {
	Went_time string
	Place_loc []interface{}
}

type ItineraryItemModel struct {
	ItineraryId uuid.UUID
	PlaceId     string
	Detail      string
	Day         string
	In_time     string
	Out_time    string
}
