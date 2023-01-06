package app

import (
	"github.com/BON4/gofeed/internal/accounts/app/usecase"
	"github.com/BON4/gofeed/internal/common/session"
)

type Application struct {
	RegisterAccount usecase.RegisterAccountHandler
	LoginAccount    usecase.LoginAccountHandler
	CreateSession   session.CreateSessionHandler
	SessionIsValid  session.IsValidSessionHandler
}
