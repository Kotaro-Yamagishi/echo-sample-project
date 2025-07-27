package initializer

import (
	"echoProject/main/domain/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
	dsn = "root:password@tcp(127.0.0.1:3306)/go_sample?charset=utf8mb4&parseTime=True&loc=Local"
)

type GormDBImpl struct {}

func NewGormDBImpl() Gorm_DB {
	return &GormDBImpl{}
}

func (g *GormDBImpl) DBinit() {
	// sqlhander と　conf でまとめれるか。
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
	}
	db.Migrator().CreateTable(entity.User{})
}