package controller

import (
	"github.com/ara-thesis/monarch-project-be/src/helper"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type NewsHandler struct{}

func (n *NewsHandler) GetNews(c *fiber.Ctx) error {

	pg := new(helper.PgHelper)
	respHelper := new(helper.ResponseHelper)

	resArr, resErr := pg.Query("SELECT * FROM newstb WHERE id='2'")

	if resErr != nil {

		respMsg, respCode := respHelper.CreateResponse(resArr, resErr)

		return c.Status(respCode).JSON(respMsg)
	}

	respMsg, respCode := respHelper.CreateResponse(resArr, resErr)

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
