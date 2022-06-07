package controller

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ara-thesis/monarch-project-be/helper"
	"github.com/ara-thesis/monarch-project-be/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type NewsHandler struct{}

// fetch all news
func (n *NewsHandler) GetNews(c *fiber.Ctx) error {

	// ReqHeader := c.GetReqHeaders()
	// AuthToken := strings.Split(ReqHeader["Authorization"], " ")[1]

	qyStr := fmt.Sprintf("SELECT * FROM %s", tbname["news"])
	resQy, resErr := db.Query(qyStr)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")
}

func (n *NewsHandler) GetNewsAdmin(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != "PLACE MANAGER" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE created_by = $1", tbname["news"])
	resQy, resErr := db.Query(qyStr, userData.UserId)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

// fetch news by id
func (n *NewsHandler) GetNewsById(c *fiber.Ctx) error {

	// ReqHeader := c.GetReqHeaders()
	// AuthToken := strings.Split(ReqHeader["Authorization"], " ")[1]

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

// add new news
func (n *NewsHandler) AddNews(c *fiber.Ctx) error {

	// check for permission
	userData := c.Locals("user").(*helper.ClaimsData)
	model := new(model.NewsModel)
	uuid := uuid.New()

	if userData.UserRole != "PLACE MANAGER" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// fetch from form-data
	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

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
		id, title, article, image, status, draft_status,
		created_at, created_by, updated_at, updated_by)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, tbname["news"])
	resMainErr := db.Command(
		cmdMainStr, uuid, model.Title, model.Article, model.Image, model.Status,
		model.Draft_status, time.Now(), userData.UserId, time.Now(), userData.UserId,
	)
	if resMainErr != nil {
		return resp.ServerError(c, "Error Adding Data: "+resMainErr.Error())
	}

	return resp.Created(c, "Success Adding Data")
}

// edit news by id
func (n *NewsHandler) EditNews(c *fiber.Ctx) error {

	// check for permission
	model := new(model.NewsModel)
	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != "PLACE MANAGER" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// fetch from form-data
	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	// check for data availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tbname["news"])
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
	if model.Title == nil {
		model.Title = checkData[0].(map[string]interface{})["title"]
	}
	if model.Article == nil {
		model.Article = checkData[0].(map[string]interface{})["article"]
	}
	if model.Image == nil {
		model.Image = checkData[0].(map[string]interface{})["image"]
	}
	if model.Status == nil {
		model.Status = checkData[0].(map[string]interface{})["status"]
	}
	if model.Draft_status == nil {
		model.Draft_status = checkData[0].(map[string]interface{})["draft_status"]
	}

	// delete data process
	cmdStr := fmt.Sprintf(
		"UPDATE %s SET title=$1, article=$2, image=$3, status=$4, draft_status=$5, updated_by=$6, updated_at=$7 WHERE id = $8",
		tbname["news"])

	cmdErr := db.Command(cmdStr, model.Title, model.Article, model.Image, model.Status,
		model.Draft_status, userData.UserId, time.Now(), c.Params("id"))
	if cmdErr != nil {
		resp.ServerError(c, "Error Updating Data: "+cmdErr.Error())
	}

	return resp.Success(c, nil, "Success Updating Data")
}

// delete news by id
func (n *NewsHandler) DeleteNews(c *fiber.Ctx) error {

	// check for permission
	userData := c.Locals("user").(*helper.ClaimsData)

	if userData.UserRole != "PLACE MANAGER" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// check for file availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", tbname["news"], c.Params("id"))
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
	cmdStr := fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", tbname["news"], c.Params("id"))
	resErr := db.Command(cmdStr)
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "Success Delete Data")

}