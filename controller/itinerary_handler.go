package controller

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ara-thesis/monarch-project-be/helper"
	"github.com/ara-thesis/monarch-project-be/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ItineraryHandler struct {
	Tbname      string
	Tbname_item string
}

func (ih *ItineraryHandler) GetItinerary(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != "TOURIST" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	row, rowErr := strconv.Atoi(c.Query("row", "10"))
	if rowErr != nil {
		row = 10
	}
	if row > 100 {
		row = 100
	}
	page, pageErr := strconv.Atoi(c.Query("page", "1"))
	if pageErr != nil {
		page = 1
	}

	qyStr := fmt.Sprintf(`SELECT * FROM %s WHERE created_by = $1 LIMIT $2 OFFSET $3`, ih.Tbname)
	resQy, resErr := db.Query(qyStr, userData.UserId, row, (page-1)*row)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

func (ih *ItineraryHandler) GetItineraryById(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != "TOURIST" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	row, rowErr := strconv.Atoi(c.Query("row", "10"))
	if rowErr != nil {
		row = 10
	}
	if row > 100 {
		row = 100
	}
	page, pageErr := strconv.Atoi(c.Query("page", "1"))
	if pageErr != nil {
		page = 1
	}

	qyStr := fmt.Sprintf(`SELECT * FROM %s WHERE created_by = $1 LIMIT $2 OFFSET $3`, ih.Tbname)
	resQy, resErr := db.Query(qyStr, userData.UserId, row, (page-1)*row)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

func (ih *ItineraryHandler) CreateItinerary(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	model := new(model.ItineraryModel)
	uuid := uuid.New()

	// permission check
	if userData.UserRole != roleId["T"] {
		return resp.Forbidden(c, "Access Forbidden")
	}

	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	// db process
	cmdMainStr := fmt.Sprintf(`
	INSERT INTO %s(id, title, detail, created_at, created_by, updated_at, updated_by) VALUES($1, $2, $3, $4, $5, $6, $7)`, ih.Tbname)
	resMainErr := db.Command(
		cmdMainStr, uuid, model.Title, model.Detail, time.Now(), userData.UserId, time.Now(), userData.UserId,
	)
	if resMainErr != nil {
		return resp.ServerError(c, "Error Adding Data: "+resMainErr.Error())
	}

	return resp.Created(c, "Success Adding Data")
}

func (ih *ItineraryHandler) UpdateItinerary(c *fiber.Ctx) error {
	return nil
}

func (ih *ItineraryHandler) DeleteItinerary(c *fiber.Ctx) error {
	return nil
}
