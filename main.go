package main

import (
	"log"

	"github.com/ara-thesis/monarch-project-be/controller"
	"github.com/ara-thesis/monarch-project-be/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func middleware(app *fiber.App) {
	app.Use(cors.New())
}

func pathstatic(app *fiber.App) {
	app.Static("/api/public/dummy/", "./public/dummy")
	app.Static("/api/public/news/", "./public/news")
	app.Static("/api/public/placeinfo/", "./public/placeinfo")
	app.Static("/api/public/banner/", "./public/banner")
	app.Static("/api/public/payment/", "./public/payment/")
}

func pathapi(app *fiber.App) {
	JwtHelper := new(helper.JwtHelper)
	AccountHandler := &controller.AccountHandler{
		Tbname: "userinfo",
	}
	NewsHandler := &controller.NewsHandler{
		Tbname:           "newstb",
		Tbname_placeinfo: "placeinfotb",
	}
	PlaceInfoHandler := &controller.PlaceInfoHandler{
		Tbname:     "placeinfotb",
		Tbname_img: "placeinfo_imgtb",
	}
	BannerHandler := &controller.BannerHandler{
		Tbname: "bannertb",
	}
	ReviewHandler := &controller.ReviewHandler{
		Tbname: "reviewtb",
		// Tbname_rating: ,
	}
	TicketHandler := &controller.TicketHandler{
		Tbname:        "tickettb",
		Tbname_place:  "placeinfotb",
		Tbname_bought: "ticketboughttb",
	}
	CartHandler := &controller.CartHandler{
		Tbname:        "ticketcarttb",
		Tbname_ticket: "tickettb",
	}
	PaymentHandler := &controller.PaymentController{
		Tbname:              "paymenttb",
		Tbname_user:         "userinfo",
		Tbname_cart:         "ticketcarttb",
		Tbname_ticketbought: "ticketboughttb",
	}
	ItineraryHandler := &controller.ItineraryHandler{
		Tbname:           "itinerarytb",
		Tbname_item:      "itineraryitemtb",
		Tbname_placeinfo: "placeinfotb",
	}

	// authorization handler
	app.Get("/api/auth/me", JwtHelper.VerifyToken, AccountHandler.GetUserInfo)
	app.Get("/api/auth/list/placemanager", JwtHelper.VerifyToken, AccountHandler.GetUserListPM)
	app.Get("/api/auth/list/tourist", JwtHelper.VerifyToken, AccountHandler.GetUserListTourist)
	// app.Get("/api/auth/:id", JwtHelper.VerifyToken, AccountHandler.GetUserInfo)
	app.Post("/api/auth/regist/placemanager", AccountHandler.CreateUserPlaceManager)
	app.Post("/api/auth/regist/admin", AccountHandler.CreateUserAdmin)
	app.Post("/api/auth/regist/tourist", JwtHelper.VerifyToken, AccountHandler.CreateUserTourist)
	app.Post("/api/auth/login", AccountHandler.UserLogin)
	app.Put("/api/auth", JwtHelper.VerifyToken, AccountHandler.EditUser)
	app.Put("/api/auth/:id", JwtHelper.VerifyToken, AccountHandler.EditUserAsAdmin)
	app.Delete("/api/auth/:id", JwtHelper.VerifyToken, AccountHandler.DeleteUser)

	// news handler
	app.Get("/api/news", NewsHandler.GetNews)
	app.Get("/api/news/:id", NewsHandler.GetNewsById)
	app.Get("/api/news/list/admin", JwtHelper.VerifyToken, NewsHandler.GetNewsAdmin)
	app.Post("/api/news", JwtHelper.VerifyToken, NewsHandler.AddNews)
	app.Put("/api/news/:id", JwtHelper.VerifyToken, NewsHandler.EditNews)
	app.Delete("/api/news/:id", JwtHelper.VerifyToken, NewsHandler.DeleteNews)

	// place info handler
	app.Get("/api/placeinfo", PlaceInfoHandler.GetPlaceInfo)
	app.Get("/api/placeinfo/:id", PlaceInfoHandler.GetPlaceInfoById)
	app.Get("/api/placeinfo/show/admin", JwtHelper.VerifyToken, PlaceInfoHandler.GetPlaceInfoAdmin)
	app.Put("/api/placeinfo", JwtHelper.VerifyToken, PlaceInfoHandler.UpdatePlaceInfoAdmin)
	app.Delete("/api/placeinfo/:userId", JwtHelper.VerifyToken, PlaceInfoHandler.DeletePlaceInfoAdmin)

	// banner handler
	app.Get("/api/banner", BannerHandler.GetBanners)
	app.Get("/api/banner/:id", BannerHandler.GetBannerById)
	app.Post("/api/banner", JwtHelper.VerifyToken, BannerHandler.AddBanner)
	app.Put("/api/banner/:id", JwtHelper.VerifyToken, BannerHandler.EditBanner)
	app.Delete("/api/banner/:id", JwtHelper.VerifyToken, BannerHandler.DeleteBanner)

	// ticket handler
	app.Get("/api/ticket", TicketHandler.GetTicketTourist)
	app.Get("/api/ticket/:id", TicketHandler.GetTicketById)
	app.Get("/api/ticket/list/tourist", JwtHelper.VerifyToken, TicketHandler.GetTicketBoughtTourist)
	app.Get("/api/ticket/list/admin", JwtHelper.VerifyToken, TicketHandler.GetTicketAdmin)
	app.Post("/api/ticket", JwtHelper.VerifyToken, TicketHandler.AddTicket)
	app.Put("/api/ticket/:id", JwtHelper.VerifyToken, TicketHandler.EditTicket)
	app.Put("/api/ticket/redeem/tourist", JwtHelper.VerifyToken, TicketHandler.RedeemTicket)
	app.Delete("/api/ticket/:id", JwtHelper.VerifyToken, TicketHandler.DeleteTicket)

	// review handler
	app.Get("/api/review", ReviewHandler.GetComment)
	app.Post("/api/review", JwtHelper.VerifyToken, ReviewHandler.AddComment)
	app.Delete("/api/review/:id", JwtHelper.VerifyToken, ReviewHandler.DeleteCommentAdmin)

	// cart handler
	app.Get("/api/cart", JwtHelper.VerifyToken, CartHandler.GetCart)
	app.Post("/api/cart", JwtHelper.VerifyToken, CartHandler.AddToCart)
	app.Delete("/api/cart/:id", JwtHelper.VerifyToken, CartHandler.RemoveFromCart)

	// payment handler
	app.Get("/api/payment", JwtHelper.VerifyToken, PaymentHandler.GetPurchaseConfirm)
	app.Get("/api/payment/:id", JwtHelper.VerifyToken, PaymentHandler.GetPurchaseConfirmById)
	app.Put("/api/payment/:id", JwtHelper.VerifyToken, PaymentHandler.AcceptPurchaseConfirmation)
	app.Delete("/api/payment/:id", JwtHelper.VerifyToken, PaymentHandler.DenyPurchaseConfirmation)
	app.Post("/api/payment/cart/buy", JwtHelper.VerifyToken, PaymentHandler.PayCart)

	// itinerary handler
	app.Get("/api/itinerary", JwtHelper.VerifyToken, ItineraryHandler.GetItinerary)
	app.Get("/api/itinerary/public/search", JwtHelper.VerifyToken, ItineraryHandler.GetItineraryPublic)
	app.Get("/api/itinerary/:id", JwtHelper.VerifyToken, ItineraryHandler.GetItineraryById)
	app.Post("/api/itinerary", JwtHelper.VerifyToken, ItineraryHandler.CreateItinerary)
	app.Put("/api/itinerary/:id", JwtHelper.VerifyToken, ItineraryHandler.UpdateItinerary)
	app.Delete("/api/itinerary/:id", JwtHelper.VerifyToken, ItineraryHandler.DeleteItinerary)

}

func main() {
	helper.SetEnv()
	app := fiber.New()
	middleware(app)
	pathstatic(app)
	pathapi(app)
	log.Fatal(app.Listen(":" + helper.GetEnv("PORT")))
}
