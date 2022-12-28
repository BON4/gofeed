package session

import (
	"time"

	"github.com/google/uuid"
)

type Session[T any] struct {
	Id           uuid.UUID
	Refreshtoken string
	UserAgent    string
	ClientIp     string
	IsBlocked    bool
	ExpiresAt    time.Time
	Instance     T
}
