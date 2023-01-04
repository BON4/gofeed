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

// Payload contains the payload data of the token
type Payload[T InstanceCredentials] struct {
	id        uuid.UUID
	issuedAt  time.Time
	expiresAt time.Time
	instance  T
}

func (p *Payload[T]) GetExpiration() time.Time {
	return p.expiresAt
}

func (p *Payload[T]) GetIssued() time.Time {
	return p.issuedAt
}

func (p *Payload[T]) GetId() uuid.UUID {
	return p.id
}

func (p *Payload[T]) GetInstance() T {
	return p.instance
}

func NewPayload[T InstanceCredentials](instance T, duration time.Duration) (*Payload[T], error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload[T]{
		id:        tokenID,
		instance:  instance,
		issuedAt:  time.Now(),
		expiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload[T]) Valid() error {
	if time.Now().After(payload.expiresAt) {
		return errExpiredToken
	}
	return nil
}

func GetPayloadFromContext[T InstanceCredentials](ctx *gin.Context, payloadkey string) (*Payload[T], error) {
	payload, ok := ctx.Get(payloadkey)
	if !ok {
		return nil, cerrors.NewAuthorizationError("no token has been provided", "no-token")
	}

	payloadParsed, ok := payload.(*Payload[T])
	if !ok {
		return nil, cerrors.NewAuthorizationError("invalid token has been provided", "invalid-token")

	}
	return payloadParsed, nil
}
