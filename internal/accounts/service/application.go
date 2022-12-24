package service

import (
	"github.com/BON4/gofeed/internal/accounts/adapters"
	"github.com/BON4/gofeed/internal/accounts/app"
)

func NewApplication() *app.Application {
	rep := adapters.NewPostgresAccountsRepository()
	return &app.Application{
		Repo: rep,
	}
}
