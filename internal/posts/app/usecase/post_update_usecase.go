package usecase

import (
	"context"

	"github.com/BON4/gofeed/internal/common/decorator"
	"github.com/sirupsen/logrus"
)

type PostUpdateModel interface {
	RatePost(ctx context.Context, params PostRateParams) error
}

type PostUpdateUsecase struct {
	repo   PostUpdateModel
	logger *logrus.Entry
}

func NewPostUpdateUsecase(
	repo PostUpdateModel,
	logger *logrus.Entry,
) *PostUpdateUsecase {
	return &PostUpdateUsecase{
		repo:   repo,
		logger: logger,
	}
}

type RatePostHandler decorator.CommandHandler[PostRateParams]

func (p *PostUpdateUsecase) HandleRatePost() RatePostHandler {
	return decorator.ApplyCommandDecorators[PostRateParams](decorator.NewCommandHandler(func(ctx context.Context, cmd PostRateParams) error {
		return p.repo.RatePost(ctx, cmd)
	}), p.logger)
}
