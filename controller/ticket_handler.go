package controller

import (
	"fmt"
	"time"

	"github.com/ara-thesis/monarch-project-be/helper"
	"github.com/ara-thesis/monarch-project-be/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TicketHandler struct {
	Tbname        string
	Tbname_place  string
	Tbname_bought string
}

/////////////////////
// fetch all ticket
/////////////////////
func (th *TicketHandler) GetTicketTourist(c *fiber.Ctx) error {

	placeId := c.Query("placeid")
	if placeId == "" {
		return resp.BadRequest(c, "Need more query (placeid is needed)")
	}

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE place_id = $1", th.Tbname)
	resQy, resErr := db.Query(qyStr, placeId)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

/////////////////////
// fetch all ticket
/////////////////////
func (th *TicketHandler) GetTicketById(c *fiber.Ctx) error {

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", th.Tbname)
	resQy, resErr := db.Query(qyStr, c.Params("id"))

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

////////////////////////////
// fetch all ticket bought
////////////////////////////
func (th *TicketHandler) GetTicketBoughtTourist(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	// check for permission
	if userData.UserRole != roleId.t {
		return resp.Forbidden(c, "Access Forbidden")
	}

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE owner = $1 AND redeemed = $2", th.Tbname_bought)
	resQy, resErr := db.Query(qyStr, userData.UserId, false)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

//////////////////
// redeem ticket
//////////////////
func (th *TicketHandler) RedeemTicket(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	model := &model.TicketBoughtModel{}

	// check for permission
	if userData.UserRole != roleId.pm {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// fetch from form-data
	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	cmdStr := fmt.Sprintf("UPDATE %s SET redeemed = $1, ticket_id=$2, updated_by = $3, updated_at = $4 WHERE code = $5", th.Tbname_bought)
	resErr := db.Command(cmdStr, true, nil, userData.UserId, time.Now(), model.TicketBought_code)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "Success Redeem Ticket")

}

///////////////////////////////
// fetch all ticket for admin
///////////////////////////////
func (th *TicketHandler) GetTicketAdmin(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	// check for permission
	if userData.UserRole != roleId.pm {
		return resp.Forbidden(c, "Access Forbidden")
	}

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE created_by = $1", th.Tbname)
	resQy, resErr := db.Query(qyStr, userData.UserId)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

///////////////////
// add new ticket
///////////////////
func (th *TicketHandler) AddTicket(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	model := new(model.TicketModel)
	uuid := uuid.New()

	// check for permission
	if userData.UserRole != roleId.pm {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// fetch from form-data
	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	// check for place id
	qyStr := fmt.Sprintf("SELECT id FROM %s WHERE created_by = $1", th.Tbname_place)
	checkData, checkErr := db.Query(qyStr, userData.UserId)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	model.Ticket_placeid = checkData[0].(map[string]interface{})["id"].(string)

	// db process
	cmdMainStr := fmt.Sprintf(`
	INSERT INTO %s(
		id, name, details, price, place_id, created_at, created_by, updated_at, updated_by)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`, th.Tbname)
	resMainErr := db.Command(
		cmdMainStr, uuid, model.Ticket_name, model.Ticket_details, model.Ticket_price, model.Ticket_placeid,
		time.Now(), userData.UserId, time.Now(), userData.UserId,
	)
	if resMainErr != nil {
		return resp.ServerError(c, "Error Adding Data: "+resMainErr.Error())
	}

	return resp.Created(c, "Success Adding Data")

}

////////////////////
// edit ticket by id
////////////////////
func (th *TicketHandler) EditTicket(c *fiber.Ctx) error {

	// check for permission
	model := new(model.TicketModel)
	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != roleId.pm {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// fetch from form-data
	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	// check for data availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", th.Tbname)
	checkData, checkErr := db.Query(qyStr, c.Params("id"))
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// fill empty data process
	if model.Ticket_name == "" {
		model.Ticket_name = checkData[0].(map[string]interface{})["name"].(string)
	}
	if model.Ticket_details == "" {
		model.Ticket_details = checkData[0].(map[string]interface{})["details"].(string)
	}
	if model.Ticket_placeid == "" {
		model.Ticket_placeid = checkData[0].(map[string]interface{})["place_id"].(string)
	}
	if model.Ticket_price == 0 {
		model.Ticket_price = checkData[0].(map[string]interface{})["price"].(int)
	}

	// update data process
	cmdStr := fmt.Sprintf(
		"UPDATE %s SET name=$1, details=$2, place_id=$3, price=$4, updated_by=$5, updated_at=$6 WHERE id = $7",
		th.Tbname)

	cmdErr := db.Command(cmdStr, model.Ticket_name, model.Ticket_details, model.Ticket_placeid,
		model.Ticket_price, userData.UserId, time.Now(), c.Params("id"))
	if cmdErr != nil {
		resp.ServerError(c, "Error Updating Data: "+cmdErr.Error())
	}

	return resp.Success(c, nil, "Success Updating Data")
}

//////////////////////
// delete ticket by id
//////////////////////
func (th *TicketHandler) DeleteTicket(c *fiber.Ctx) error {

	// check for permission
	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != roleId.pm {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// check for file availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", th.Tbname, c.Params("id"))
	checkData, checkErr := db.Query(qyStr)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// delete data process
	cmdStr := fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", th.Tbname, c.Params("id"))
	resErr := db.Command(cmdStr)
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "Success Delete Data")

}
