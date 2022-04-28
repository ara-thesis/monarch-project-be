package main

import (
	"log"
	"os"

	env_helper "github.com/ara-thesis/monarch-project-be/src/helper"

	"github.com/ara-thesis/monarch-project-be/src/controller"
	"github.com/gofiber/fiber/v2"
)

func pathapi(app *fiber.App) {
	AccountHandler := new(controller.AccountHandler)
	NewsHandler := new(controller.NewsHandler)

	app.Get("/api/auth", AccountHandler.GetUserInfo)
	app.Post("/api/auth/regist/placemanager", AccountHandler.CreateUserPlaceManager)
	app.Post("/api/auth/regist/tourist", AccountHandler.CreateUserTourist)
	app.Post("/api/auth/login", AccountHandler.UserLogin)
	app.Put("/api/auth", AccountHandler.EditUser)
	app.Put("/api/auth/:id", AccountHandler.EditUserAsAdmin)
	app.Delete("/api/auth/:id", AccountHandler.DeleteUser)

	app.Get("/api/news", NewsHandler.GetNews)
	app.Get("/api/news/:id", NewsHandler.GetNewsById)
	app.Post("/api/news", NewsHandler.AddNews)
	app.Put("/api/news/:id", NewsHandler.EditNews)
	app.Delete("/api/news/:id", NewsHandler.DeleteNews)

}

func main() {
	env_helper.SetEnv()
	app := fiber.New()
	pathapi(app)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
