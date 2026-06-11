package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	dbsql "letsgo/db/sqlc"
	"letsgo/internal/config"
	"letsgo/internal/database"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	dbPool := database.Connect(cfg)
	defer dbPool.Close()

	sqlDB := stdlib.OpenDBFromPool(dbPool)
	defer sqlDB.Close()

	name := getEnv("USER_NAME", "admin")
	email := getEnv("USER_EMAIL", "admin@example.com")
	password := getEnv("USER_PASSWORD", "password")

	if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		log.Fatal("USER_NAME, USER_EMAIL, and USER_PASSWORD must be set")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	ctx := context.Background()
	queries := dbsql.New(sqlDB)

	user, err := queries.GetUserByEmail(ctx, email)
	if err == nil {
		fmt.Printf("user %s already exists with id %d\n", email, user.ID)
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		log.Fatalf("failed to query existing user: %v", err)
	}

	created, err := queries.CreateUser(ctx, dbsql.CreateUserParams{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	})
	if err != nil {
		log.Fatalf("failed to seed user: %v", err)
	}

	fmt.Printf("seeded user %s with id %d\n", email, created.ID)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
