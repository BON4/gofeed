// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: posts.sql

package sqlc

import (
	"context"
	"time"
)

const addPostScore = `-- name: AddPostScore :one
UPDATE Posts
SET
	score = score + $1
WHERE
	post_id = $2
RETURNING post_id, content, posted_on, posted_by, score
`

type AddPostScoreParams struct {
	Score  int32 `json:"score"`
	PostID int64 `json:"post_id"`
}

func (q *Queries) AddPostScore(ctx context.Context, arg AddPostScoreParams) (*Post, error) {
	row := q.db.QueryRowContext(ctx, addPostScore, arg.Score, arg.PostID)
	var i Post
	err := row.Scan(
		&i.PostID,
		&i.Content,
		&i.PostedOn,
		&i.PostedBy,
		&i.Score,
	)
	return &i, err
}

const createPost = `-- name: CreatePost :one
INSERT INTO Posts (
       content,
       posted_on,
       posted_by,
       score
) VALUES (
  $1, $2, $3, $4
) RETURNING post_id
`

type CreatePostParams struct {
	Content  string    `json:"content"`
	PostedOn time.Time `json:"posted_on"`
	PostedBy string    `json:"posted_by"`
	Score    int32     `json:"score"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.Content,
		arg.PostedOn,
		arg.PostedBy,
		arg.Score,
	)
	var post_id int64
	err := row.Scan(&post_id)
	return post_id, err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM Posts WHERE post_id = $1
`

func (q *Queries) DeletePost(ctx context.Context, postID int64) error {
	_, err := q.db.ExecContext(ctx, deletePost, postID)
	return err
}

const getRatePost = `-- name: GetRatePost :one
SELECT post_id, account, rated_score FROM RatedPosts WHERE post_id = $1 and account = $2
`

type GetRatePostParams struct {
	PostID  int64  `json:"post_id"`
	Account string `json:"account"`
}

func (q *Queries) GetRatePost(ctx context.Context, arg GetRatePostParams) (*Ratedpost, error) {
	row := q.db.QueryRowContext(ctx, getRatePost, arg.PostID, arg.Account)
	var i Ratedpost
	err := row.Scan(&i.PostID, &i.Account, &i.RatedScore)
	return &i, err
}

const listPosts = `-- name: ListPosts :many
SELECT post_id, content, posted_on, posted_by, score FROM Posts offset $2 limit $1
`

type ListPostsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListPosts(ctx context.Context, arg ListPostsParams) ([]*Post, error) {
	rows, err := q.db.QueryContext(ctx, listPosts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.PostID,
			&i.Content,
			&i.PostedOn,
			&i.PostedBy,
			&i.Score,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const ratePost = `-- name: RatePost :exec
INSERT INTO RatedPosts (
       post_id,
       account,
       rated_score
) VALUES (
  $1, $2, $3
)
ON CONFLICT (post_id, account)
DO UPDATE
   SET rated_score = $3
`

type RatePostParams struct {
	PostID     int64  `json:"post_id"`
	Account    string `json:"account"`
	RatedScore int32  `json:"rated_score"`
}

func (q *Queries) RatePost(ctx context.Context, arg RatePostParams) error {
	_, err := q.db.ExecContext(ctx, ratePost, arg.PostID, arg.Account, arg.RatedScore)
	return err
}
