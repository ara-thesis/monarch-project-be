package helper

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

// var (
// 	host     = "localhost"
// 	port     = "5432"
// 	user     = "postgres"
// 	password = "Raflis2001"
// 	dbname   = "monarch-thesis"
// )

var (
	host     = os.Getenv("PG_HOST")
	port     = os.Getenv("PG_PORT")
	user     = os.Getenv("PG_USER")
	password = os.Getenv("PG_PASS")
	dbname   = os.Getenv("PG_DB")
)

type PgHelper struct{}

// sql query helper
func (pg *PgHelper) Query(query string, args ...interface{}) ([]interface{}, error) {

	connConfig := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	pgqlconn, errConn := pgx.ParseDSN(connConfig)
	if errConn != nil {
		return nil, errConn
	}
	db, dbErr := pgx.Connect(pgqlconn)
	if dbErr != nil {
		return nil, dbErr
	}
	rows, qyErr := db.Query(query, args...)
	if qyErr != nil {
		db.Close()
		return nil, qyErr
	}
	// valArray := make(map[int]map[string]interface{})
	valArray := make([]interface{}, 1)
	colName := rows.FieldDescriptions()
	valArrCount := 0
	for rows.Next() {
		valResult, valErr := rows.Values()
		tempArr := make(map[string]interface{})
		if valErr != nil {
			db.Close()
			return nil, valErr
		}
		for id, val := range valResult {
			tempArr[colName[id].Name] = val
		}
		if valArrCount == 0 {
			valArray[0] = tempArr
		} else {
			valArray = append(valArray, tempArr)
		}
		valArrCount++
	}
	db.Close()
	return valArray, nil

}

// sql command helper
func (pg *PgHelper) Command(query string, args ...interface{}) error {

	connConfig := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	pgqlconn, errConn := pgx.ParseDSN(connConfig)
	if errConn != nil {
		return errConn
	}
	db, dbErr := pgx.Connect(pgqlconn)
	if dbErr != nil {
		return dbErr
	}
	tx, txErr := db.Begin()
	if txErr != nil {
		db.Close()
		return txErr
	}
	_, cmdErr := tx.Exec(query, args...)
	if cmdErr != nil {
		tx.Rollback()
		db.Close()
		return cmdErr
	}
	commErr := tx.Commit()
	if commErr != nil {
		tx.Rollback()
		db.Close()
		return commErr
	}
	db.Close()
	return nil

}
