package city

import (
	"echoProject/domain/entity"
	"echoProject/domain/repository"
	"echoProject/domain/usecase"
)

type CityImpl struct {
	repo repository.City
}

func NewCityService(repo repository.City) usecase.City {
	return &CityImpl{repo: repo}
}

func (s *CityImpl) Select() []entity.City {
	return s.repo.Select()
}
