package model

type AccountModel struct {
	Name     string `json:"name" xml:"name" form:"name"`
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
	Mobile   string `json:"mobile" xml:"mobile" form:"mobile"`
	Role     string `json:"role" xml:"role" form:"role"`
}
