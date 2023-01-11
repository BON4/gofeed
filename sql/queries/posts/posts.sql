-- name: CreatePost :one
INSERT INTO Posts (
       content,
       posted_on,
       posted_by,
       score
) VALUES (
  $1, $2, $3, $4
) RETURNING post_id;

-- name: DeletePost :exec
DELETE FROM Posts WHERE post_id = $1;

-- name: ListPosts :many
SELECT * FROM Posts offset $2 limit $1;

-- name: GetRatePost :one
SELECT * FROM RatedPosts WHERE post_id = $1 and account = $2;

-- name: RatePost :exec
INSERT INTO RatedPosts (
       post_id,
       account,
       rated_score
) VALUES (
  $1, $2, $3
)
ON CONFLICT (post_id, account)
DO UPDATE
   SET rated_score = $3;

-- name: AddPostScore :one
UPDATE Posts
SET
	score = score + sqlc.arg('score')
WHERE
	post_id = sqlc.arg('post_id')
RETURNING *;
