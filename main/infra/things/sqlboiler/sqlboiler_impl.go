package sqlboiler

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type SQLBoilerImpl struct {
}

func NewSQLBoilerImpl() SQLBoiler {
	sqlboiler := new(SQLBoilerImpl)
    return  sqlboiler
}

func (s *SQLBoilerImpl) ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/go_sample?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
