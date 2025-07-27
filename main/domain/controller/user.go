package controller

import (
	"echoProject/main/domain/entity"	
"github.com/labstack/echo/v4"
)


type User interface {
	Create(c echo.Context)
	GetUser() []entity.User
	Delete(id string)
}