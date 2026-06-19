// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"log"

	"letsgo/internal/app"
	"letsgo/internal/routes"
)

func main() {
	// Initialize the application
	application := app.New()

	// Register routes
	routes.Register(application)

	// Run the application
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
