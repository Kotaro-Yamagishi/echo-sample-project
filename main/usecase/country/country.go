package country

import (
	"echoProject/main/domain/entity"
	"echoProject/main/domain/repository"
	"echoProject/main/domain/usecase"
	"fmt"
)

type CountryImpl struct {
	repo repository.Country
}

func NewCountryService(repo repository.Country) usecase.Country {
	return &CountryImpl{repo: repo}
}

func (s *CountryImpl) Select() ([]entity.Country, error) {
	countries, err := s.repo.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to select countries: %w", err)
	}
	return countries, nil
}

func (s *CountryImpl) Insert(country entity.Country) error {
	if err := s.repo.Insert(country); err != nil {
		return fmt.Errorf("failed to insert country: %w", err)
	}
	return nil
}
