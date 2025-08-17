package country

import (
	"echoProject/main/domain/controller"
	"echoProject/main/domain/entity"
	"echoProject/main/domain/output"
	"echoProject/main/domain/usecase"
	"net/http"

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
// @Success 200 {object} output.SuccessResponse
// @Failure 500 {object} output.ErrorResponse
// @Router /countries [get]
func (c *CountryImpl) GetCountry(ctx echo.Context) error {
	countries, err := c.uc.Select()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, output.NewErrorResponse(
			http.StatusInternalServerError,
			"failed to get countries",
		))
	}

	return ctx.JSON(http.StatusOK, output.NewSuccessResponse(countries))
}

// @Summary Create a new country
// @Description Add a new country to the database
// @Tags countries
// @Accept json
// @Produce json
// @Param country body object true "Country object"
// @Success 201 {object} output.SuccessResponse
// @Failure 400 {object} output.ErrorResponse
// @Failure 500 {object} output.ErrorResponse
// @Router /countries [post]
func (c *CountryImpl) Create(ctx echo.Context) error {
	var request struct {
		Country string `json:"country"`
	}

	// ここなんか良い案ないかな。イメージとしては、requestで送られてきた内容がフォーマットに沿ってるかくらいのチェックでいい気がする
	// そして、それはmiddlewareでやりたい
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, output.NewErrorResponse(
			http.StatusBadRequest,
			"invalid request format",
		))
	}

	// 引数どのように渡すかは考えもの
	country, err := entity.NewValidatedCountry(request.Country)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, output.NewErrorResponse(
			http.StatusBadRequest,
			err.Error(),
		))
	}

	if err := c.uc.Insert(country); err != nil {
		return ctx.JSON(http.StatusInternalServerError, output.NewErrorResponse(
			http.StatusInternalServerError,
			"failed to create country",
		))
	}

	return ctx.JSON(http.StatusCreated, output.NewSuccessResponse(country))
}
