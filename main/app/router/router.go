package router

import (
	di "echoProject/app/DI"
	"fmt"
	"net/http"
	"os"

	_ "echoProject/docs" // This is required for swagger

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init() {
	// Echo instancea
	e := echo.New()
	ctr, err := di.InitializeController()

	if err != nil {
		fmt.Printf("failed to create event: %s\n", err)
		os.Exit(2)
	}

	e.GET("/cities", func(c echo.Context) error {
		cities := ctr.CityController.GetCity()
		c.Bind(&cities)
		return c.JSON(http.StatusOK, cities)
	})

	e.GET("/countries", func(c echo.Context) error {
		countries := ctr.CountryController.GetCountry(c)
		c.Bind(&countries)
		return c.JSON(http.StatusOK, countries)
	})

	e.POST("/countries", func(c echo.Context) error {
		return ctr.CountryController.Create(c)
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	e.Logger.Fatal(e.Start(":1324"))
}
