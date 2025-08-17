package country

import (
	"context"
	"database/sql"
	"echoProject/main/domain/datasource"
	"echoProject/main/domain/model"
	"echoProject/main/infra/things/sqlboiler"
	"log"

	"github.com/aarondl/sqlboiler/v4/boil"
)

type CountryImpl struct {
	db        *sql.DB
	sqlBoiler sqlboiler.SQLBoiler
}

func NewCountryDataSource(sqlBoiler sqlboiler.SQLBoiler) datasource.Country {
	db := sqlBoiler.ConnectDB()
	if db == nil {
		log.Fatal("Failed to connect to the database")
	}
	return &CountryImpl{db: db, sqlBoiler: sqlBoiler}
}

func (ds *CountryImpl) Select() model.CountrySlice {
	countries, err := model.Countries().All(context.Background(), ds.db)
	if err != nil {
		log.Fatal(err)
	}
	return countries
}

func (ds *CountryImpl) Insert(country *model.Country) error {
	return country.Insert(context.Background(), ds.db, boil.Infer())
}
