package jwt

import (
	"time"

	gjwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

type Provider struct {
	secret []byte
	issuer string
	expire time.Duration
}

func New(
	secret string,
	issuer string,
	expire time.Duration,
) *Provider {
	return &Provider{
		secret: []byte(secret),
		issuer: issuer,
		expire: expire,
	}
}

// Encode creates a JWT access token with the given user information and default expiration time.
func (p *Provider) Encode(
	userID int64,
	email string,
	role string,
) (string, error) {
	return p.encodeWithType(userID, email, role, p.expire, AccessTokenType)
}

// EncodeRefresh creates a JWT refresh token with the given user information and custom expiration time.
func (p *Provider) EncodeRefresh(
	userID int64,
	email string,
	role string,
	expire time.Duration,
) (string, error) {
	return p.encodeWithType(userID, email, role, expire, RefreshTokenType)
}

func (p *Provider) encodeWithType(
	userID int64,
	email string,
	role string,
	expire time.Duration,
	tokenType string,
) (string, error) {

	now := time.Now()

	claims := Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: gjwt.RegisteredClaims{
			ID:       uuid.NewString(),
			Issuer:   p.issuer,
			IssuedAt: gjwt.NewNumericDate(now),
			ExpiresAt: gjwt.NewNumericDate(
				now.Add(expire),
			),
		},
	}

	token := gjwt.NewWithClaims(
		gjwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		p.secret,
	)
}

// Decode parses and validates the given JWT token string and returns the claims if valid.
func (p *Provider) Decode(
	tokenString string,
) (*Claims, error) {

	token, err := gjwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *gjwt.Token) (interface{}, error) {
			return p.secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// Validate checks if the given JWT token string is valid and not expired.
func (p *Provider) Validate(
	tokenString string,
) error {

	_, err := p.Decode(tokenString)

	return err
}
