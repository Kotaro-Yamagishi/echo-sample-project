package entity

import "time"

type Country struct {
	CountryID  uint16
	Country    string
	LastUpdate time.Time
}

// NewCountry creates a new Country entity
func NewCountry(country string) Country {
	return Country{
		Country:    country,
		LastUpdate: time.Now(),
	}
}

// NewCountryWithID creates a new Country entity with ID
func NewCountryWithID(countryID uint16, country string, lastUpdate time.Time) Country {
	return Country{
		CountryID:  countryID,
		Country:    country,
		LastUpdate: lastUpdate,
	}
}
