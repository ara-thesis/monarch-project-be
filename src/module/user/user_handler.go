package user_handler

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func InitNews() {

	connStr := os.Getenv("DB_URI")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

}

func GetNews(c *fiber.Ctx) error {

	rows, err := db.Query("SELECT * FROM newstb")
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(rows)
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
