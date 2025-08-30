package country

import (
	"context"
	"database/sql"
	"echoProject/domain/datasource"
	"echoProject/domain/model"
	"echoProject/infra/things/sqlboiler"
	"fmt"

	"github.com/aarondl/sqlboiler/v4/boil"
)

type CountryImpl struct {
	db        *sql.DB
	sqlBoiler sqlboiler.SQLBoiler
}

func NewCountryDataSource(sqlBoiler sqlboiler.SQLBoiler) (datasource.Country, error) {
	db := sqlBoiler.ConnectDB()
	if db == nil {
		return nil, fmt.Errorf("database error: failed to connect to the database")
	}
	return &CountryImpl{db: db, sqlBoiler: sqlBoiler}, nil
}

func (ds *CountryImpl) Select() (model.CountrySlice, error) {
	countries, err := model.Countries().All(context.Background(), ds.db)
	if err != nil {
		return nil, fmt.Errorf("database error: failed to select countries: %w", err)
	}
	return countries, nil
}

func (ds *CountryImpl) Insert(country *model.Country) error {
	return country.Insert(context.Background(), ds.db, boil.Infer())
}
