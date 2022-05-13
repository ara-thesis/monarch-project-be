package model

type PlaceInfoModel struct {
	Place_name      interface{} `json:"place_name" xml:"place_name" form:"place_name"`
	Place_info      interface{} `json:"place_info" xml:"place_info" form:"place_info"`
	Place_city      interface{} `json:"place_city" xml:"place_city" form:"place_city"`
	Place_stateprov interface{} `json:"place_stateprov" xml:"place_stateprov" form:"place_stateprov"`
	Place_country   interface{} `json:"place_country" xml:"place_country" form:"place_country"`
	Place_postal    interface{} `json:"place_postal" xml:"place_postal" form:"place_postal"`
	Place_loc       interface{} `json:"place_loc" xml:"place_loc" form:"place_loc"`
	Place_opentime  interface{} `json:"place_opentime" xml:"place_opentime" form:"place_opentime"`
	Place_closetime interface{} `json:"place_closetime" xml:"place_closetime" form:"place_closetime"`
}