package controller

import (
	"fmt"
	"time"

	"github.com/ara-thesis/monarch-project-be/src/model"
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

// fetch news by id
func (n *NewsHandler) GetNewsById(c *fiber.Ctx) error {

	// ReqHeader := c.GetReqHeaders()
	// AuthToken := strings.Split(ReqHeader["Authorization"], " ")[1]

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", tbname["news"], c.Params("id"))
	resQy, resErr := db.Query(qyStr)
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

	model := new(model.NewsModel)
	uuid := uuid.New()

	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	cmdStr := fmt.Sprintf(`INSERT INTO %s(id, title, article, status, draft_status, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7)`, tbname["news"])
	resErr := db.Command(cmdStr, uuid, model.Title, model.Article, model.Status, model.Draft_status, time.Now(), time.Now())
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "Data added")
}

// edit news by id
func (n *NewsHandler) EditNews(c *fiber.Ctx) error {

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", tbname["news"], c.Params("id"))
	checkData, checkErr := db.Query(qyStr)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// cmdStr := fmt.Sprintf("")

	return c.SendString("Test")
}

// delete news by id
func (n *NewsHandler) DeleteNews(c *fiber.Ctx) error {

	// ReqHeader := c.GetReqHeaders()
	// AuthToken := strings.Split(ReqHeader["Authorization"], " ")[1]

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", tbname["news"], c.Params("id"))
	checkData, checkErr := db.Query(qyStr)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	cmdStr := fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", tbname["news"], c.Params("id"))
	resErr := db.Command(cmdStr)
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "Success Delete Data")

}
