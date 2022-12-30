package domain

import "context"

type Repository interface {
	CreateAccount(ctx context.Context, acc *Account) (*Account, error)
	GetAccount(ctx context.Context, username string) (*Account, error)
}

type Store interface {
	Get(context.Context, string) (*Session, error)
	Set(context.Context, string, *Session) error
}
