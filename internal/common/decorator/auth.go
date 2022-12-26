package decorator

import (
	"context"

	"github.com/BON4/gofeed/internal/common/errors"
)

type AuthManager interface {
	//Checks payload from context. If its not empty, that means that request passed the auth middliware.
	CheckAuth(ctx context.Context) error
}

type commandAuthDecorator[C any] struct {
	base        CommandHandler[C]
	authManager AuthManager
}

func (c commandAuthDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	if err = c.authManager.CheckAuth(ctx); err != nil {
		return errors.NewAuthorizationError("not authenticated", "not-allowed")
	}

	return c.base.Handle(ctx, cmd)
}

type queryAuthDecorator[Q any, R any] struct {
	base        QueryHandler[Q, R]
	authManager AuthManager
}

func (q queryAuthDecorator[Q, R]) Handle(ctx context.Context, cmd Q) (res R, err error) {
	if err = q.authManager.CheckAuth(ctx); err != nil {
		return res, errors.NewAuthorizationError("wrong username or password", "invalid-credentialss")
	}

	return q.base.Handle(ctx, cmd)
}
