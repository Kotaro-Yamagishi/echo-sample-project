package dsIF

import (
	"echoProject/main/domain/model"
)

type City interface {
	Select() model.CitySlice
}
