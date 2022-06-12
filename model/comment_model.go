package model

type CommentModel struct {
	Comment  interface{} `json:"comment" xml:"comment" form:"comment"`
	Place_Id interface{} `json:"place_id" xml:"place_id" form:"place_id"`
}
