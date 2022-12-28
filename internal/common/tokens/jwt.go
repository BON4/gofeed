package tokens

import (
	"errors"
	"fmt"
	"time"

	cerrors "github.com/BON4/gofeed/internal/common/errors"
	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTGenerator[T any] struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTGenerator[T any](secretKey string) (Generator[T], error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTGenerator[T]{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTGenerator[T]) CreateToken(instance T, duration time.Duration) (string, *Payload[T], error) {
	payload, err := NewPayload(instance, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *JWTGenerator[T]) VerifyToken(token string) (*Payload[T], error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload[T]{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errExpiredToken) {
			return nil, cerrors.NewAuthorizationError("no token has been provided", "no-token")
		}
		return nil, cerrors.NewAuthorizationError("invalid token has been provided", "invalid-token")
	}

	payload, ok := jwtToken.Claims.(*Payload[T])
	if !ok {
		return nil, cerrors.NewAuthorizationError("invalid token has been provided", "invalid-token")
	}

	return payload, nil
}
