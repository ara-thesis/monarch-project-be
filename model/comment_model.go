package model

type CommentModel struct {
	Comment  string `json:"comment" xml:"comment" form:"comment"`
	Place_Id string `json:"place_id" xml:"place_id" form:"place_id"`
}
