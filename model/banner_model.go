package model

type BannerModel struct {
	Title  interface{} `json:"title" xml:"title" form:"title"`
	Detail interface{} `json:"detail" xml:"detail" form:"detail"`
	Image  interface{} `json:"image" xml:"image" form:"image"`
	Status interface{} `json:"status" xml:"status" form:"status"`
}
