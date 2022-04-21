package news_handler

import (

	// "github.com/ara-thesis/monarch-project-be/src/helper"

	pgHelper "github.com/ara-thesis/monarch-project-be/src/helper"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type data struct {
	id    string
	title string
}

func GetNews(c *fiber.Ctx) error {

	var value data

	pgHelper.Qy("SELECT id, title FROM newstb", &value.id, &value.title)

	// result, _ := json.Marshal(resultsql)

	return c.SendString("" + value.id + ":" + value.title)
}

func GetNewsById(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func AddNews(c *fiber.Ctx) error {

	pgHelper.Cmd("")

	return c.SendString("Test")
}

func EditNews(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func DeleteNews(c *fiber.Ctx) error {
	return c.SendString("Test")
}

// func randID(c *fiber.Ctx) error {
// 	uuid := uuid.New()

// 	return c.SendString(uuid.String())
// }
