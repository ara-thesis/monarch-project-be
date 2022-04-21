package helper

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("PG_HOST")
	port     = os.Getenv("PG_PORT")
	user     = os.Getenv("PG_USER")
	password = os.Getenv("PG_PASS")
	dbname   = os.Getenv("PG_DB")
)

func Qy(query string, select_column ...interface{}) error {

	fmt.Print(host + " " + port + " " + user + " " + dbname)
	pgqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, dbErr := sql.Open("postgres", pgqlconn)
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()
	var (
		rows  *sql.Rows
		qyErr error
	)
	rows, qyErr = db.Query(query)
	if qyErr != nil {
		return qyErr
	}
	for rows.Next() {
		if scanErr := rows.Scan(select_column...); scanErr != nil {
			return scanErr
		}
	}
	return nil

}

func Cmd(query string, args ...interface{}) error {

	pgqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, dbErr := sql.Open("postgres", pgqlconn)
	if dbErr != nil {
		return dbErr
	}
	db.Exec("BEGIN")
	_, cmdErr := db.Exec(query, args[len(args)-1])
	if cmdErr != nil {
		db.Exec("ROLLBACK")
		db.Close()
		return cmdErr
	}
	db.Exec("COMMIT")
	db.Close()
	return nil

}
