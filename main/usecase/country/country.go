package country

import (
	"echoProject/main/domain/entity"
	"echoProject/main/domain/repository"
	"echoProject/main/domain/usecase"
)

type CountryImpl struct {
	repo repository.Country
}

func NewCountryService(repo repository.Country) usecase.Country {
	return &CountryImpl{repo: repo}
}

func (s *CountryImpl) Select() []entity.Country {
	return s.repo.Select()
}

func (s *CountryImpl) Insert(country entity.Country) error {
	return s.repo.Insert(country)
}
