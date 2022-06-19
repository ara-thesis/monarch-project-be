package model

type NewsModel struct {
	Title        string      `json:"title" xml:"title" form:"title"`
	Article      string      `json:"article" xml:"article" form:"article"`
	Image        interface{} `json:"image" xml:"image" form:"image"`
	Status       bool        `json:"status" xml:"status" form:"status"`
	Draft_status bool        `json:"draft_status" xml:"draft_status" form:"draft_status"`
}
