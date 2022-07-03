package model

import (
	"time"
)

type PlaceInfoModel struct {
	Place_name      string    `json:"place_name" xml:"place_name" form:"place_name"`
	Place_info      string    `json:"place_info" xml:"place_info" form:"place_info"`
	Place_street    string    `json:"place_street" xml:"place_street" form:"place_street"`
	Place_city      string    `json:"place_city" xml:"place_city" form:"place_city"`
	Place_stateprov string    `json:"place_stateprov" xml:"place_stateprov" form:"place_stateprov"`
	Place_country   string    `json:"place_country" xml:"place_country" form:"place_country"`
	Place_postal    string    `json:"place_postal" xml:"place_postal" form:"place_postal"`
	Place_loc_long  float32   `json:"place_loc_long" xml:"place_loc_long" form:"place_loc_long"`
	Place_loc_lat   float32   `json:"place_loc_lat" xml:"place_loc_lat" form:"place_loc_lat"`
	Place_images    []string  ``
	Place_opentime  time.Time `json:"place_opentime" xml:"place_opentime" form:"place_opentime"`
	Place_closetime time.Time `json:"place_closetime" xml:"place_closetime" form:"place_closetime"`
}
