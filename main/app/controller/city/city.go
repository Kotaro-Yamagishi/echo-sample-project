package city

import (
	"echoProject/main/domain/controller"
	"echoProject/main/domain/entity"
	"echoProject/main/domain/usecase"
)

type CityImpl struct {
	uc usecase.City
}

func NewCityController(uc usecase.City) controller.City {
	return &CityImpl{uc: uc}
}

// @Summary Get all cities
// @Description Retrieve all cities from the database
// @Tags cities
// @Accept json
// @Produce json
// @Success 200 {array} entity.City
// @Router /cities [get]
func (c *CityImpl) GetCity() []entity.City {
	return c.uc.Select()
}
