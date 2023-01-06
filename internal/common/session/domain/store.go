package domain

import "context"

type Store interface {
	Get(context.Context, string) (*Session, error)
	Set(context.Context, string, *Session) error
}
