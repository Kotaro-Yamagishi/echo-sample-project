package datasource

import (
	"echoProject/domain/model"
)

type Country interface {
	Select() (model.CountrySlice, error)
	Insert(country *model.Country) error
}
