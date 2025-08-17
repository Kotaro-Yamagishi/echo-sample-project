package country

import (
	"echoProject/main/domain/controller"
	"echoProject/main/domain/entity"
	"echoProject/main/domain/usecase"

	"github.com/labstack/echo/v4"
)

type CountryImpl struct {
	uc usecase.Country
}

func NewCountryController(uc usecase.Country) controller.Country {
	return &CountryImpl{uc: uc}
}

// @Summary Get all countries
// @Description Retrieve all countries from the database
// @Tags countries
// @Accept json
// @Produce json
// @Success 200 {array} entity.Country
// @Router /countries [get]
func (c *CountryImpl) GetCountry() []entity.Country {
	return c.uc.Select()
}

// @Summary Create a new country
// @Description Add a new country to the database
// @Tags countries
// @Accept json
// @Produce json
// @Param country body object true "Country object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /countries [post]
func (c *CountryImpl) Create(ctx echo.Context) error {
	var request struct {
		Country string `json:"country"`
	}
	if err := ctx.Bind(&request); err != nil {
		return err
	}

	country := entity.NewCountry(request.Country)
	return c.uc.Insert(country)
}
