package types

import "errors"

type CountryName string

func ValidateCountryName(name string) error {
	if name == "" {
		return errors.New("country name is required")
	}
	return nil
}

func (c CountryName) String() string {
	return string(c)
}
