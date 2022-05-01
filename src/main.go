package main

import (
	"log"
	"os"

	"github.com/ara-thesis/monarch-project-be/src/helper"
	env_helper "github.com/ara-thesis/monarch-project-be/src/helper"

	"github.com/ara-thesis/monarch-project-be/src/controller"
	"github.com/gofiber/fiber/v2"
)

func pathapi(app *fiber.App) {
	AccountHandler := new(controller.AccountHandler)
	NewsHandler := new(controller.NewsHandler)
	JwtHelper := new(helper.JwtHelper)

	app.Get("/api/auth", AccountHandler.GetUserInfo)
	app.Post("/api/auth/regist/placemanager", AccountHandler.CreateUserPlaceManager)
	app.Post("/api/auth/regist/tourist", AccountHandler.CreateUserTourist)
	app.Post("/api/auth/login", AccountHandler.UserLogin)
	app.Put("/api/auth", AccountHandler.EditUser)
	app.Put("/api/auth/:id", AccountHandler.EditUserAsAdmin)
	app.Delete("/api/auth/:id", AccountHandler.DeleteUser)

	app.Get("/api/news", JwtHelper.VerifyToken, NewsHandler.GetNews)
	app.Get("/api/news/:id", JwtHelper.VerifyToken, NewsHandler.GetNewsById)
	app.Post("/api/news", JwtHelper.VerifyToken, NewsHandler.AddNews)
	app.Put("/api/news/:id", JwtHelper.VerifyToken, NewsHandler.EditNews)
	app.Delete("/api/news/:id", JwtHelper.VerifyToken, NewsHandler.DeleteNews)

}

func main() {
	env_helper.SetEnv()
	app := fiber.New()
	pathapi(app)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
