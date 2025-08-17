package datasource

import (
	"echoProject/main/domain/model"
)

type Country interface {
	Select() model.CountrySlice
	Insert(country *model.Country) error
}
