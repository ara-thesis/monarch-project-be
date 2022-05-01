package controller

import (
	"fmt"
	"time"

	"github.com/ara-thesis/monarch-project-be/src/helper"
	"github.com/ara-thesis/monarch-project-be/src/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PlaceInfoHandler struct{}

func (pinf *PlaceInfoHandler) GetPlaceInfo(c *fiber.Ctx) error {

	qyStr := fmt.Sprintf("SELECT * FROM %s", tbname["placeinfo"])
	resQy, resErr := db.Query(qyStr)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

func (pinf *PlaceInfoHandler) GetPlaceInfoById(c *fiber.Ctx) error {

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tbname["news"])
	resQy, resErr := db.Query(qyStr, c.Params("id"))
	if resErr != nil {
		return resp.ServerError(c, "Server Error")
	}
	if resQy[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

func (pinf *PlaceInfoHandler) GetPlaceInfoAdmin(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != "PLACE MANAGER" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE created_by = $1", tbname["placeinfo"])
	resQy, resErr := db.Query(qyStr, userData.UserId)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

func (pinf *PlaceInfoHandler) AddPlaceInfoAdmin(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	model := new(model.PlaceInfoModel)
	uuid := uuid.New()

	if userData.UserRole != "PLACE MANAGER" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	cmdStr := fmt.Sprintf(`
	INSERT INTO %s(
		id, place_name, place_info, place_city, place_stateprov, place_coutnry,
		place_postal, palce_loc, place_opentime, place_closetime,
		created_at, created_by, updated_at, updated_by)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, tbname["placeinfo"])
	resErr := db.Command(cmdStr, uuid, model.Place_name, model.Place_info, model.Place_city, model.Place_stateprov, model.Place_country,
		model.Place_postal, model.Place_loc, model.Place_opentime, model.Place_closetime, time.Now(), userData.UserId, time.Now(), userData.UserId)
	if resErr != nil {
		return resp.ServerError(c, "Error Adding Data: "+resErr.Error())
	}

	return resp.Created(c, "Success Adding Data")
}

func (pinf *PlaceInfoHandler) EditPlaceInfoAdmin(c *fiber.Ctx) error {

	return nil
}

func (pinf *PlaceInfoHandler) DeletePlaceInfoAdmin(c *fiber.Ctx) error {

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tbname["news"])
	checkData, checkErr := db.Query(qyStr, c.Params("id"))
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	cmdStr := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tbname["news"])
	resErr := db.Command(cmdStr, c.Params("id"))
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "Success Delete Data")

}
