package controller

import (
	"echoProject/main/domain/controller"
	"echoProject/main/domain/entity"
	"echoProject/main/domain/service"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	service service.UserService
}

// 依存関係

func NewUserController(userService service.UserService) controller.User {
	return &UserController{service: userService}
}

func (ControllerUserController *UserController) Create(c echo.Context) {
	u := entity.User{}
	c.Bind(&u)

	ControllerUserController.service.Add(u)
	createdUsers := ControllerUserController.service.GetInfo()
	// response を json 形式に変換
	c.JSON(201, createdUsers)
}

func (ControllerUserController *UserController) GetUser() []entity.User {
	res := ControllerUserController.service.GetInfo()
	return res
}

func (ControllerUserController *UserController) Delete(id string) {
	ControllerUserController.service.Delete(id)
}
