package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func ApplyCommandDecorators[H any](handler CommandHandler[H], logger *logrus.Entry) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base:   handler,
		logger: logger,
	}
}

func ApplyAuthCommandDecorators[H any](handler CommandHandler[H], authManager AuthManager, logger *logrus.Entry) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base: commandAuthDecorator[H]{
			base:        handler,
			authManager: authManager,
		},
		logger: logger,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

type concreteCommandHandler[C any] struct {
	handler func(ctx context.Context, cmd C) error
}

func (c *concreteCommandHandler[C]) Handle(ctx context.Context, cmd C) error {
	return c.handler(ctx, cmd)
}

func NewCommandHandler[C any](handler func(ctx context.Context, cmd C) error) *concreteCommandHandler[C] {
	return &concreteCommandHandler[C]{handler: handler}
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
