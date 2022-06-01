package model

type BannerModel struct {
	Title      interface{} `json:"title" xml:"title" form:"title"`
	Detail     interface{} `json:"detail" xml:"detail" form:"detail"`
	Image      interface{} `json:"image" xml:"image" form:"image"`
	Status     interface{} `json:"status" xml:"status" form:"status"`
	Created_by interface{} `json:"created_by" xml:"created_by" form:"created_by"`
	Created_at interface{} `json:"created_at" xml:"created_at" form:"created_at"`
	Updated_by interface{} `json:"updated_by" xml:"updated_by" form:"updated_by"`
	Updated_at interface{} `json:"updated_at" xml:"updated_at" form:"updated_at"`
}
