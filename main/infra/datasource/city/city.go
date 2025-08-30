package city

import (
	"context"
	"database/sql"
	"echoProject/domain/datasource"
	"echoProject/domain/model"
	"echoProject/infra/things/sqlboiler"
	"log"
)

type CityImpl struct {
	db        *sql.DB
	sqlBoiler sqlboiler.SQLBoiler
}

func NewCityDataSource(sqlBoiler sqlboiler.SQLBoiler) datasource.City {
	db := sqlBoiler.ConnectDB()
	if db == nil {
		log.Fatal("Failed to connect to the database")
	}
	return &CityImpl{db: db, sqlBoiler: sqlBoiler}
}

func (ds *CityImpl) Select() model.CitySlice {
	cities, err := model.Cities().All(context.Background(), ds.db)
	if err != nil {
		log.Fatal(err)
	}
	return cities
}
