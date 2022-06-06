package main

import (
	"log"

	"github.com/ara-thesis/monarch-project-be/src/helper"

	"github.com/ara-thesis/monarch-project-be/src/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func middleware(app *fiber.App) {
	app.Use(cors.New())
}

func pathstatic(app *fiber.App) {
	app.Static("/api/public/news/", "./public/news")
	app.Static("/api/public/placeinfo/", "./public/placeinfo")
	app.Static("/api/public/banner/", "./public/banner")
}

func pathapi(app *fiber.App) {
	JwtHelper := new(helper.JwtHelper)
	AccountHandler := new(controller.AccountHandler)
	NewsHandler := new(controller.NewsHandler)
	PlaceInfoHandler := new(controller.PlaceInfoHandler)
	BannerHandler := new(controller.BannerHandler)
	ReviewHandler := new(controller.ReviewHandler)

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

	app.Get("/api/banner", BannerHandler.GetBanners)
	app.Get("/api/banner/:id", BannerHandler.GetBannerById)
	app.Post("/api/banner", JwtHelper.VerifyToken, BannerHandler.AddBanner)
	app.Put("/api/banner/:id", JwtHelper.VerifyToken, BannerHandler.EditBanner)
	app.Delete("/api/banner/:id", JwtHelper.VerifyToken, BannerHandler.DeleteBanner)

	app.Get("/api/review", ReviewHandler.GetComment)
	app.Post("/api/review", JwtHelper.VerifyToken, ReviewHandler.AddComment)
	app.Delete("/api/review/:id", JwtHelper.VerifyToken, ReviewHandler.DeleteCommentAdmin)

}

func main() {
	helper.SetEnv()
	app := fiber.New()
	middleware(app)
	pathstatic(app)
	pathapi(app)
	log.Fatal(app.Listen(":" + helper.GetEnv("PORT")))
}
