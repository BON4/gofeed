package usecase

import "context"

type PostInfoModel interface {
	List(ctx context.Context, params FindPostParams) ([]*Post, error)
}
