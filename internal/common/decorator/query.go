package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R], logger *logrus.Entry) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base:   handler,
		logger: logger,
	}
}

func ApplyAuthQueryDecorators[H any, R any](handler QueryHandler[H, R], authManager AuthManager, logger *logrus.Entry) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: queryAuthDecorator[H, R]{
			base:        handler,
			authManager: authManager,
		},
		logger: logger,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}

type concreteQueryHandler[Q any, R any] struct {
	handler func(ctx context.Context, query Q) (R, error)
}

func (c *concreteQueryHandler[Q, R]) Handle(ctx context.Context, query Q) (R, error) {
	return c.handler(ctx, query)
}

func NewQueryHandler[Q any, R any](handler func(ctx context.Context, qurey Q) (R, error)) *concreteQueryHandler[Q, R] {
	return &concreteQueryHandler[Q, R]{handler: handler}
}
