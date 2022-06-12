package model

type NewsModel struct {
	Title        interface{} `json:"title" xml:"title" form:"title"`
	Article      interface{} `json:"article" xml:"article" form:"article"`
	Image        interface{} `json:"image" xml:"image" form:"image"`
	Status       interface{} `json:"status" xml:"status" form:"status"`
	Draft_status interface{} `json:"draft_status" xml:"draft_status" form:"draft_status"`
}
