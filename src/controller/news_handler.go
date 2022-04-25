package controller

import (
	"fmt"
	"time"

	"github.com/ara-thesis/monarch-project-be/src/helper"
	"github.com/ara-thesis/monarch-project-be/src/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type NewsHandler struct{}

var (
	tbname = "newstb"
	db     = new(helper.PgHelper)
	resp   = new(helper.ResponseHelper)
)

func (n *NewsHandler) GetNews(c *fiber.Ctx) error {

	// ReqHeader := c.GetReqHeaders()
	// AuthToken := strings.Split(ReqHeader["Authorization"], " ")[1]

	qyStr := fmt.Sprintf("SELECT * FROM %s", tbname)
	resArr, resErr := db.Query(qyStr)

	if resErr != nil {
		return resp.Failed(c, resErr.Error(), 500)
	}

	return resp.Success(c, resArr, "Success Fetching Data", 200)
}

func (n *NewsHandler) GetNewsById(c *fiber.Ctx) error {

	// ReqHeader := c.GetReqHeaders()
	// AuthToken := strings.Split(ReqHeader["Authorization"], " ")[1]

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", tbname, c.Params("id"))
	resArr, resErr := db.Query(qyStr)
	if resErr != nil {
		return resp.Failed(c, "Server Error", 500)
	}
	if resArr[0] == nil {
		return resp.Failed(c, "Data Not Found", 404)
	}

	return resp.Success(c, resArr, "Success Fetching Data", 200)

}

func (n *NewsHandler) AddNews(c *fiber.Ctx) error {

	model := new(model.NewsModel)
	uuid := uuid.New()

	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.Failed(c, reqErr.Error(), 500)
	}

	cmdStr := fmt.Sprintf(`INSERT INTO %s(id, title, article, status, draft_status, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7)`, tbname)
	resErr := db.Command(cmdStr, uuid, model.Title, model.Article, model.Status, model.Draft_status, time.Now(), time.Now())
	if resErr != nil {
		return resp.Failed(c, resErr.Error(), 500)
	}

	return resp.Success(c, nil, "Data added", 201)
}

func (n *NewsHandler) EditNews(c *fiber.Ctx) error {

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", tbname, c.Params("id"))
	checkData, checkErr := db.Query(qyStr)
	if checkErr != nil {
		return resp.Failed(c, "Server Error", 500)
	}
	if checkData[0] == nil {
		return resp.Failed(c, "Data Not Found", 404)
	}

	// cmdStr := fmt.Sprintf("")

	return c.SendString("Test")
}

func (n *NewsHandler) DeleteNews(c *fiber.Ctx) error {

	// ReqHeader := c.GetReqHeaders()
	// AuthToken := strings.Split(ReqHeader["Authorization"], " ")[1]

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", tbname, c.Params("id"))
	checkData, checkErr := db.Query(qyStr)
	if checkErr != nil {
		return resp.Failed(c, "Server Error", 500)
	}
	if checkData[0] == nil {
		return resp.Failed(c, "Data Not Found", 404)
	}

	cmdStr := fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", tbname, c.Params("id"))
	resErr := db.Command(cmdStr)
	if resErr != nil {
		return resp.Failed(c, "Server Error", 500)
	}

	return resp.Success(c, nil, "Success Delete Data", 200)

}
