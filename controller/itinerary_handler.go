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

	if userData.UserRole != roleId.t {
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

func (ih *ItineraryHandler) GetItineraryPublic(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != roleId.t {
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

	// check for search
	qyStr := fmt.Sprintf(`SELECT * FROM %s WHERE POSITION(upper($1) IN upper(title))>0 LIMIT $2 OFFSET $3`, ih.Tbname)
	resQy, resErr := db.Query(qyStr, c.Query("title"), row, (page-1)*row)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

func (ih *ItineraryHandler) GetItineraryById(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != roleId.t {
		return resp.Forbidden(c, "Access Forbidden")
	}

	qyStr := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, ih.Tbname)
	idLoc := c.Params("id")
	resQy, resErr := db.Query(qyStr, idLoc)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	if resQy[0] == nil {
		return resp.NotFound(c, "Itinerary not found")
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

func (ih *ItineraryHandler) CreateItinerary(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	model := new(model.ItineraryModel)
	uuidMain := uuid.New()

	// permission check
	if userData.UserRole != roleId.t {
		return resp.Forbidden(c, "Access Forbidden")
	}

	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	// db process for itinerary main
	cmdMainStr := fmt.Sprintf(`
	INSERT INTO %s(id, title, detail, created_at, created_by, updated_at, updated_by)
	VALUES($1, $2, $3, $4, $5, $6, $7)`, ih.Tbname)
	resMainErr := db.Command(
		cmdMainStr, uuidMain, model.Title, model.Detail, time.Now(), userData.UserId, time.Now(), userData.UserId,
	)
	if resMainErr != nil {
		return resp.ServerError(c, "Error Adding Data: "+resMainErr.Error())
	}

	// db process for itinerary item
	if len(model.Items) > 0 {

		for i := 0; i < len(model.Items); i++ {

			uuidItem := uuid.New()
			itemData := model.Items[i].(map[string]interface{})

			cmdItemStr := fmt.Sprintf(`
			INSERT INTO %s(id, itinerary_id, place_id, detail, went_time, created_at, created_by, updated_at, updated_by)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`, ih.Tbname_item)
			resItemErr := db.Command(
				cmdItemStr, uuidItem, uuidMain,
				itemData["place_id"], itemData["detail"], itemData["went_time"],
				time.Now(), userData.UserId, time.Now(), userData.UserId,
			)
			if resItemErr != nil {
				return resp.ServerError(c, "Error Adding Data: "+resItemErr.Error())
			}

		}

	}

	return resp.Created(c, "Success Adding Data")
}

func (ih *ItineraryHandler) UpdateItinerary(c *fiber.Ctx) error {
	return nil
}

func (ih *ItineraryHandler) DeleteItinerary(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	// permission check
	if userData.UserRole != roleId.t {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// check for file availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", ih.Tbname, c.Params("id"))
	checkData, checkErr := db.Query(qyStr)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// delete data process for itinerary item
	cmdItemStr := fmt.Sprintf("DELETE FROM %s WHERE itinerary_id = '%s'", ih.Tbname_item, c.Params("id"))
	resItemErr := db.Command(cmdItemStr)
	if resItemErr != nil {
		return resp.ServerError(c, resItemErr.Error())
	}

	// delete data process for itinerary main
	cmdMainStr := fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", ih.Tbname, c.Params("id"))
	resMainErr := db.Command(cmdMainStr)
	if resMainErr != nil {
		return resp.ServerError(c, resMainErr.Error())
	}

	return resp.Success(c, nil, "Success Delete Data")

}
