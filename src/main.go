package main

import (
	"log"

	"github.com/ara-thesis/monarch-project-be/src/helper"

	"github.com/ara-thesis/monarch-project-be/src/controller"
	"github.com/gofiber/fiber/v2"
)

func middleware(app *fiber.App) {

}

func pathstatic(app *fiber.App) {
	app.Static("/api/public/news", "./public/news")
	app.Static("/api/public/placeinfo", "./public/placeinfo")
}

func pathapi(app *fiber.App) {
	JwtHelper := new(helper.JwtHelper)
	AccountHandler := new(controller.AccountHandler)
	NewsHandler := new(controller.NewsHandler)
	PlaceInfoHandler := new(controller.PlaceInfoHandler)

	app.Get("/api/auth", AccountHandler.GetUserInfo)
	app.Post("/api/auth/regist/placemanager", AccountHandler.CreateUserPlaceManager)
	app.Post("/api/auth/regist/tourist", AccountHandler.CreateUserTourist)
	app.Post("/api/auth/login", AccountHandler.UserLogin)
	app.Put("/api/auth", AccountHandler.EditUser)
	app.Put("/api/auth/:id", AccountHandler.EditUserAsAdmin)
	app.Delete("/api/auth/:id", AccountHandler.DeleteUser)

	app.Get("/api/news", NewsHandler.GetNews)
	app.Get("/api/news/:id", NewsHandler.GetNewsById)
	app.Get("/api/news/admin", JwtHelper.VerifyToken, NewsHandler.GetNewsAdmin)
	app.Post("/api/news", JwtHelper.VerifyToken, NewsHandler.AddNews)
	app.Put("/api/news/:id", JwtHelper.VerifyToken, NewsHandler.EditNews)
	app.Delete("/api/news/:id", JwtHelper.VerifyToken, NewsHandler.DeleteNews)

	app.Get("/api/placeinfo", PlaceInfoHandler.GetPlaceInfo)
	app.Get("/api/placeinfo/:id", PlaceInfoHandler.GetPlaceInfoById)
	app.Get("/api/placeinfo/admin", JwtHelper.VerifyToken, PlaceInfoHandler.GetPlaceInfoAdmin)
	app.Put("/api/placeinfo", JwtHelper.VerifyToken, PlaceInfoHandler.AddAndEditPlaceInfoAdmin)
	app.Delete("/api/placeinfo/:userId", JwtHelper.VerifyToken, PlaceInfoHandler.DeletePlaceInfoAdmin)

	// app.Get("")

}

func main() {
	helper.SetEnv()
	app := fiber.New()
	pathstatic(app)
	middleware(app)
	pathapi(app)
	log.Fatal(app.Listen(":" + helper.GetEnv("PORT")))
}
