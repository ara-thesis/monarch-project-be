package controller

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/ara-thesis/monarch-project-be/helper"
	"github.com/ara-thesis/monarch-project-be/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type AccountHandler struct{}

func (u *AccountHandler) GetUserInfo(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func (u *AccountHandler) CreateUserPlaceManager(c *fiber.Ctx) error {

	role := "PLACE MANAGER"
	model := new(model.AccountModel)

	model.Name = c.FormValue("name")
	model.Username = c.FormValue("username")
	model.Email = c.FormValue("email")
	model.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(c.FormValue("password"))))
	model.Mobile = c.FormValue("mobile")

	cmdStr := fmt.Sprintf(
		`INSERT INTO %s(id, nameperson, username, useremail, userpassword, usermobile, userrole, created_at, updated_at)
	 	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`, tbname["account_userinfo"],
	)
	cmdErr := db.Command(cmdStr, uuid.New(), model.Name, model.Username, model.Email, model.Password, model.Mobile, role, time.Now(), time.Now())
	if cmdErr != nil {
		return resp.ServerError(c, "Failed to create user")
	}

	return resp.Created(c, "Place manager account created")
}

func (u *AccountHandler) CreateUserTourist(c *fiber.Ctx) error {

	role := "TOURIST"
	model := new(model.AccountModel)

	model.Name = c.FormValue("name")
	model.Username = c.FormValue("username")
	model.Email = c.FormValue("email")
	model.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(c.FormValue("password"))))
	model.Mobile = c.FormValue("mobile")

	cmdStr := fmt.Sprintf(
		`INSERT INTO %s(id, nameperson, username, useremail, userpassword, usermobile, userrole, created_at, updated_at)
	 	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`, tbname["account_userinfo"],
	)
	cmdErr := db.Command(cmdStr, uuid.New(), model.Name, model.Username, model.Email, model.Password, model.Mobile, role, time.Now(), time.Now())
	if cmdErr != nil {
		return resp.ServerError(c, "Failed to create user")
	}

	return resp.Created(c, "Tourist account created")
}

func (u *AccountHandler) UserLogin(c *fiber.Ctx) error {

	respMap := make([]interface{}, 0)

	model := new(model.AccountModel)
	model.Username = c.FormValue("username")
	model.Password = c.FormValue("password")

	password := fmt.Sprintf("%x", sha256.Sum256([]byte(model.Password)))

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE username = $1 AND userpassword = $2", tbname["account_userinfo"])
	resQy, resErr := db.Query(qyStr, model.Username, password)
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}
	if resQy[0] == nil {
		return resp.NotFound(c, "Wrong username or password")
	}

	token, err := jwthelper.GenerateToken(map[string]interface{}{
		"userId":   resQy[0].(map[string]interface{})["id"],
		"username": resQy[0].(map[string]interface{})["username"],
		"userRole": resQy[0].(map[string]interface{})["userrole"],
	})
	if err != nil {
		return resp.ServerError(c, "Problem generating token")
	}

	respMap = append(respMap, map[string]string{
		"token": token,
	})

	return resp.Success(c, respMap, "LOGIN SUCCESS")
}

func (u *AccountHandler) EditUser(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func (u *AccountHandler) EditUserAsAdmin(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func (u *AccountHandler) DeleteUser(c *fiber.Ctx) error {

	// check for permission
	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != "ADMIN" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	return c.SendString("Test")
}
