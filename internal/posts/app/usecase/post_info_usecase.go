package usecase

import (
	"context"

	"github.com/BON4/gofeed/internal/common/decorator"
	"github.com/BON4/gofeed/internal/common/errors"
	"github.com/sirupsen/logrus"
)

type PostInfoModel interface {
	// TODO: it is dumb, that List method uses usecase.Post, better be domain.Post
	List(ctx context.Context, params FindPostParams) ([]*Post, error)
}

type PostInfoUsecase struct {
	repo   PostInfoModel
	logger *logrus.Entry
}

func NewPostInfoUsecase(
	repo PostInfoModel,
	logger *logrus.Entry,
) *PostInfoUsecase {
	return &PostInfoUsecase{
		repo:   repo,
		logger: logger,
	}
}

type ListPostHandler decorator.QueryHandler[FindPostParams, []*Post]

func (p *PostInfoUsecase) HandleListPosts() ListPostHandler {
	return decorator.ApplyQueryDecorators[FindPostParams, []*Post](decorator.NewQueryHandler(func(ctx context.Context, qurey FindPostParams) ([]*Post, error) {
		if qurey.PageNumber < 0 || qurey.PageSize <= 0 {
			return nil, errors.NewSlugError("page pagination params are not valid", "not-valid-pagination")
		}

		posts, err := p.repo.List(ctx, qurey)
		if err != nil {
			return nil, err
		}

		return posts, nil
	}), p.logger)
}
