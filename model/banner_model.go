package model

type BannerModel struct {
	Title  string      `json:"title" xml:"title" form:"title"`
	Detail string      `json:"detail" xml:"detail" form:"detail"`
	Image  interface{} //`json:"image" xml:"image" form:"image"`
	Status bool        `json:"status" xml:"status" form:"status"`
}
