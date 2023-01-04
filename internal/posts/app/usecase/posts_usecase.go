package usecase

import (
	"context"

	"github.com/BON4/gofeed/internal/common/decorator"
	"github.com/BON4/gofeed/internal/posts/domain"

	"github.com/sirupsen/logrus"
)

type PostsUsecase struct {
	repo   domain.Repository
	logger *logrus.Entry
}

func NewPostsUsecase(
	repo domain.Repository,
	logger *logrus.Entry,
) *PostsUsecase {
	return &PostsUsecase{
		repo:   repo,
		logger: logger,
	}
}

type CreatePostQuery struct {
	Content string
	Account string
}

type CreatePostHandler decorator.QueryHandler[CreatePostQuery, int64]

func (p *PostsUsecase) HandleCreatePost() CreatePostHandler {
	return decorator.ApplyQueryDecorators[CreatePostQuery, int64](decorator.NewQueryHandler(func(ctx context.Context, qurey CreatePostQuery) (int64, error) {
		post, err := domain.NewPost(qurey.Account, qurey.Content)
		if err != nil {
			return -1, err
		}

		created, err := p.repo.Create(ctx, post)
		if err != nil {
			return -1, err
		}

		return created.ID(), nil
	}), p.logger)
}

type DeletePostCommand struct {
	PostId int64
}

type DeletePostHandler decorator.CommandHandler[DeletePostCommand]

func (p *PostsUsecase) HadleDeletePost() DeletePostHandler {
	return decorator.ApplyCommandDecorators[DeletePostCommand](decorator.NewCommandHandler(func(ctx context.Context, cmd DeletePostCommand) error {
		return p.repo.Delete(ctx, cmd.PostId)
	}), p.logger)
}

type ListPostHandler decorator.QueryHandler[domain.FindPostParams]
