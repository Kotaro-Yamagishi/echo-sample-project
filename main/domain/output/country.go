package output

import (
	"echoProject/main/domain/types"
	"time"
)

type Country struct {
	CountryID  uint16
	Country    types.CountryName
	LastUpdate time.Time
}

func NewCountry(countryID uint16, countryName types.CountryName, lastUpdate time.Time) Country {
	return Country{
		CountryID:  countryID,
		Country:    countryName,
		LastUpdate: lastUpdate,
	}
}
