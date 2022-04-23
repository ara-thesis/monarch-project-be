package module

import (
	pg "github.com/ara-thesis/monarch-project-be/src/helper"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type NewsHandler struct {
}

func (n *NewsHandler) GetNews(c *fiber.Ctx) error {

	respMsg := make(map[string]interface{})

	resArr, resErr := pg.Query("SELECT * FROM newstb WHERE id='2'")

	if resErr != nil {

		respCode := 500

		respMsg["success"] = false
		respMsg["data"] = resErr.Error()
		respMsg["message"] = "SQL ERROR"
		respMsg["code"] = respCode

		return c.Status(respCode).JSON(respMsg)
	}

	respCode := 200

	respMsg["success"] = true
	respMsg["data"] = resArr
	respMsg["message"] = "Fetching News Data"
	respMsg["code"] = respCode

	return c.Status(respCode).JSON(respMsg)
}

func (n *NewsHandler) GetNewsById(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func (n *NewsHandler) AddNews(c *fiber.Ctx) error {

	return c.SendString("Test")
}

func (n *NewsHandler) EditNews(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func (n *NewsHandler) DeleteNews(c *fiber.Ctx) error {
	return c.SendString("Test")
}
