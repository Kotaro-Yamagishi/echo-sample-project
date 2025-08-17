package datasource

import (
	"echoProject/main/domain/model"
)

type Country interface {
	Select() (model.CountrySlice, error)
	Insert(country *model.Country) error
}
