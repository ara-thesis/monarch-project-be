package news_handler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	// "github.com/ara-thesis/monarch-project-be/src/helper"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func Init() {

	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASS")
	dbname := os.Getenv("PG_DB")

	pgconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", pgconn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("DB CONNECTED")

}

func GetNews(c *fiber.Ctx) error {

	rows, err1 := db.Query("SELECT * FROM newstb")

	if err1 != nil {
		log.Fatal(err)
	}

	res, err2 := rows.Columns()

	if err2 != nil {
		log.Fatal(err)
	}

	// result := map[string]string{}

	// for rows.Next() {
	// 	var (
	// 		id int64
	// 		name string
	// 	)

	// 	err := rows.Scan()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.
	// }

	rows.Close()

	// return c.JSON(rows)
	return c.SendString(strings.Join(res, " "))
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

// func randID(c *fiber.Ctx) error {
// 	uuid := uuid.New()

// 	return c.SendString(uuid.String())
// }
