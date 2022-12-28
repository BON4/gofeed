package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/BON4/gofeed/internal/accounts/domain"
	"github.com/BON4/gofeed/internal/common/decorator"
	"github.com/BON4/gofeed/internal/common/errors"
	pswrd "github.com/BON4/gofeed/internal/common/password"
	"github.com/sirupsen/logrus"
)

type AccountUsecase struct {
	repo    domain.Repository
	accFc   *domain.Factory
	tokenFc *domain.TokenFactory
	logger  *logrus.Entry
}

func NewAccountUsecase(
	repo domain.Repository,
	fc *domain.Factory,
	tokenFc *domain.TokenFactory,
	logger *logrus.Entry,
) *AccountUsecase {
	return &AccountUsecase{
		repo:    repo,
		tokenFc: tokenFc,
		accFc:   fc,
		logger:  logger,
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

				acc, err := au.accFc.NewAccount(cmd.Username, cmd.Email, cmd.Password, domain.AccountRoleBasic)
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

type LoginQuery struct {
	Username string
	Password string
}

type LoginResponseInstanse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	AccessToken           string                `json:"access_token"`
	AccessTokenExpiresAt  time.Time             `json:"access_token_expires_at"`
	RefreshToken          string                `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time             `json:"refresh_token_expires_at"`
	Instance              LoginResponseInstanse `json:"instance"`
}

type LoginAccountHandler decorator.QueryHandler[LoginQuery, *LoginResponse]

func (au *AccountUsecase) HandleLogin() LoginAccountHandler {
	return decorator.ApplyQueryDecorators[LoginQuery, *LoginResponse](
		decorator.NewQueryHandler(
			func(ctx context.Context, query LoginQuery) (*LoginResponse, error) {

				found, err := au.repo.GetAccount(ctx, query.Username)
				if err != nil {
					if err == sql.ErrNoRows {
						return nil, errors.NewDoesNotExistsError("account with this username does not exists", "account-does-not-exists")
					}
					return nil, err
				}

				if err := pswrd.CheckPassword(query.Password, found.GetPassword()); err != nil {
					return nil, errors.NewInvalidCredError("wrong username or password", "invalid-credentials")
				}

				token, err := au.tokenFc.NewTokenPair(domain.AccountCredentials{
					Username: found.GetUsername(),
					Role:     string(found.GetRole()),
				})

				if err != nil {
					return nil, err
				}

				return &LoginResponse{
					AccessToken:           token.AccessToken,
					AccessTokenExpiresAt:  token.AccessPayload.GetExpiration(),
					RefreshToken:          token.RefreshToken,
					RefreshTokenExpiresAt: token.RefreshPayload.GetExpiration(),
					Instance: LoginResponseInstanse{
						Username: token.RefreshPayload.GetInstance().Username,
						Role:     token.RefreshPayload.GetInstance().Role,
					},
				}, nil
			}),
		au.logger)
}

type SessionCreateQ struct {
}

//type SessionCreateHandler decorator.QueryHandler[]
