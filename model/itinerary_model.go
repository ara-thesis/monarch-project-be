package model

type ItineraryModel struct {
	Title  string      `json:"title" xml:"title" form:"title"`
	Detail string      `json:"detail" xml:"detail" form:"detail"`
	Items  interface{} `json:"items" xml:"items" form:"items"`
}
