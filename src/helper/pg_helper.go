package helper

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func InitDB() {
	db, err = sql.Open("postgres", os.Getenv("DB_URI"))
	if err != nil {
		log.Fatal(err)
	}
}

// func Qy(query string, args ...string) string, error {
// 	row, err := db.Query(query, args)
// 	if err != nil {
// 		return "error"
// 	}
// 	return row.Columns()
// }

func Cmd(command string) {

}
