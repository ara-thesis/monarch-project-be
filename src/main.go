package main

import (
	"log"
	"os"

	news_handler "github.com/ara-thesis/monarch-project-be/src/module/news"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func pathapi(app *fiber.App) {
	app.Get("/api/news/v1", news_handler.GetNews)
	app.Get("/api/news/v1/:id", news_handler.GetNewsById)
	app.Post("/api/news/v1", news_handler.AddNews)
	app.Put("/api/news/v1/:id", news_handler.EditNews)
	app.Delete("/api/news/v1/:id", news_handler.DeleteNews)
}

func main() {
	env_err := godotenv.Load("./.env")
	if env_err != nil {
		log.Fatalf("failed to load env file: %s", env_err)
	}
	app := fiber.New()
	pathapi(app)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
