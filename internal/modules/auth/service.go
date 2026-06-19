package auth

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	dbsql "letsgo/db/sqlc"
	sharedjwt "letsgo/shared/jwt"
)

const (
	defaultUserRole      = "user"
	refreshTokenDuration = 7 * 24 * time.Hour
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	repository *Repository
	jwt        *sharedjwt.Provider
}

func NewService(repository *Repository, jwtProvider *sharedjwt.Provider) *Service {
	return &Service{
		repository: repository,
		jwt:        jwtProvider,
	}
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (dbsql.User, string, string, error) {
	user, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil {
		return dbsql.User{}, "", "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return dbsql.User{}, "", "", ErrInvalidCredentials
	}

	return s.issuePair(user)
}

func (s *Service) Refresh(ctx context.Context, tokenString string) (dbsql.User, string, string, error) {
	claims, err := s.jwt.Decode(tokenString)
	if err != nil {
		return dbsql.User{}, "", "", err
	}

	if claims.TokenType != sharedjwt.RefreshTokenType {
		return dbsql.User{}, "", "", ErrInvalidCredentials
	}

	user, err := s.repository.FindUserByID(ctx, claims.UserID)
	if err != nil {
		return dbsql.User{}, "", "", err
	}

	return s.issuePair(user)
}

func (s *Service) FetchUser(ctx context.Context, tokenString string) (dbsql.User, error) {
	claims, err := s.jwt.Decode(tokenString)
	if err != nil {
		return dbsql.User{}, err
	}

	return s.repository.FindUserByID(ctx, claims.UserID)
}

func (s *Service) FindUserByID(ctx context.Context, id int64) (dbsql.User, error) {
	return s.repository.FindUserByID(ctx, id)
}

func (s *Service) issuePair(user dbsql.User) (dbsql.User, string, string, error) {
	accessToken, err := s.jwt.Encode(int64(user.ID), user.Email, defaultUserRole)
	if err != nil {
		return dbsql.User{}, "", "", err
	}

	refreshToken, err := s.jwt.EncodeRefresh(int64(user.ID), user.Email, defaultUserRole, refreshTokenDuration)
	if err != nil {
		return dbsql.User{}, "", "", err
	}

	return user, accessToken, refreshToken, nil
}
