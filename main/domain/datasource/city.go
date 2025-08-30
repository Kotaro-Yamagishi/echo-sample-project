package datasource

import (
	"echoProject/domain/model"
)

type City interface {
	Select() model.CitySlice
}
