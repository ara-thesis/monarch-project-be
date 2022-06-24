package controller

import "github.com/ara-thesis/monarch-project-be/helper"

var (
	db        = new(helper.PgHelper)
	resp      = new(helper.ResponseHelper)
	jwthelper = new(helper.JwtHelper)
	roleId    = map[string]string{
		"T":   "TOURIST",
		"PM":  "PLACE MANAGER",
		"ADM": "ADMIN",
	}
)
