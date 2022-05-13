package controller

import "github.com/ara-thesis/monarch-project-be/src/helper"

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
		"account_roleuser": "roleuser",
		"account_roleinfo": "roleinfo",
		"placeinfo":        "placeinfotb",
	}
	db               = new(helper.PgHelper)
	resp             = new(helper.ResponseHelper)
	jwthelper        = new(helper.JwtHelper)
	miniocredentials = minioCred{
		"localhost:9000",
		"YV2n9zZ1rSqDbNZF",
		"V1ynQ8Xx6RTMsHx9rDTZFNp61pdWDHFM",
		false,
	}
)
