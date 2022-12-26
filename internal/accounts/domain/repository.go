package domain

import "context"

type Repository interface {
	CreateAccount(ctx context.Context, acc *Account) (*Account, error)
	GetAccount(ctx context.Context, username string) (*Account, error)
}
