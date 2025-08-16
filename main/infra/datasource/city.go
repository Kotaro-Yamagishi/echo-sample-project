package ds

import (
	"context"
	"database/sql"
	"echoProject/main/domain/datasource"
	"echoProject/main/domain/model"
	"echoProject/main/infra/things/sqlboiler"
	"log"
)

type CityDs struct {
	db        *sql.DB
	sqlBoiler sqlboiler.SQLBoiler
}

func NewCityDataSource(sqlBoiler sqlboiler.SQLBoiler) dsIF.City {
	db := sqlBoiler.ConnectDB()
	if db == nil {
		log.Fatal("Failed to connect to the database")
	}
	return &CityDs{db: db, sqlBoiler: sqlBoiler}
}

func (ds *CityDs) Select() model.CitySlice {
	cities, err := model.Cities().All(context.Background(), ds.db)
	if err != nil {
		log.Fatal(err)
	}
	return cities
}
