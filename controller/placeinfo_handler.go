package controller

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ara-thesis/monarch-project-be/helper"
	"github.com/ara-thesis/monarch-project-be/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PlaceInfoHandler struct {
	Tbname     string
	Tbname_img string
}

/////////////////////////
// fetch all place info
/////////////////////////
func (pinf *PlaceInfoHandler) GetPlaceInfo(c *fiber.Ctx) error {

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

	searchKey := c.Query("search")

	qyStr := fmt.Sprintf(`
	SELECT *, st_x(place_loc) AS place_loc_long, st_y(place_loc) AS place_loc_lat
	FROM %s
	WHERE POSITION(upper($1) IN upper(place_name))>0
	LIMIT $2 OFFSET $3`, pinf.Tbname)
	resQy, resErr := db.Query(qyStr, searchKey, row, (page-1)*row)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	if resQy[0] != nil {
		for i := 0; i < len(resQy); i++ {
			qyImgStr := fmt.Sprintf(`
			SELECT image
			FROM %s
			WHERE place_id = $1
			`, pinf.Tbname_img)
			resImgQy, resImgErr := db.Query(qyImgStr, resQy[i].(map[string]interface{})["id"])

			if resImgErr != nil {
				return resp.ServerError(c, resErr.Error())
			}

			resQy[i].(map[string]interface{})["images"] = resImgQy

		}
	}

	if resQy[0] == nil && len(searchKey) != 0 {
		return resp.NotFound(c, "Place not found!")
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

///////////////////////////
// fetch place info by id
///////////////////////////
func (pinf *PlaceInfoHandler) GetPlaceInfoById(c *fiber.Ctx) error {

	qyStr := fmt.Sprintf(`SELECT *, st_x(place_loc) AS place_loc_long, st_y(place_loc) AS place_loc_lat
		FROM %s WHERE id = $1`, pinf.Tbname)
	resQy, resErr := db.Query(qyStr, c.Params("id"))
	if resErr != nil {
		return resp.ServerError(c, "Server Error")
	}
	if resQy[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	qyImgStr := fmt.Sprintf(`
	SELECT image FROM %s WHERE place_id = $1;
	`, pinf.Tbname_img)
	resQyImg, resErrImg := db.Query(qyImgStr, c.Params("id"))
	if resErrImg != nil {
		return resp.ServerError(c, "Server Error")
	}

	resQy[0].(map[string]interface{})["images"] = resQyImg

	return resp.Success(c, resQy, "Success Fetching Data")

}

///////////////////////////
// fetch place info admin
///////////////////////////
func (pinf *PlaceInfoHandler) GetPlaceInfoAdmin(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != roleId.pm {
		return resp.Forbidden(c, "Access Forbidden")
	}

	qyStr := fmt.Sprintf(`SELECT *, st_x(place_loc) AS place_loc_long, st_y(place_loc) AS place_loc_lat
	FROM %s WHERE created_by = $1`, pinf.Tbname)
	resQy, resErr := db.Query(qyStr, userData.UserId)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

//////////////////////////////
// add and update place info
//////////////////////////////
func (pinf *PlaceInfoHandler) UpdatePlaceInfoAdmin(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	model := new(model.PlaceInfoModel)
	uuid_main := uuid.New()

	if userData.UserRole != roleId.pm {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// bodyparser process
	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	// check for placeinfo
	qyStr := fmt.Sprintf(`SELECT * FROM %s WHERE created_by = $1`, pinf.Tbname)
	qyRes, qyErr := db.Query(qyStr, userData.UserId)
	if qyErr != nil {
		return resp.ServerError(c, "Server Error")
	}

	// geom data process
	place_loc := ""
	if model.Place_loc_long != 0 && model.Place_loc_lat != 0 {
		place_loc = fmt.Sprintf("POINT(%f %f)", model.Place_loc_long, model.Place_loc_lat)
	} else {
		place_loc = "POINT(0 0)"
	}

	// file process
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["images"]
		for _, file := range files {
			fileName := fmt.Sprintf("%s-%s", uuid_main, file.Filename)
			for {
				pathDir := "./public/placeinfo"
				saveErr := c.SaveFile(file, fmt.Sprintf("%s/%s", pathDir, fileName))
				if saveErr != nil {
					os.MkdirAll(pathDir, 0777)
					continue
				}
				break
			}
			model.Place_images = append(model.Place_images, fmt.Sprintf("/api/public/placeinfo/%s", fileName))
		}
	}

	// insert mode when placeinfo not found
	if qyRes[0] == nil {

		cmdStr := fmt.Sprintf(`
		INSERT INTO %s(
			id, place_name, place_info, place_street, place_city, place_stateprov,
			place_country, place_postal, place_loc, place_opentime, place_closetime,
			created_at, created_by, updated_at, updated_by)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`, pinf.Tbname)
		resErr := db.Command(cmdStr, uuid_main, model.Place_name, model.Place_info, model.Place_street, model.Place_city,
			model.Place_stateprov, model.Place_country, model.Place_postal, place_loc, model.Place_opentime,
			model.Place_closetime, time.Now(), userData.UserId, time.Now(), userData.UserId)
		if resErr != nil {
			return resp.ServerError(c, "Error Adding Data: "+resErr.Error())
		}

		for _, image_url := range model.Place_images {
			uuid_img := uuid.New()
			cmdImgStr := fmt.Sprintf(`
			INSERT INTO %s(
				id, place_id, image, created_at, created_by, updated_at, updated_by
			)
			VALUES($1, $2, $3, $4, $5, $6, $7)`, pinf.Tbname_img)
			resImgErr := db.Command(cmdImgStr, uuid_img, uuid_main, image_url, time.Now(), userData.UserId, time.Now(), userData.UserId)
			if resImgErr != nil {
				return resp.ServerError(c, "Error Adding Data: "+resErr.Error())
			}
		}

		return resp.Created(c, "Success Creating Data")
	}

	// update mode when placeinfo is found

	qyFinData := qyRes[0]

	if model.Place_name == "" {
		model.Place_name = qyFinData.(map[string]interface{})["place_name"].(string)
	}
	if model.Place_info == "" {
		model.Place_info = qyFinData.(map[string]interface{})["place_info"].(string)
	}
	if model.Place_street == "" {
		model.Place_street = qyFinData.(map[string]interface{})["place_street"].(string)
	}
	if model.Place_city == "" {
		model.Place_city = qyFinData.(map[string]interface{})["place_city"].(string)
	}
	if model.Place_stateprov == "" {
		model.Place_stateprov = qyFinData.(map[string]interface{})["place_stateprov"].(string)
	}
	if model.Place_country == "" {
		model.Place_country = qyFinData.(map[string]interface{})["place_country"].(string)
	}
	if model.Place_postal == "" {
		model.Place_postal = qyFinData.(map[string]interface{})["place_postal"].(string)
	}
	if model.Place_loc_lat == 0 || model.Place_loc_long == 0 {
		place_loc = fmt.Sprintf("%s", qyFinData.(map[string]interface{})["place_loc"])
	}
	if fmt.Sprintf("%v", model.Place_opentime) == "0001-01-01 00:00:00 +0000 UTC" {
		model.Place_opentime = qyFinData.(map[string]interface{})["place_opentime"].(time.Time)
	}
	if fmt.Sprintf("%v", model.Place_closetime) == "0001-01-01 00:00:00 +0000 UTC" {
		model.Place_closetime = qyFinData.(map[string]interface{})["place_closetime"].(time.Time)
	}
	// if fmt.Sprintf("%v", model.Place_opentime) == "" {
	// 	model.Place_opentime = qyFinData.(map[string]interface{})["place_opentime"].(string)
	// }
	// if fmt.Sprintf("%v", model.Place_closetime) == "" {
	// 	model.Place_closetime = qyFinData.(map[string]interface{})["place_closetime"].(string)
	// }

	cmdStr := fmt.Sprintf(`
	UPDATE %s SET
	place_name=$1, place_info=$2, place_street=$3, place_city=$4, place_stateprov=$5, place_country=$6,
	place_postal=$7, place_loc=$8, place_opentime=$9, place_closetime=$10,
	updated_at=$11, updated_by=$12`, pinf.Tbname)
	resErr := db.Command(cmdStr, model.Place_name, model.Place_info, model.Place_street, model.Place_city, model.Place_stateprov,
		model.Place_country, model.Place_postal, place_loc, model.Place_opentime, model.Place_closetime, time.Now(), userData.UserId)
	if resErr != nil {
		return resp.ServerError(c, "Server Error: "+resErr.Error())
	}

	cmdDelImgStr := fmt.Sprintf(`
	DELETE FROM %s WHERE created_by = $1
	`, pinf.Tbname_img)
	resDelImgStr := db.Command(cmdDelImgStr, userData.UserId)
	if resDelImgStr != nil {
		return resp.ServerError(c, "Error Deleting Data: "+resErr.Error())
	}

	for _, image_url := range model.Place_images {
		uuid_img := uuid.New()
		cmdAddImgStr := fmt.Sprintf(`
		INSERT INTO %s(
			id, place_id, image, created_at, created_by, updated_at, updated_by
		)
		VALUES($1, $2, $3, $4, $5, $6, $7)`, pinf.Tbname_img)
		resAddImgErr := db.Command(cmdAddImgStr, uuid_img, qyRes[0].(map[string]interface{})["id"].(string), image_url, time.Now(), userData.UserId, time.Now(), userData.UserId)
		if resAddImgErr != nil {
			return resp.ServerError(c, "Error Adding Data: "+resErr.Error())
		}
	}

	return resp.Success(c, nil, "Success Updating Data")
}

//////////////////////
// delete place info
//////////////////////
func (pinf *PlaceInfoHandler) DeletePlaceInfoAdmin(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != roleId.pm {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// check for data
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE created_by = $1", pinf.Tbname)
	checkData, checkErr := db.Query(qyStr, c.Params("userId"))
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// delte placeinfo image data
	cmdImgStr := fmt.Sprintf("DELETE FROM %s WHERE place_id = $1", pinf.Tbname_img)
	resImgErr := db.Command(cmdImgStr, c.Params("id"))
	if resImgErr != nil {
		return resp.ServerError(c, resImgErr.Error())
	}
	// delete placeinfo data
	cmdStr := fmt.Sprintf("DELETE FROM %s WHERE id = $1", pinf.Tbname_img)
	resErr := db.Command(cmdStr, c.Params("id"))
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "Success Delete Data")

}
