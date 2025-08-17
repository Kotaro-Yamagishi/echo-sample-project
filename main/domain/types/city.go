package types

import "errors"

type CityName string

func ValidateCityName(name string) error {
	if name == "" {
		return errors.New("city name is required")
	}
	return nil
}

func (c CityName) Validate() error {
	return ValidateCityName(string(c))
}

func (c CityName) String() string {
	return string(c)
}
