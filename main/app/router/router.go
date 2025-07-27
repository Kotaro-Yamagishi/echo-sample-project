package router

import (
	"echoProject/main/app/DI"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func Init() {
	// Echo instancea
	e := echo.New()
	controller,err := di.InitializeController()

	if err != nil {
		fmt.Printf("failed to create event: %s\n", err)
		os.Exit(2)
	}

	e.GET("/users", func(c echo.Context) error {
		users := controller.UserController.GetUser()
		c.Bind(&users)
		return c.JSON(http.StatusOK, users)
	})

	e.POST("/users", func(c echo.Context) error {
		controller.UserController.Create(c)
		return c.String(http.StatusOK, "created")
	})

	e.DELETE("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		controller.UserController.Delete(id)
		return c.String(http.StatusOK, "deleted")
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
