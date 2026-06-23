// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"

	"letsgo/internal/app"
	"letsgo/internal/routes"
	"letsgo/internal/scheduler"
)

func main() {
	// Initialize the application
	application := app.New()

	// Register routes
	routes.Register(application)

	// Start the scheduler
	scheduler.Start(application)

	// Run the application
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
