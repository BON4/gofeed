package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/BON4/gofeed/internal/accounts/domain"
	"github.com/BON4/gofeed/internal/common/decorator"
	"github.com/BON4/gofeed/internal/common/errors"
	pswrd "github.com/BON4/gofeed/internal/common/password"
	sessDomain "github.com/BON4/gofeed/internal/common/session/domain"
	"github.com/BON4/gofeed/internal/common/tokens"
	"github.com/sirupsen/logrus"
)

type AuthUsecase struct {
	repo    domain.Repository
	accFc   *domain.Factory
	tokenFc *sessDomain.TokenFactory
	logger  *logrus.Entry
}

func NewAccountUsecase(
	repo domain.Repository,
	fc *domain.Factory,
	tokenFc *sessDomain.TokenFactory,
	logger *logrus.Entry,
) *AuthUsecase {
	return &AuthUsecase{
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

func (au *AuthUsecase) HandleRegister() RegisterAccountHandler {
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
	Username string
	Role     string
}

type LoginResponse struct {
	AccessTokenId         string
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshTokenId        string
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	Instance              LoginResponseInstanse
}

type LoginAccountHandler decorator.QueryHandler[LoginQuery, *LoginResponse]

func (au *AuthUsecase) HandleLogin() LoginAccountHandler {
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

				token, err := au.tokenFc.NewTokenPair(tokens.InstanceCredentials{
					Username: found.GetUsername(),
					Role:     string(found.GetRole()),
				})

				if err != nil {
					return nil, err
				}

				return &LoginResponse{
					AccessTokenId:         token.AccessPayload.Id.String(),
					AccessToken:           token.AccessToken,
					AccessTokenExpiresAt:  token.AccessPayload.ExpiresAt,
					RefreshTokenId:        token.RefreshPayload.Id.String(),
					RefreshToken:          token.RefreshToken,
					RefreshTokenExpiresAt: token.RefreshPayload.ExpiresAt,
					Instance: LoginResponseInstanse{
						Username: token.RefreshPayload.Instance.Username,
						Role:     token.RefreshPayload.Instance.Role,
					},
				}, nil
			}),
		au.logger)
}
