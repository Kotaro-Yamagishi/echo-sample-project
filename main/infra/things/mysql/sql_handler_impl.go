package mysql

// これはgormパッケージだね

import (
	"context"
	"database/sql"
	"echoProject/main/domain/model"
	"fmt"
	"log"
)

type SqlHandlerImpl struct {
	db *sql.DB
}

func NewSqlHandler() SqlHandler {
	dsn := "root:password@tcp(127.0.0.1:3306)/go_sample?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }

    // DB接続テスト
    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }

	sqlHandler := new(SqlHandlerImpl)
	sqlHandler.db = db
	return sqlHandler
}

func (handler *SqlHandlerImpl) Create(obj interface{}) {
	// handler.db.Create(obj)
}

func (handler *SqlHandlerImpl) FindAll(obj interface{}) {
	users, err := model.Users().All(context.Background(), handler.db)
	if err != nil {
        log.Fatal(err)
    }
	fmt.Println("Users:", users)
	obj = users
}

func (handler *SqlHandlerImpl) DeleteById(obj interface{}, id string) {
	// handler.db.Delete(obj, id)
}