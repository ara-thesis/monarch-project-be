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
	Tbname           string
	Tbname_item      string
	Tbname_placeinfo string
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
	idItinerary := c.Params("id")
	resQy, resErr := db.Query(qyStr, idItinerary)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	if resQy[0] == nil {
		return resp.NotFound(c, "Itinerary not found")
	}

	qyItemStr := fmt.Sprintf(`
		SELECT i.*, p.place_name, st_x(p.place_loc) AS long, st_y(p.place_loc) AS lat
		FROM %s i JOIN %s p
		ON i.place_id = p.id
		WHERE itinerary_id = $1
		ORDER BY i.in_time ASC`,
		ih.Tbname_item,
		ih.Tbname_placeinfo,
	)
	resItemQy, resItemErr := db.Query(qyItemStr, idItinerary)

	if resItemErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	if resItemQy[0] != nil {
		resQy[0].(map[string]interface{})["item"] = resItemQy
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

func (ih *ItineraryHandler) CreateItinerary(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	modelMain := new(model.ItineraryModel)
	uuidMain := uuid.New()

	// permission check
	if userData.UserRole != roleId.t {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// get body json
	if reqErr := c.BodyParser(modelMain); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	// db process for itinerary main
	cmdMainStr := fmt.Sprintf(`
	INSERT INTO %s(id, title, detail, created_at, created_by, updated_at, updated_by)
	VALUES($1, $2, $3, $4, $5, $6, $7)`, ih.Tbname)
	resMainErr := db.Command(
		cmdMainStr, uuidMain, modelMain.Title, modelMain.Detail, time.Now(), userData.UserId, time.Now(), userData.UserId,
	)
	if resMainErr != nil {
		return resp.ServerError(c, "Error Adding Data: "+resMainErr.Error())
	}

	// db process for itinerary item
	if len(modelMain.Items) > 0 && modelMain.Items[0] != nil {

		for i := 0; i < len(modelMain.Items); i++ {

			uuidItem := uuid.New()
			modelItem := &model.ItineraryItemModel{
				ItineraryId: uuidMain,
				PlaceId:     modelMain.Items[i].(map[string]interface{})["place_id"].(string),
				Detail:      modelMain.Items[i].(map[string]interface{})["detail"].(string),
				In_time:     modelMain.Items[i].(map[string]interface{})["in_time"].(string),
				Out_time:    modelMain.Items[i].(map[string]interface{})["out_time"].(string),
			}

			cmdItemStr := fmt.Sprintf(`
			INSERT INTO %s(id, itinerary_id, place_id, detail, in_time, out_time, created_at, created_by, updated_at, updated_by)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, ih.Tbname_item)
			resItemErr := db.Command(
				cmdItemStr, uuidItem, uuidMain, modelItem.PlaceId, modelItem.Detail, modelItem.In_time, modelItem.Out_time,
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

	userData := c.Locals("user").(*helper.ClaimsData)
	model := new(model.ItineraryModel)
	itineraryId := c.Params("id")

	// permission check
	if userData.UserRole != roleId.t {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// check for file availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", ih.Tbname)
	checkData, checkErr := db.Query(qyStr, itineraryId)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// get body json
	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	if model.Title == "" {
		model.Title = checkData[0].(map[string]string)["title"]
	}
	if model.Detail == "" {
		model.Detail = checkData[0].(map[string]string)["detail"]
	}

	// update main itinerary data process
	cmdMainStr := fmt.Sprintf(`UPDATE %s SET title = $1, detail = $2, updated_at = $3, updated_by = $4 WHERE id = $5`, ih.Tbname)
	resMainErr := db.Command(cmdMainStr, model.Title, model.Detail, time.Now(), userData.UserId, itineraryId)

	if resMainErr != nil {
		return resp.ServerError(c, "Error Updating Data: "+resMainErr.Error())
	}

	// update itinerary item data process
	if model.Items != nil {

		// delete all the data first
		cmdItemDelStr := fmt.Sprintf(`DELETE FROM %s WHERE itinerary_id = $1`, ih.Tbname_item)
		errItemDel := db.Command(cmdItemDelStr, itineraryId)

		if errItemDel != nil {
			return resp.ServerError(c, "Error Updating Data: "+resMainErr.Error())
		}

		// create new one
		for i := 0; i < len(model.Items); i++ {

			if model.Items[0] == nil {
				break
			}

			uuidItem := uuid.New()
			itemData := model.Items[i].(map[string]interface{})

			cmdItemAddStr := fmt.Sprintf(`
			INSERT INTO %s(id, itinerary_id, place_id, detail, in_time, out_time, created_at, created_by, updated_at, updated_by)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, ih.Tbname_item)
			resItemAddErr := db.Command(
				cmdItemAddStr, uuidItem, itineraryId, itemData["place_id"], itemData["detail"], itemData["in_time"], itemData["out_time"],
				time.Now(), userData.UserId, time.Now(), userData.UserId,
			)
			if resItemAddErr != nil {
				return resp.ServerError(c, "Error Adding Data: "+resItemAddErr.Error())
			}

		}

	}

	return resp.Success(c, nil, "Success Updating Data")
}

func (ih *ItineraryHandler) DeleteItinerary(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	itineraryId := c.Params("id")

	// permission check
	if userData.UserRole != roleId.t {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// check for file availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", ih.Tbname)
	checkData, checkErr := db.Query(qyStr, itineraryId)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// delete data process for itinerary item
	cmdItemStr := fmt.Sprintf("DELETE FROM %s WHERE itinerary_id = $1", ih.Tbname_item)
	resItemErr := db.Command(cmdItemStr, itineraryId)
	if resItemErr != nil {
		return resp.ServerError(c, resItemErr.Error())
	}

	// delete data process for itinerary main
	cmdMainStr := fmt.Sprintf("DELETE FROM %s WHERE id = $1", ih.Tbname)
	resMainErr := db.Command(cmdMainStr, itineraryId)
	if resMainErr != nil {
		return resp.ServerError(c, resMainErr.Error())
	}

	return resp.Success(c, nil, "Success Delete Data")

}
