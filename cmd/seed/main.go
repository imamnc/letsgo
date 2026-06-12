package main

import (
	"context"
	"log"

	"letsgo/db/seeds"
	"letsgo/internal/app"
)

func main() {
	// Initialize the app singleton to get access to the database connection and queries
	application := app.New()
	// Get the SQLC queries from the application instance
	queries := application.Queries

	// Seed the user data
	if err := seeds.SeedUser(context.Background(), queries); err != nil {
		log.Fatal(err)
	}
}
