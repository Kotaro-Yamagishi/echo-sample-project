package router

import (
	di "echoProject/main/app/DI"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
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

	// Start server
	e.Logger.Fatal(e.Start(":1324"))
}
