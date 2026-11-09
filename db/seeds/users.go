package seeds

import (
	"context"
	"fmt"
	"strings"

	dbsql "letsgo/db/sqlc"

	"golang.org/x/crypto/bcrypt"
)

type SeedUserParams struct {
	Name     string
	Gender   string
	Email    string
	Password string
}

func PrepareUsers() []SeedUserParams {
	password, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Sprintf("failed to hash password: %v", err))
	}

	return []SeedUserParams{
		{
			Name:     "John Doe",
			Gender:   "Male",
			Email:    "johndoe@mail.com",
			Password: string(password),
		},
		{
			Name:     "Lusy Ariana",
			Gender:   "Female",
			Email:    "lusy.ariana@mail.com",
			Password: string(password),
		},
	}
}

func SeedUser(ctx context.Context, queries *dbsql.Queries) error {
	users := PrepareUsers()

	for _, user := range users {
		_, err := queries.CreateUser(ctx, dbsql.CreateUserParams{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			Gender:   user.Gender,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				fmt.Printf("User with email %s already exists, skipping...\n", user.Email)
				continue
			}
			return fmt.Errorf("failed to create user %s: %w", user.Email, err)
		}
		fmt.Printf("User with email %s created successfully\n", user.Email)
	}

	return nil
}
