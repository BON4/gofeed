package tokens

import (
	"errors"
	"time"

	cerrors "github.com/BON4/gofeed/internal/common/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	errInvalidToken = errors.New("token is invalid")
	errExpiredToken = errors.New("token has expired")
)

// TODO: maby make payload fields private and add MarshalJSON()

// Payload contains the payload data of the token
type Payload[T InstanceCredentials] struct {
	Id        uuid.UUID `json:"id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Instance  T         `json:"instance"`
}

func NewPayload[T InstanceCredentials](instance T, duration time.Duration) (*Payload[T], error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload[T]{
		Id:        tokenID,
		Instance:  instance,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload[T]) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return errExpiredToken
	}
	return nil
}

type ctxKey int

const (
	payloadContextKey ctxKey = iota
)

func GetPayloadFromContext[T InstanceCredentials](ctx *gin.Context, payloadkey string) (*Payload[T], error) {
	// TODO: change payloadkey to payloadContextKey
	payload, ok := ctx.Get(payloadkey)
	if !ok {
		return nil, cerrors.NewAuthorizationError("no token has been provided", "get-no-token")
	}

	payloadParsed, ok := payload.(*Payload[T])
	if !ok {
		return nil, cerrors.NewAuthorizationError("invalid token has been provided", "type-invalid-token")

	}
	return payloadParsed, nil
}
