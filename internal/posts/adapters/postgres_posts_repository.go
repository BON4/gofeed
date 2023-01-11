package adapters

import (
	"context"
	"database/sql"

	"github.com/BON4/gofeed/internal/posts/adapters/sqlc"
	"github.com/BON4/gofeed/internal/posts/app/usecase"
	"github.com/BON4/gofeed/internal/posts/domain"
	_ "github.com/lib/pq"
)

type PostgresPostsRepository struct {
	querys sqlc.Store
}

func NewPostgresPostsRepository(dbcon *sql.DB) *PostgresPostsRepository {
	return &PostgresPostsRepository{
		querys: sqlc.NewStore(dbcon),
	}
}

func (p *PostgresPostsRepository) RatePost(ctx context.Context, params usecase.PostRateParams) error {
	var newRate int32 = 0
	rate, err := p.querys.GetRatePost(ctx, sqlc.GetRatePostParams{
		PostID:  params.PostId,
		Account: params.Account,
	})

	if err == nil {
		newRate = rate.RatedScore
	} else {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if params.Rate == rate.RatedScore {
		return nil
	}

	// TODO: begining transaction to like/dislike post is too expansive
	return p.querys.ExecTx(ctx, func(q *sqlc.Queries) error {
		var err error

		err = q.RatePost(ctx, sqlc.RatePostParams{
			PostID:     params.PostId,
			Account:    params.Account,
			RatedScore: params.Rate,
		})

		if err != nil {
			return err
		}

		_, err = q.AddPostScore(ctx, sqlc.AddPostScoreParams{
			PostID: params.PostId,
			Score:  params.Rate - newRate,
		})

		return err
	})
}

func (p *PostgresPostsRepository) Create(ctx context.Context, post *domain.Post) (int64, error) {
	createdId, err := p.querys.CreatePost(ctx, sqlc.CreatePostParams{
		Content:  post.Content(),
		PostedOn: post.PostedOn(),
		PostedBy: post.PostedBy(),
		Score:    post.Score(),
	})

	if err != nil {
		return -1, err
	}

	return createdId, nil
}

func (p *PostgresPostsRepository) Delete(ctx context.Context, postId int64) error {
	return p.querys.DeletePost(ctx, postId)
}

func (p *PostgresPostsRepository) List(ctx context.Context, params usecase.FindPostParams) ([]*usecase.Post, error) {
	dbPosts, err := p.querys.ListPosts(ctx, sqlc.ListPostsParams{
		Limit:  int32(params.PageSize),
		Offset: int32(params.PageNumber),
	})

	if err != nil {
		return nil, err
	}

	posts := make([]*usecase.Post, len(dbPosts))

	for i := 0; i < len(dbPosts); i++ {
		posts[i] = &usecase.Post{}
		convertPosts(dbPosts[i], posts[i])
	}

	return posts, nil
}

func convertPosts(src *sqlc.Post, dst *usecase.Post) {
	dst.PostId = src.PostID
	dst.Content = src.Content
	dst.PostedOn = src.PostedOn
	dst.PostedBy = src.PostedBy
	dst.Score = int(src.Score)
}

func NewPostgresConnection(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
