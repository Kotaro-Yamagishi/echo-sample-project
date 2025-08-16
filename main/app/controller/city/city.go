package ctrcity

import (
	"echoProject/main/domain/controller"
	"echoProject/main/domain/entity"
	"echoProject/main/domain/usecase"
)

type CityImpl struct {
	uc ucIF.City
}

func NewCityController(uc ucIF.City) ctrIF.City {
	return &CityImpl{uc: uc}
}

func (c *CityImpl) GetCity() []entity.City {
	return c.uc.Select()
}
