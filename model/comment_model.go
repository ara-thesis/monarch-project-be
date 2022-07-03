package model

type CommentModel struct {
	Comment  string `json:"comment" xml:"comment" form:"comment"`
	Score    int    `json:"score" xml:"score" form:"score"`
	Place_Id string `json:"place_id" xml:"place_id" form:"place_id"`
}
