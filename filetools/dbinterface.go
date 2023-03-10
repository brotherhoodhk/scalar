package filetools

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func connectdatabase(dburl string) (*sql.DB, error) {
	dbcon, err := sql.Open("mysql", dburl)
	if err != nil {
		return nil, err
	}
	return dbcon, nil
}
