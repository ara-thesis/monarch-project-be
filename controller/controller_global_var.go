package controller

import "github.com/ara-thesis/monarch-project-be/helper"

type minioCred struct {
	endpoint        string
	accessKeyId     string
	secretAccessKey string
	useSSL          bool
}

var (
	tbname = map[string]string{
		"news":             "newstb",
		"account_userinfo": "userinfo",
		"placeinfo":        "placeinfotb",
		"banner":           "bannertb",
		"review":           "reviewtb",
		"ticket":           "tickettb",
	}
	db        = new(helper.PgHelper)
	resp      = new(helper.ResponseHelper)
	jwthelper = new(helper.JwtHelper)
	roleId    = map[string]string{
		"T":   "TOURIST",
		"PM":  "PLACE MANAGER",
		"ADM": "ADMIN",
	}
)
