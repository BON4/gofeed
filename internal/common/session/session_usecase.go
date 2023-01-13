package session

import (
	"context"
	"time"

	"github.com/BON4/gofeed/internal/common/decorator"
	"github.com/BON4/gofeed/internal/common/session/domain"
	"github.com/BON4/gofeed/internal/common/tokens"
	"github.com/sirupsen/logrus"
)

// TODO: Hanlde refresh access token
// TODO: Make separate microsrevice for sessions

type SessionUsecase struct {
	fc     *domain.SessionFactory
	store  domain.Store
	logger *logrus.Entry
}

func NewSessionUsecase(
	store domain.Store,
	fc *domain.SessionFactory,
	logger *logrus.Entry,
) *SessionUsecase {
	return &SessionUsecase{
		store:  store,
		fc:     fc,
		logger: logger,
	}
}

type CreateSessionCommand struct {
	ID           string
	Refreshtoken string
	UserAgent    string
	ClientIp     string
	ExpiresAt    time.Time
	Instance     tokens.InstanceCredentials
}

type CreateSessionHandler decorator.CommandHandler[CreateSessionCommand]

func (su *SessionUsecase) HandleCreateSession() CreateSessionHandler {
	return decorator.ApplyCommandDecorators[CreateSessionCommand](
		decorator.NewCommandHandler(func(ctx context.Context, cmd CreateSessionCommand) error {
			ss, err := su.fc.NewSession(
				cmd.ID,
				cmd.Refreshtoken,
				cmd.UserAgent,
				cmd.ClientIp,
				time.Now(),
				cmd.ExpiresAt,
				cmd.Instance)

			if err != nil {
				return err
			}

			return su.store.Set(ctx, ss.Id, ss)
		}), su.logger)
}

type IsValidSessionQuery struct {
	ID string
}

type IsValidSessionHandler decorator.QueryHandler[IsValidSessionQuery, bool]

func (su *SessionUsecase) HadleIsValidSession() IsValidSessionHandler {
	return decorator.ApplyQueryDecorators[IsValidSessionQuery, bool](
		decorator.NewQueryHandler(func(ctx context.Context, qurey IsValidSessionQuery) (bool, error) {
			ss, err := su.store.Get(ctx, qurey.ID)
			if err != nil {
				return false, err
			}

			return !ss.IsBlocked(), nil
		}), su.logger)
}
