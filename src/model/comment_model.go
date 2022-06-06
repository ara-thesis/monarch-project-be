package model

type CommentModel struct {
	Comment    interface{} `json:"comment" xml:"comment" form:"comment"`
	Place_Id   interface{} `json:"place_id" xml:"place_id" form:"place_id"`
	Created_by interface{} `json:"created_by" xml:"created_by" form:"created_by"`
	Created_at interface{} `json:"created_at" xml:"created_at" form:"created_at"`
	Updated_by interface{} `json:"updated_by" xml:"updated_by" form:"updated_by"`
	Updated_at interface{} `json:"updated_at" xml:"updated_at" form:"updated_at"`
}
