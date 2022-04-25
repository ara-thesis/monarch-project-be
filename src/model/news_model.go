package model

import (
	"time"
)

type NewsModel struct {
	// id           string    `json:"id"`
	Title        string    `json:"title" xml:"title" form:"title"`
	Article      string    `json:"article" xml:"article" form:"article"`
	Image        string    `json:"image" xml:"image" form:"image"`
	Status       bool      `json:"status" xml:"status" form:"status"`
	Draft_status bool      `json:"draft_status" xml:"draft_status" form:"draft_status"`
	Created_by   string    `json:"created_by" xml:"created_by" form:"created_by"`
	Created_at   time.Time `json:"created_at" xml:"created_at" form:"created_at"`
	Updated_by   string    `json:"updated_by" xml:"updated_by" form:"updated_by"`
	Updated_at   time.Time `json:"updated_at" xml:"updated_at" form:"updated_at"`
}
