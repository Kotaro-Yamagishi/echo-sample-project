package country

import (
	"echoProject/main/domain/datasource"
	"echoProject/main/domain/entity"
	"echoProject/main/domain/model"
	"echoProject/main/domain/repository"
)

type CountryImpl struct {
	ds datasource.Country
}

func NewCountryRepository(ds datasource.Country) repository.Country {
	return &CountryImpl{ds: ds}
}

func (db *CountryImpl) Select() []entity.Country {
	countries := db.ds.Select()
	var entities []entity.Country
	for _, c := range countries {
		entities = append(entities, entity.Country{
			CountryID:  c.CountryID,
			Country:    c.Country,
			LastUpdate: c.LastUpdate,
		})
	}
	return entities
}

func (db *CountryImpl) Insert(country entity.Country) error {
	modelCountry := &model.Country{
		Country:    country.Country,
		LastUpdate: country.LastUpdate,
	}
	return db.ds.Insert(modelCountry)
}
