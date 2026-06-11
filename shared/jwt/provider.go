package jwt

import (
	"time"

	gjwt "github.com/golang-jwt/jwt/v5"
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

func (p *Provider) Encode(
	userID int64,
	email string,
	role string,
) (string, error) {

	now := time.Now()

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: gjwt.RegisteredClaims{
			Issuer:   p.issuer,
			IssuedAt: gjwt.NewNumericDate(now),
			ExpiresAt: gjwt.NewNumericDate(
				now.Add(p.expire),
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

func (p *Provider) Validate(
	tokenString string,
) error {

	_, err := p.Decode(tokenString)

	return err
}
