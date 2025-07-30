package main

import (
	di "echoProject/main/app/DI"
	"echoProject/main/app/router"

	"github.com/aarondl/sqlboiler/v4/boil"
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
	boil.SetDB(db.SqlBoiler.ConnectDB())
}
