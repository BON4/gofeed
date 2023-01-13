package service

import (
	"os"
	"strconv"

	"github.com/BON4/gofeed/internal/accounts/adapters"
	"github.com/BON4/gofeed/internal/accounts/app"
	"github.com/BON4/gofeed/internal/accounts/app/usecase"
	"github.com/BON4/gofeed/internal/accounts/config"
	"github.com/BON4/gofeed/internal/accounts/domain"
	"github.com/BON4/gofeed/internal/common/session"
	sessAdapters "github.com/BON4/gofeed/internal/common/session/adapters"
	sessDomain "github.com/BON4/gofeed/internal/common/session/domain"
	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewApplication(cfg config.ServerConfig) (*app.Application, func()) {
	logger := logrus.NewEntry(logrus.StandardLogger())

	accfc, err := domain.NewFactory(domain.FactoryConfig{
		MinUsernameLen: 4,
		MinPasswordLen: 6,
		DefaultRole:    domain.AccountRoleBasic,
	})
	if err != nil {
		panic(err)
	}

	db, err := adapters.NewPostgresConnection(os.Getenv("DB_SOURCE"))
	if err != nil {
		panic(err)
	}

	// logger.Info("Applying migrations")

	// if err := runDBMigration(cfg.MigrationsPath, cfg.DBconn); err != nil {
	// 	panic(err)
	// }

	rep := adapters.NewPostgresAccountsRepository(db, accfc)

	tokenfc, err := sessDomain.NewAuthTokenFactory(sessDomain.TokenFactoryConfig{
		AccessTokenDuration:  cfg.AcessDuration,
		RefreshTokenDuration: cfg.RefreshDuration,
		TokenSecret:          cfg.SecretToken,
	})
	if err != nil {
		panic(err)
	}

	accUcc := usecase.NewAccountUsecase(rep, accfc, tokenfc, logger)

	sessFc, err := sessDomain.NewSessionFactory(nil)

	if err != nil {
		panic(err)
	}

	redis_db, _ := strconv.Atoi(os.Getenv(""))

	redisCli, err := sessAdapters.NewRedisConnection(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"), redis_db)
	if err != nil {
		panic(err)
	}

	redisStore := sessAdapters.NewRedisStore(redisCli, &sessFc)

	sessUc := session.NewSessionUsecase(redisStore, &sessFc, logger)

	return &app.Application{
			LoginAccount:    accUcc.HandleLogin(),
			RegisterAccount: accUcc.HandleRegister(),
			CreateSession:   sessUc.HandleCreateSession(),
			SessionIsValid:  sessUc.HadleIsValidSession(),
		},
		func() {
			_ = redisCli.Close()
			_ = db.Close()
		}
}

func runDBMigration(migrationURL string, dbSource string) error {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		return err
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
