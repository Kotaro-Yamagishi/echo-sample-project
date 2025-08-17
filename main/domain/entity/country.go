package entity

import (
	"echoProject/main/domain/types"
	"time"
)

type Country struct {
	CountryID  uint16
	Country    types.CountryName
	LastUpdate time.Time
}

func NewCountry(countryName string) (Country, error) {
	return Country{
		Country:    types.CountryName(countryName),
		LastUpdate: time.Now(),
	}, nil
}

func NewValidatedCountry(countryName string) (Country, error) {
	if err := Validate(countryName); err != nil {
		return Country{}, err
	}

	return NewCountry(countryName)
}

func Validate(countryName string) error {
	if err := types.ValidateCountryName(countryName); err != nil {
		return err
	}
	return nil
}
