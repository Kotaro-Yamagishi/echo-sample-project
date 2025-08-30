package country

import (
	"echoProject/domain/datasource"
	"echoProject/domain/entity"
	"echoProject/domain/model"
	"echoProject/domain/repository"
	"echoProject/domain/types"
	"fmt"
)

type CountryImpl struct {
	ds datasource.Country
}

func NewCountryRepository(ds datasource.Country) repository.Country {
	return &CountryImpl{ds: ds}
}

func (db *CountryImpl) Select() ([]entity.Country, error) {
	countries, err := db.ds.Select()
	if err != nil {
		return nil, fmt.Errorf("database error: failed to select countries: %w", err)
	}
	var entities []entity.Country
	for _, c := range countries {
		entities = append(entities, entity.Country{
			CountryID:  c.CountryID,
			Country:    types.CountryName(c.Country),
			LastUpdate: c.LastUpdate,
		})
	}
	return entities, nil
}

func (db *CountryImpl) Insert(country entity.Country) error {
	modelCountry := &model.Country{
		Country:    string(country.Country),
		LastUpdate: country.LastUpdate,
	}

	if err := db.ds.Insert(modelCountry); err != nil {
		return fmt.Errorf("database error: failed to insert country: %w", err)
	}
	return nil
}
