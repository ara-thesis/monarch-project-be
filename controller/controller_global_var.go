package controller

import "github.com/ara-thesis/monarch-project-be/helper"

type roleSet struct {
	t   string
	pm  string
	adm string
}

var (
	db        = new(helper.PgHelper)
	resp      = new(helper.ResponseHelper)
	jwthelper = new(helper.JwtHelper)
	// roleId    = map[string]string{
	// 	"T":   "TOURIST",
	// 	"PM":  "PLACE MANAGER",
	// 	"ADM": "ADMIN",
	// }
	roleId = roleSet{
		t:   "TOURIST",
		pm:  "PLACE MANAGER",
		adm: "ADMIN",
	}
)
