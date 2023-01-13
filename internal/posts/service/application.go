package service

import (
	"os"

	"github.com/BON4/gofeed/internal/common/session"
	sessDomain "github.com/BON4/gofeed/internal/common/session/domain"
	"github.com/BON4/gofeed/internal/posts/adapters"
	"github.com/BON4/gofeed/internal/posts/app"
	"github.com/BON4/gofeed/internal/posts/app/usecase"
	"github.com/BON4/gofeed/internal/posts/config"
	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewSessionMiddleware(cfg config.ServerConfig) (*session.SessionMiddleware, func()) {
	logger := logrus.NewEntry(logrus.StandardLogger())

	tokenVer, err := sessDomain.NewTokenVerifier(cfg.SecretToken)
	if err != nil {
		panic(err)
	}

	sessMdwr := session.NewSessionMiddleware(tokenVer, logger, cfg.HeaderKey)

	return sessMdwr, func() {
	}
}

func NewApplication(cfg config.ServerConfig) (*app.Application, func()) {
	logger := logrus.NewEntry(logrus.StandardLogger())

	db, err := adapters.NewPostgresConnection(os.Getenv("DB_SOURCE"))
	if err != nil {
		panic(err)
	}

	// logger.Info("Applying migrations")

	// if err := runDBMigration(cfg.MigrationsPath, cfg.DBconn); err != nil {
	// 	panic(err)
	// }

	repo := adapters.NewPostgresPostsRepository(db)

	postUc := usecase.NewPostsUsecase(repo, logger)
	postInfoUc := usecase.NewPostInfoUsecase(repo, logger)
	postUpdUc := usecase.NewPostUpdateUsecase(repo, logger)

	return &app.Application{
			CreatePost: postUc.HandleCreatePost(),
			DeletePost: postUc.HadleDeletePost(),
			ListPost:   postInfoUc.HandleListPosts(),
			RatePost:   postUpdUc.HandleRatePost(),
		},
		func() {
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
