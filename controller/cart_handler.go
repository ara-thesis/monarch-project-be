package controller

import (
	"fmt"
	"time"

	"github.com/ara-thesis/monarch-project-be/helper"
	"github.com/ara-thesis/monarch-project-be/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CartHandler struct {
	Tbname        string
	Tbname_ticket string
}

func (ch *CartHandler) GetCart(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != roleId.t {
		return resp.Forbidden(c, "Access Forbidden")
	}

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE created_by = $1 AND wait = $2 ORDER BY updated_at DESC", ch.Tbname)
	resQy, resErr := db.Query(qyStr, userData.UserId, false)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")
}

func (ch *CartHandler) AddToCart(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	model := new(model.TicketCartModel)
	uuid := uuid.New()

	// check for permission
	if userData.UserRole != roleId.t {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// fetch from form-data
	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", ch.Tbname_ticket)
	qyRes, qyErr := db.Query(qyStr, model.Ticket_id)
	if qyErr != nil {
		return resp.NotFound(c, "Ticket not found")
	}

	itemPrice := qyRes[0].(map[string]interface{})["price"].(float64)
	model.Total_price = itemPrice * float64(model.Qty)

	// db process
	cmdMainStr := fmt.Sprintf(`
	INSERT INTO %s(
		id, ticket_id, qty, total_price, wait, created_at, created_by, updated_at, updated_by)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`, ch.Tbname)
	resMainErr := db.Command(
		cmdMainStr, uuid, model.Ticket_id, model.Qty, model.Total_price, model.Wait,
		time.Now(), userData.UserId, time.Now(), userData.UserId,
	)
	if resMainErr != nil {
		return resp.ServerError(c, "Error Adding Data: "+resMainErr.Error())
	}

	return resp.Created(c, "Success Adding Data")
}

func (ch *CartHandler) RemoveFromCart(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	// permission check
	if userData.UserRole != roleId.t {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// check file availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", ch.Tbname, c.Params("id"))
	checkData, checkErr := db.Query(qyStr)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// db process
	cmdStr := fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", ch.Tbname, c.Params("id"))
	resErr := db.Command(cmdStr)
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "Success Delete Data")
}
