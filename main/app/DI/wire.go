//go:build wireinject
// +build wireinject

package di

import (
	"echoProject/main/app/controller"
	con "echoProject/main/domain/controller"
	"echoProject/main/infra/datasource"
	"echoProject/main/infra/repository/user"
	"echoProject/main/infra/things/mysql"
	"echoProject/main/usecase"

	"github.com/google/wire"
)

var infrastructureSet = wire.NewSet(
	mysql.NewSqlHandler,
	datasource.NewUserDataSource,
)

var repositorySet = wire.NewSet(
	user.NewUserRepository,
)

var usecaseSet = wire.NewSet(
	usecase.NewUserService,
)


var controllerSet = wire.NewSet(
	controller.NewUserController,
)

type ControllersSet struct {
	UserController con.User
}

func InitializeController() (*ControllersSet,error) {
	wire.Build(
		infrastructureSet,
		repositorySet,
		usecaseSet,
		controllerSet,
		wire.Struct(new(ControllersSet), "*"),
	)

	return nil, nil
}

