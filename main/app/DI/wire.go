//go:build wireinject
// +build wireinject

package di

import (
	ctrcity "echoProject/main/app/controller/city"
	country "echoProject/main/app/controller/country"
	"echoProject/main/domain/controller"
	dscity "echoProject/main/infra/datasource/city"
	dscountry "echoProject/main/infra/datasource/country"
	repocity "echoProject/main/infra/repository/city"
	repocountry "echoProject/main/infra/repository/country"
	"echoProject/main/infra/things/sqlboiler"
	uccity "echoProject/main/usecase/city"
	uccountry "echoProject/main/usecase/country"

	"github.com/google/wire"
)

var infrastructureSet = wire.NewSet(
	dscity.NewCityDataSource,
	dscountry.NewCountryDataSource,
	sqlboiler.NewSQLBoilerImpl,
)

var repositorySet = wire.NewSet(
	repocity.NewCityRepository,
	repocountry.NewCountryRepository,
)

var usecaseSet = wire.NewSet(
	uccity.NewCityService,
	uccountry.NewCountryService,
)

var controllerSet = wire.NewSet(
	ctrcity.NewCityController,
	country.NewCountryController,
)

type ControllersSet struct {
	CityController    controller.City
	CountryController controller.Country
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
