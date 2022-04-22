package news_handler

import (
	pg "github.com/ara-thesis/monarch-project-be/src/helper/database"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func GetNews(c *fiber.Ctx) error {

	resArr, resErr := pg.Query("SELECT * FROM ewstb")

	if resErr != nil {
		return c.Status(500).SendString(resErr.Error())
	}

	return c.JSON(resArr)
}

func GetNewsById(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func AddNews(c *fiber.Ctx) error {

	return c.SendString("Test")
}

func EditNews(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func DeleteNews(c *fiber.Ctx) error {
	return c.SendString("Test")
}
