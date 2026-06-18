package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	dbsql "letsgo/db/sqlc"
	"letsgo/shared/json"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func cacheKeyForUserID(id int32) string {
	return fmt.Sprintf("user:%d", id)
}

func marshalUser(user dbsql.User) (string, error) {
	return json.Marshal(user)
}

func unmarshalUser(data string) (dbsql.User, error) {
	return json.Unmarshal[dbsql.User](data)
}

func (s *Service) CreateUser(ctx context.Context, request CreateUserRequest) (dbsql.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(request.Password)), bcrypt.DefaultCost)
	if err != nil {
		return dbsql.User{}, err
	}

	return s.repository.CreateUser(ctx, dbsql.CreateUserParams{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
	})
}

func (s *Service) FindUserByID(ctx context.Context, id int32) (dbsql.User, error) {
	cacheKey := cacheKeyForUserID(id)
	cached, err := s.repository.app.Redis.Remember(ctx, cacheKey, 5*time.Minute, func() (string, error) {
		user, err := s.repository.FindUserByID(ctx, id)
		if err != nil {
			return "", err
		}
		return marshalUser(user)
	})
	if err != nil {
		return dbsql.User{}, err
	}

	return unmarshalUser(cached)
}

func (s *Service) ListUsers(ctx context.Context, limit, offset int32) ([]dbsql.User, error) {
	return s.repository.ListUsers(ctx, limit, offset)
}

func (s *Service) UpdateUser(ctx context.Context, id int32, request UpdateUserRequest) (dbsql.User, error) {
	existing, err := s.repository.FindUserByID(ctx, id)
	if err != nil {
		return dbsql.User{}, err
	}

	name := existing.Name
	if request.Name != nil {
		trimmed := strings.TrimSpace(*request.Name)
		if trimmed == "" {
			return dbsql.User{}, errors.New("name cannot be empty")
		}
		name = trimmed
	}

	email := existing.Email
	if request.Email != nil {
		trimmed := strings.TrimSpace(*request.Email)
		if trimmed == "" {
			return dbsql.User{}, errors.New("email cannot be empty")
		}
		email = trimmed
	}

	password := existing.Password
	if request.Password != nil {
		trimmed := strings.TrimSpace(*request.Password)
		if trimmed == "" {
			return dbsql.User{}, errors.New("password cannot be empty")
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(trimmed), bcrypt.DefaultCost)
		if err != nil {
			return dbsql.User{}, err
		}
		password = string(hashedPassword)
	}

	user, err := s.repository.UpdateUser(ctx, dbsql.UpdateUserParams{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return dbsql.User{}, err
	}

	return user, s.repository.app.Redis.Forget(ctx, cacheKeyForUserID(id))
}

func (s *Service) DeleteUser(ctx context.Context, id int32) error {
	err := s.repository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return s.repository.app.Redis.Forget(ctx, cacheKeyForUserID(id))
}
