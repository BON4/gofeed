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
