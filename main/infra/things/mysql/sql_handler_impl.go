package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SqlHandlerImpl struct {
	db *gorm.DB
}

func NewSqlHandler() SqlHandler {
	dsn := "root:password@tcp(127.0.0.1:3306)/go_sample?charset=utf8mb4&parseTime=True&loc=Local"
	open, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	sqlHandler := new(SqlHandlerImpl)
	sqlHandler.db = open
	return sqlHandler
}

func (handler *SqlHandlerImpl) Create(obj interface{}) {
	handler.db.Create(obj)
}

func (handler *SqlHandlerImpl) FindAll(obj interface{}) {
	handler.db.Find(obj)
}

func (handler *SqlHandlerImpl) DeleteById(obj interface{}, id string) {
	handler.db.Delete(obj, id)
}