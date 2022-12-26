package app

import (
	"github.com/BON4/gofeed/internal/accounts/app/usecase"
)

type Application struct {
	RegisterAccount usecase.RegisterAccountHandler
	LoginAccount    usecase.LoginAccountHandler
}
