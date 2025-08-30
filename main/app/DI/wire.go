//go:build wireinject
// +build wireinject

package di

import (
	ctrcity "echoProject/app/controller/city"
	country "echoProject/app/controller/country"
	"echoProject/domain/controller"
	dscity "echoProject/infra/datasource/city"
	dscountry "echoProject/infra/datasource/country"
	repocity "echoProject/infra/repository/city"
	repocountry "echoProject/infra/repository/country"
	"echoProject/infra/things/sqlboiler"
	uccity "echoProject/usecase/city"
	uccountry "echoProject/usecase/country"

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
