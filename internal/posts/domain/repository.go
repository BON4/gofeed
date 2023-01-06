package domain

import "context"

type Repository interface {
	Create(ctx context.Context, post *Post) (int64, error)
	Delete(ctx context.Context, postId int64) error
}
