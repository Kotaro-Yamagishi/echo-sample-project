package controller

import (
	"github.com/labstack/echo/v4"
)

type Country interface {
	GetCountry(ctx echo.Context) error
	Create(ctx echo.Context) error
}
