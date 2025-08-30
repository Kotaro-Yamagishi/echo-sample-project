package city

import (
	"echoProject/domain/datasource"
	"echoProject/domain/entity"
	"echoProject/domain/repository"
)

type CityImpl struct {
	ds datasource.City
}

func NewCityRepository(ds datasource.City) repository.City {
	return &CityImpl{ds: ds}
}

func (db *CityImpl) Select() []entity.City {
	cities := db.ds.Select()
	var entities []entity.City
	for _, c := range cities {
		entities = append(entities, entity.City{
			CityID:     c.CityID,
			City:       string(c.City),
			CountryID:  c.CountryID,
			LastUpdate: c.LastUpdate,
		})
	}
	return entities
}
