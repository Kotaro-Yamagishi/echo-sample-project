// @title Echo Sample Project API
// @version 1.0
// @description This is a sample Echo server with Sakila database.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1324
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"echoProject/app/router"
)

// @Summary Health check
// @Description Get server health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func main() {
	// Initialize router
	router.Init()
}
