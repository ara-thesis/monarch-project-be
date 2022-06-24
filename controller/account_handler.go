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

type AccountHandler struct {
	Tbname string
}

//////////////////////////////
// fetch personal user info
//////////////////////////////
func (u *AccountHandler) GetUserInfo(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	qyStr := fmt.Sprintf("SELECT id, nameperson, username, useremail, profilepics FROM %s WHERE id = $1", u.Tbname)
	resQy, resErr := db.Query(qyStr, userData.UserId)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")
}

//////////////////////////////////
// fetch place manager user info
//////////////////////////////////
func (u *AccountHandler) GetUserListPM(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != "ADMIN" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	qyStr := fmt.Sprintf("SELECT id, nameperson, username, useremail, profilepics FROM %s WHERE userrole = $1", u.Tbname)
	resQy, resErr := db.Query(qyStr, "PLACE MANAGER")

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")
}

////////////////////////////
// fetch tourist user info
////////////////////////////
func (u *AccountHandler) GetUserListTourist(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != "ADMIN" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	qyStr := fmt.Sprintf("SELECT id, nameperson, username, useremail, profilepics FROM %s WHERE userrole = $1", u.Tbname)
	resQy, resErr := db.Query(qyStr, "TOURIST")

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")
}

///////////////////////////////
// create user admin by admin
///////////////////////////////
func (u *AccountHandler) CreateUserAdmin(c *fiber.Ctx) error {

	role := "ADMIN"
	model := new(model.AccountModel)

	model.Name = c.FormValue("name")
	model.Username = c.FormValue("username")
	model.Email = c.FormValue("email")
	model.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(c.FormValue("password"))))
	model.Mobile = c.FormValue("mobile")

	cmdStr := fmt.Sprintf(
		`INSERT INTO %s(id, nameperson, username, useremail, userpassword, usermobile, profilepics, userrole, created_at, updated_at)
	 	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, u.Tbname,
	)
	cmdErr := db.Command(cmdStr, uuid.New(), model.Name, model.Username, model.Email,
		model.Password, model.Mobile, "/api/public/dummy/dummy-pics.jpg", role, time.Now(), time.Now())
	if cmdErr != nil {
		return resp.ServerError(c, "Failed to create user")
	}

	return resp.Created(c, "Place manager account created")
}

///////////////////////////////////////
// create user place manager by admin
///////////////////////////////////////
func (u *AccountHandler) CreateUserPlaceManager(c *fiber.Ctx) error {

	role := "PLACE MANAGER"
	model := new(model.AccountModel)

	model.Name = c.FormValue("name")
	model.Username = c.FormValue("username")
	model.Email = c.FormValue("email")
	model.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(c.FormValue("password"))))
	model.Mobile = c.FormValue("mobile")

	cmdStr := fmt.Sprintf(
		`INSERT INTO %s(id, nameperson, username, useremail, userpassword, usermobile, profilepics, userrole, created_at, updated_at)
	 	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, u.Tbname,
	)
	cmdErr := db.Command(cmdStr, uuid.New(), model.Name, model.Username, model.Email,
		model.Password, model.Mobile, "/api/public/dummy/dummy-pics.jpg", role, time.Now(), time.Now())
	if cmdErr != nil {
		return resp.ServerError(c, "Failed to create user")
	}

	return resp.Created(c, "Place manager account created")
}

////////////////////////
// create user tourist
////////////////////////
func (u *AccountHandler) CreateUserTourist(c *fiber.Ctx) error {

	role := "TOURIST"
	model := new(model.AccountModel)

	model.Name = c.FormValue("name")
	model.Username = c.FormValue("username")
	model.Email = c.FormValue("email")
	model.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(c.FormValue("password"))))
	model.Mobile = c.FormValue("mobile")

	cmdStr := fmt.Sprintf(
		`INSERT INTO %s(id, nameperson, username, useremail, userpassword, usermobile, profilepics, userrole, created_at, updated_at)
	 	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, u.Tbname,
	)
	cmdErr := db.Command(cmdStr, uuid.New(), model.Name, model.Username, model.Email,
		model.Password, model.Mobile, "/api/public/dummy/dummy-pics.jpg", role, time.Now(), time.Now())
	if cmdErr != nil {
		return resp.ServerError(c, "Failed to create user")
	}

	return resp.Created(c, "Tourist account created")
}

///////////////
// user login
///////////////
func (u *AccountHandler) UserLogin(c *fiber.Ctx) error {

	respMap := make([]interface{}, 0)

	model := new(model.AccountModel)
	model.Username = c.FormValue("username")
	model.Password = c.FormValue("password")

	password := fmt.Sprintf("%x", sha256.Sum256([]byte(model.Password)))

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE username = $1 AND userpassword = $2", u.Tbname)
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

//////////////////////////////
// update personal user info
//////////////////////////////
func (u *AccountHandler) EditUser(c *fiber.Ctx) error {
	return c.SendString("Test")
}

//////////////////////////////
// update user info by admin
//////////////////////////////
func (u *AccountHandler) EditUserAsAdmin(c *fiber.Ctx) error {
	return c.SendString("Test")
}

//////////////////////////////
// delete user info by admin
//////////////////////////////
func (u *AccountHandler) DeleteUser(c *fiber.Ctx) error {

	// check for permission
	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != "ADMIN" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// delete data process
	cmdStr := fmt.Sprintf("DELETE FROM %s WHERE id = $1", u.Tbname)
	resErr := db.Command(cmdStr, c.Params("id"))
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "DELETE SUCCESS")
}
