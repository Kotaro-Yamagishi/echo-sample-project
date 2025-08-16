//go:build wireinject
// +build wireinject

package di

import (
	ctrcity "echoProject/main/app/controller/city"
	"echoProject/main/domain/controller"
	dscity "echoProject/main/infra/datasource/city"
	repocity "echoProject/main/infra/repository/city"
	"echoProject/main/infra/things/sqlboiler"
	uccity "echoProject/main/usecase/city"

	"github.com/google/wire"
)

var infrastructureSet = wire.NewSet(
	dscity.NewCityDataSource,
	sqlboiler.NewSQLBoilerImpl,
)

var repositorySet = wire.NewSet(
	repocity.NewCityRepository,
)

var usecaseSet = wire.NewSet(
	uccity.NewCityService,
)

var controllerSet = wire.NewSet(
	ctrcity.NewCityController,
)

type ControllersSet struct {
	CityController controller.City
}

type initializeDBSet struct {
	SqlBoiler sqlboiler.SQLBoiler
}

func InitializeController() (*ControllersSet, error) {
	wire.Build(
		infrastructureSet,
		repositorySet,
		usecaseSet,
		controllerSet,
		wire.Struct(new(ControllersSet), "*"),
	)

	return nil, nil
}

func InitializeDB() (*initializeDBSet, error) {
	wire.Build(
		infrastructureSet,
		wire.Struct(new(initializeDBSet), "*"),
	)

	return nil, nil
}
