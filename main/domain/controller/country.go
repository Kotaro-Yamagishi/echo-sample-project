package controller

import (
	"echoProject/main/domain/entity"

	"github.com/labstack/echo/v4"
)

type Country interface {
	GetCountry() []entity.Country
	Create(ctx echo.Context) error
}
