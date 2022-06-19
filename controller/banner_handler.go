package controller

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ara-thesis/monarch-project-be/helper"
	"github.com/ara-thesis/monarch-project-be/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type BannerHandler struct{}

////////////////////
// fetch all banner
////////////////////
func (n *BannerHandler) GetBanners(c *fiber.Ctx) error {

	// ReqHeader := c.GetReqHeaders()
	// AuthToken := strings.Split(ReqHeader["Authorization"], " ")[1]

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

	qyStr := fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", tbname["banner"])
	resQy, resErr := db.Query(qyStr, row, (page-1)*row)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")
}

///////////////////////
// fetch banner by id
///////////////////////
func (n *BannerHandler) GetBannerById(c *fiber.Ctx) error {

	// ReqHeader := c.GetReqHeaders()
	// AuthToken := strings.Split(ReqHeader["Authorization"], " ")[1]

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tbname["banner"])
	resQy, resErr := db.Query(qyStr, c.Params("id"))
	if resErr != nil {
		return resp.ServerError(c, "Server Error")
	}
	if resQy[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

///////////////////
// add new banner
//////////////////
func (n *BannerHandler) AddBanner(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	// model := new(model.BannerModel)
	model := &model.BannerModel{}
	uuid := uuid.New()

	// permission check
	if userData.UserRole != roleId["ADM"] {
		return resp.Forbidden(c, "Access Forbidden")
	}

	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	test := c.FormValue("title")

	// file process
	fileForm, _ := c.FormFile("image")
	fileName := fmt.Sprintf("%s-%s", uuid, fileForm.Filename)
	for {
		pathDir := "public/banner" + test
		saveErr := c.SaveFile(fileForm, fmt.Sprintf("%s/%s", pathDir, fileName))
		if saveErr != nil {
			os.MkdirAll(pathDir, 0777)
			continue
		}
		break
	}
	model.Image = fmt.Sprintf("/api/public/banner/%s", fileName)

	// db process
	cmdMainStr := fmt.Sprintf(`
	INSERT INTO %s(
		id, title, detail, image, status, created_at, created_by, updated_at, updated_by)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`, tbname["banner"])
	resMainErr := db.Command(
		cmdMainStr, uuid, model.Title, model.Detail, model.Image, model.Status,
		time.Now(), userData.UserId, time.Now(), userData.UserId,
	)
	if resMainErr != nil {
		return resp.ServerError(c, "Error Adding Data: "+resMainErr.Error())
	}

	return resp.Created(c, "Success Adding Data")
}

/////////////////////
// edit banner by id
/////////////////////
func (n *BannerHandler) EditBanner(c *fiber.Ctx) error {

	model := new(model.BannerModel)
	userData := c.Locals("user").(*helper.ClaimsData)

	// permission check
	if userData.UserRole != "ADMIN" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	// check file availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tbname["banner"])
	checkData, checkErr := db.Query(qyStr, c.Params("id"))
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// file process
	fileForm, fileErr := c.FormFile("image")
	if fileErr == nil {
		fileName := fmt.Sprintf("%s-%s", checkData[0].(map[string]interface{})["id"], fileForm.Filename)
		c.SaveFile(fileForm, fmt.Sprintf("public/banner/%s", fileName))
		model.Image = fmt.Sprintf("/api/public/banner/%s", fileName)
	}

	// replace data process
	if model.Title == "" {
		model.Title = checkData[0].(map[string]interface{})["title"].(string)
	}
	if model.Detail == "" {
		model.Detail = checkData[0].(map[string]interface{})["article"].(string)
	}
	if model.Image == nil {
		model.Image = checkData[0].(map[string]interface{})["image"]
	}
	if !model.Status {
		model.Status = checkData[0].(map[string]interface{})["status"].(bool)
	}

	// update data process
	cmdStr := fmt.Sprintf("UPDATE %s SET title=$1, detail=$2, image=$3, status=$4, updated_by=$5, updated_at=$6 WHERE id = $7", tbname["banner"])

	cmdErr := db.Command(cmdStr, model.Title, model.Detail, model.Image, model.Status, userData.UserId, time.Now(), c.Params("id"))
	if cmdErr != nil {
		resp.ServerError(c, "Error Updating Data: "+cmdErr.Error())
	}

	return resp.Success(c, nil, "Success Updating Data")
}

///////////////////////
// delete banner by id
///////////////////////
func (n *BannerHandler) DeleteBanner(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	// permission check
	if userData.UserRole != "ADMIN" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// check file availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", tbname["banner"], c.Params("id"))
	checkData, checkErr := db.Query(qyStr)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// file process
	fileNameRaw := checkData[0].(map[string]interface{})["image"]
	fileName := strings.Split(fileNameRaw.(string), "/")

	os.Remove(fmt.Sprintf("./public/banner/%s", fileName[4]))

	// db process
	cmdStr := fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", tbname["banner"], c.Params("id"))
	resErr := db.Command(cmdStr)
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "Success Delete Data")

}
