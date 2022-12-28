package service

import (
	"github.com/BON4/gofeed/internal/accounts/adapters"
	"github.com/BON4/gofeed/internal/accounts/app"
	"github.com/BON4/gofeed/internal/accounts/app/usecase"
	"github.com/BON4/gofeed/internal/accounts/config"
	"github.com/BON4/gofeed/internal/accounts/domain"
	"github.com/sirupsen/logrus"
)

func NewApplication(cfg config.ServerConfig) *app.Application {
	accfc, err := domain.NewFactory(domain.FactoryConfig{
		MinUsernameLen: 4,
		MinPasswordLen: 6,
		DefaultRole:    domain.AccountRoleBasic,
	})
	if err != nil {
		panic(err)
	}

	db, err := adapters.NewPostgresConnection(cfg.DBconn)
	if err != nil {
		panic(err)
	}

	rep := adapters.NewPostgresAccountsRepository(db, accfc)

	tokenfc, err := domain.NewAuthTokenFactory(domain.TokenFactoryConfig{
		AccessTokenDuration:  cfg.AcessDuration,
		RefreshTokenDuration: cfg.RefreshDuration,
		TokenSecret:          cfg.SecretToken,
	})
	if err != nil {
		panic(err)
	}

	logger := logrus.NewEntry(logrus.StandardLogger())

	accUcc := usecase.NewAccountUsecase(rep, accfc, tokenfc, logger)
	return &app.Application{
		LoginAccount:    accUcc.HandleLogin(),
		RegisterAccount: accUcc.HandleRegister(),
	}
}
