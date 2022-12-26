package usecase

import (
	"context"
	"database/sql"

	"github.com/BON4/gofeed/internal/accounts/domain"
	"github.com/BON4/gofeed/internal/common/decorator"
	"github.com/BON4/gofeed/internal/common/errors"
	pswrd "github.com/BON4/gofeed/internal/common/password"
	"github.com/sirupsen/logrus"
)

type AccountUsecase struct {
	repo   domain.Repository
	fc     *domain.Factory
	logger *logrus.Entry
}

func NewAccountUsecase(
	repo domain.Repository,
	fc *domain.Factory,
	logger *logrus.Entry,
) *AccountUsecase {
	return &AccountUsecase{
		repo:   repo,
		logger: logger,
	}
}

type RegisterCommand struct {
	Username string
	Email    string
	Password string
}

type RegisterAccountHandler decorator.CommandHandler[RegisterCommand]

func (au *AccountUsecase) HandleRegister() RegisterAccountHandler {
	return decorator.ApplyCommandDecorators[RegisterCommand](
		decorator.NewCommandHandler(
			func(ctx context.Context, cmd RegisterCommand) error {

				_, err := au.repo.GetAccount(ctx, cmd.Username)
				if err == nil {
					return errors.NewAlreadyExistsError("account with this username already exists", "register-already-exists")
				}

				if err != nil {
					if err != sql.ErrNoRows {
						return err
					}
				}

				acc, err := au.fc.NewAccount(cmd.Username, cmd.Password, cmd.Email, domain.AccountRoleBasic)
				if err != nil {
					return err
				}

				_, err = au.repo.CreateAccount(ctx, acc)
				if err != nil {
					return err
				}

				return nil
			}),
		au.logger)
}

type LoginCommand struct {
	Username string
	Password string
}

type LoginAccountHandler decorator.CommandHandler[LoginCommand]

func (au *AccountUsecase) HandleLogin() LoginAccountHandler {
	return decorator.ApplyCommandDecorators[LoginCommand](
		decorator.NewCommandHandler(
			func(ctx context.Context, cmd LoginCommand) error {

				found, err := au.repo.GetAccount(ctx, cmd.Username)
				if err != nil {
					if err == sql.ErrNoRows {
						return errors.NewDoesNotExistsError("account with this username does not exists", "account-does-not-exists")
					}
					return err
				}

				if err := pswrd.CheckPassword(cmd.Password, found.GetPassword()); err != nil {
					return errors.NewInvalidCredError("wrong username or password", "invalid-credentials")
				}

				return nil
			}),
		au.logger)
}
