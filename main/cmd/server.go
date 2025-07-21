package main

import (
	"echoProject/main/internal/models"
	"echoProject/main/internal/app/router"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
	dsn = "root:password@tcp(127.0.0.1:3306)/go_sample?charset=utf8mb4&parseTime=True&loc=Local"
)

func main() {
	dbinit()
	router.Init()
}

func dbinit() {
	// sqlhander と　conf でまとめれるか。
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
	}
	db.Migrator().CreateTable(models.User{})
}
