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

type NewsHandler struct {
	Tbname           string
	Tbname_placeinfo string
}

///////////////////
// fetch all news
///////////////////
func (n *NewsHandler) GetNews(c *fiber.Ctx) error {

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

	placeIDCheckStr := ""

	if c.Query("place_id") != "" {
		placeIDCheckStr = fmt.Sprintf(" place_id = '%s' AND ", c.Query("place_id"))
	}

	qyStr := fmt.Sprintf(`
	SELECT id, title, image FROM %s
	WHERE POSITION(upper($1) IN upper(title))>0 AND %s status = $2
	ORDER BY updated_at DESC LIMIT $3 OFFSET $4`, n.Tbname, placeIDCheckStr)
	resQy, resErr := db.Query(qyStr, c.Query("search"), true, row, (page-1)*row)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	if resQy[0] == nil {
		return resp.NotFound(c, "News not found")
	}

	return resp.Success(c, resQy, "Success Fetching Data")
}

/////////////////////
// fetch news admin
/////////////////////
func (n *NewsHandler) GetNewsAdmin(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != roleId.pm {
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

	qyStr := fmt.Sprintf(`
	SELECT id, title, image FROM %s
	WHERE POSITION(upper($1) IN upper(title))>0 AND created_by = $2
	ORDER BY updated_at DESC LIMIT $3 OFFSET $4
	`, n.Tbname)
	resQy, resErr := db.Query(qyStr, c.Query("search"), userData.UserId, row, (page-1)*row)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

/////////////////////
// fetch news by id
/////////////////////
func (n *NewsHandler) GetNewsById(c *fiber.Ctx) error {

	// ReqHeader := c.GetReqHeaders()
	// AuthToken := strings.Split(ReqHeader["Authorization"], " ")[1]

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", n.Tbname)
	resQy, resErr := db.Query(qyStr, c.Params("id"))
	if resErr != nil {
		return resp.ServerError(c, "Server Error")
	}
	if resQy[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

/////////////////
// add new news
/////////////////
func (n *NewsHandler) AddNews(c *fiber.Ctx) error {

	// return c.Send(c.Body())

	// check for permission
	userData := c.Locals("user").(*helper.ClaimsData)
	model := new(model.NewsModel)
	uuid := uuid.New()

	if userData.UserRole != roleId.pm {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// // fetch from form-data
	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	// check for place id
	qyStr := fmt.Sprintf("SELECT id FROM %s WHERE created_by = $1", n.Tbname_placeinfo)
	checkData, checkErr := db.Query(qyStr, userData.UserId)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	model.Place_id = checkData[0].(map[string]interface{})["id"].(string)

	// file process
	fileForm, _ := c.FormFile("image")
	fileName := fmt.Sprintf("%s-%s", uuid, fileForm.Filename)
	for {
		pathDir := "./public/news"
		saveErr := c.SaveFile(fileForm, fmt.Sprintf("%s/%s", pathDir, fileName))
		if saveErr != nil {
			os.MkdirAll(pathDir, 0777)
			continue
		}
		break
	}
	model.Image = fmt.Sprintf("/api/public/news/%s", fileName)

	// db process
	cmdMainStr := fmt.Sprintf(`
	INSERT INTO %s(
		id, place_id, title, article, image, status, draft_status,
		created_at, created_by, updated_at, updated_by)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, n.Tbname)
	resMainErr := db.Command(
		cmdMainStr, uuid, model.Place_id, model.Title, model.Article, model.Image,
		model.Status, model.Draft_status, time.Now(), userData.UserId, time.Now(), userData.UserId,
	)
	if resMainErr != nil {
		return resp.ServerError(c, "Error Adding Data: "+resMainErr.Error())
	}

	return resp.Created(c, "Success Adding Data")
}

////////////////////
// edit news by id
////////////////////
func (n *NewsHandler) EditNews(c *fiber.Ctx) error {

	// check for permission
	model := new(model.NewsModel)
	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != roleId.pm {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// fetch from form-data
	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	// check for data availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", n.Tbname)
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
		c.SaveFile(fileForm, fmt.Sprintf("public/news/%s", fileName))
		model.Image = fmt.Sprintf("/api/public/news/%s", fileName)
	}

	// fill empty data process
	if model.Title == "" {
		model.Title = checkData[0].(map[string]interface{})["title"].(string)
	}
	if model.Article == "" {
		model.Article = checkData[0].(map[string]interface{})["article"].(string)
	}
	if model.Image == nil {
		model.Image = checkData[0].(map[string]interface{})["image"]
	}
	if !model.Status {
		model.Status = checkData[0].(map[string]interface{})["status"].(bool)
	}
	if !model.Draft_status {
		model.Draft_status = checkData[0].(map[string]interface{})["draft_status"].(bool)
	}

	// delete data process
	cmdStr := fmt.Sprintf(
		"UPDATE %s SET title=$1, article=$2, image=$3, status=$4, draft_status=$5, updated_by=$6, updated_at=$7 WHERE id = $8", n.Tbname)

	cmdErr := db.Command(cmdStr, model.Title, model.Article, model.Image, model.Status,
		model.Draft_status, userData.UserId, time.Now(), c.Params("id"))
	if cmdErr != nil {
		resp.ServerError(c, "Error Updating Data: "+cmdErr.Error())
	}

	return resp.Success(c, nil, "Success Updating Data")
}

//////////////////////
// delete news by id
//////////////////////
func (n *NewsHandler) DeleteNews(c *fiber.Ctx) error {

	// check for permission
	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != roleId.pm {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// check for file availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", n.Tbname, c.Params("id"))
	checkData, checkErr := db.Query(qyStr)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// delete file process
	fileNameRaw := checkData[0].(map[string]interface{})["image"]
	fileName := strings.Split(fileNameRaw.(string), "/")

	os.Remove(fmt.Sprintf("./public/news/%s", fileName[4]))

	// delete data process
	cmdStr := fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", n.Tbname, c.Params("id"))
	resErr := db.Command(cmdStr)
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "Success Delete Data")

}
