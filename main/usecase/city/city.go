package uccity

import (
	"echoProject/main/domain/entity"
	repoIF "echoProject/main/domain/repository"
	ucIF "echoProject/main/domain/usecase"
)

type CityImpl struct {
	repo repoIF.City
}

func NewCityService(repo repoIF.City) ucIF.City {
	return &CityImpl{repo: repo}
}

func (s *CityImpl) Select() []entity.City {
	return s.repo.Select()
}
