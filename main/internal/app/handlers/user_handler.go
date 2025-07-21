package handlers

import (
	"echoProject/main/internal/app/repositories"
	"echoProject/main/internal/app/services"
	domain "echoProject/main/internal/models"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService services.UserServiceIF
}

func NewUserHandler(sqlHandler repositories.SqlHandler) *UserHandler {
	return &UserHandler{
		UserService: services.NewUserService(
			repositories.NewUserRepository(sqlHandler),
		),
	}
}

func (handler *UserHandler) Create(c echo.Context) {
	u := domain.User{}
	c.Bind(&u)

	handler.UserService.Add(u)
	createdUsers := handler.UserService.GetInfo()
	// response を json 形式に変換
	c.JSON(201, createdUsers)
}

func (handler *UserHandler) GetUser() []domain.User {
	res := handler.UserService.GetInfo()
	return res
}

func (handler *UserHandler) Delete(id string) {
	handler.UserService.Delete(id)
}
