package controller

import "github.com/ara-thesis/monarch-project-be/src/helper"

var (
	tbname = map[string]string{
		"news":             "newstb",
		"account_userinfo": "userinfo",
		"account_roleuser": "roleuser",
		"account_roleinfo": "roleinfo",
	}
	db        = new(helper.PgHelper)
	resp      = new(helper.ResponseHelper)
	jwthelper = new(helper.JwtHelper)
)
