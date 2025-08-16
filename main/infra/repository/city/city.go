package repocity

import (
	dsIF "echoProject/main/domain/datasource"
	"echoProject/main/domain/entity"
	repoIF "echoProject/main/domain/repository"
)

type CityImpl struct {
	ds dsIF.City
}

func NewCityRepository(ds dsIF.City) repoIF.City {
	return &CityImpl{ds: ds}
}

func (db *CityImpl) Select() []entity.City {
	cities := db.ds.Select()
	var entities []entity.City
	for _, c := range cities {
		entities = append(entities, entity.City{
			CityID:     c.CityID,
			City:       c.City,
			CountryID:  c.CountryID,
			LastUpdate: c.LastUpdate,
		})
	}
	return entities
}
