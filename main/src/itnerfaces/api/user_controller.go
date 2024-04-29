package api

import (
	"echoProject/main/src/domain"
	"echoProject/main/src/itnerfaces/database"
	"echoProject/main/src/usecase"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	Interactor usecase.UserInteractor
}

// NewUserController routerで参照されるため
// ここで依存性の注入を行なっている？
func NewUserController(sqlHandler database.SqlHandler) *UserController {
	return &UserController{
		Interactor: usecase.UserInteractor{
			UserRepository: &database.UserRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *UserController) Create(c echo.Context) {
	u := domain.User{}
	c.Bind(&u)

	controller.Interactor.Add(u)
	createdUsers := controller.Interactor.GetInfo()
	// response を json 形式に変換
	c.JSON(201, createdUsers)
	return
}

func (controller *UserController) GetUser() []domain.User {
	res := controller.Interactor.GetInfo()
	return res
}

func (controller *UserController) Delete(id string) {
	controller.Interactor.Delete(id)
}
