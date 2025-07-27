package main

import (
	di "echoProject/main/app/DI"
	"echoProject/main/app/router"
)

func main() {
	dbinit()
	router.Init()
}

func dbinit() {
	db, err := di.InitializeDB()
	if err != nil {
		panic(err)
	}
	db.Gorm_DB.DBinit()
}
