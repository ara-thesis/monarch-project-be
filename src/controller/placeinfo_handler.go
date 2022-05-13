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

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tbname["placeinfo"])
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

func (pinf *PlaceInfoHandler) AddAndEditPlaceInfoAdmin(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	model := new(model.PlaceInfoModel)
	uuid := uuid.New()

	if userData.UserRole != "PLACE MANAGER" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	qyStr := fmt.Sprintf(`SELECT * FROM %s WHERE created_by = $1`, tbname["placeinfo"])
	qyRes, qyErr := db.Query(qyStr, userData.UserId)
	if qyErr != nil {
		return resp.ServerError(c, "Server Error")
	}

	if qyRes[0] == nil {
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
		return resp.Created(c, "Success Creating Data")
	}

	qyFinData := qyRes[0]

	if model.Place_name == nil {
		model.Place_name = qyFinData.(map[string]interface{})["place_name"]
	}
	if model.Place_info == nil {
		model.Place_info = qyFinData.(map[string]interface{})["place_info"]
	}
	if model.Place_city == nil {
		model.Place_city = qyFinData.(map[string]interface{})["place_city"]
	}
	if model.Place_stateprov == nil {
		model.Place_stateprov = qyFinData.(map[string]interface{})["place_stateprov"]
	}
	if model.Place_country == nil {
		model.Place_country = qyFinData.(map[string]interface{})["place_country"]
	}
	if model.Place_postal == nil {
		model.Place_postal = qyFinData.(map[string]interface{})["place_postal"]
	}
	if model.Place_loc == nil {
		model.Place_loc = qyFinData.(map[string]interface{})["place_loc"]
	}
	if model.Place_opentime == nil {
		model.Place_opentime = qyFinData.(map[string]interface{})["place_opentime"]
	}
	if model.Place_closetime == nil {
		model.Place_closetime = qyFinData.(map[string]interface{})["place_closetime"]
	}

	cmdStr := fmt.Sprintf(`
	UPDATE %s SET
	place_name=$1, place_info=$2, place_city=$3, place_stateprov=$4, place_country=$5,
	place_postal=$6, place_loc=$7, place_opentime=$8, place_closetime=$9,
	updated_at=$10, updated_by=$11`, tbname["placeinfo"])
	resErr := db.Command(cmdStr, model.Place_name, model.Place_info, model.Place_city, model.Place_stateprov, model.Place_country,
		model.Place_postal, model.Place_loc, model.Place_opentime, model.Place_closetime, time.Now(), userData.UserId)
	if resErr != nil {
		return resp.ServerError(c, "Server Error: "+resErr.Error())
	}

	return resp.Success(c, nil, "Success Updating Data")
}

func (pinf *PlaceInfoHandler) DeletePlaceInfoAdmin(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != "PLACE MANAGER" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE created_by = $1", tbname["placeinfo"])
	checkData, checkErr := db.Query(qyStr, c.Params("userId"))
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
