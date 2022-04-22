package main

import (
	"log"
	"os"

	env_helper "github.com/ara-thesis/monarch-project-be/src/helper/environment"

	news_handler "github.com/ara-thesis/monarch-project-be/src/module/news"
	"github.com/gofiber/fiber/v2"
)

func pathapi(app *fiber.App) {
	app.Get("/api/news", news_handler.GetNews)
	app.Get("/api/news/:id", news_handler.GetNewsById)
	app.Post("/api/news", news_handler.AddNews)
	app.Put("/api/news/:id", news_handler.EditNews)
	app.Delete("/api/news/:id", news_handler.DeleteNews)
}

func main() {
	env_helper.SetEnv()
	app := fiber.New()
	pathapi(app)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
