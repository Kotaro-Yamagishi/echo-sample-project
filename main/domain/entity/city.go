package entity

import "time"

type City struct {
	CityID     uint16    `json:"city_id"`
	City       string    `json:"city"`
	CountryID  uint16    `json:"country_id"`
	LastUpdate time.Time `json:"last_update"`
}
